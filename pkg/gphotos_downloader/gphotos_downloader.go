package gp2app

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/korovkin/limiter"

	"github.com/google/uuid"
	internal "github.com/j4ng5y/gphotos_downloader/internal/gphotos_downloader"
	callback "github.com/j4ng5y/gphotos_downloader/pkg/gphotos_downloader/callbacks/google"
	"golang.org/x/oauth2"
)

type Google struct {
	client *http.Client
}

// Run creates a new instance of the Google struct.
func Run() (err error) {
	G := new(Google)
	if err := G.doOauth(); err != nil {
		return err
	}

	if err := G.getMediaItems(&GetMediaItemsRequest{
		PageSize: 50,
	}); err != nil {
		return err
	}

	return err
}

// openBrowser simply opens the OS default browser to the Google sign-in page
func (Google) openBrowser(u string) error {
	switch runtime.GOOS {
	case "linux":
		if err := exec.Command("xdg-open", u).Start(); err != nil {
			return err
		}
		return nil
	case "windows":
		if err := exec.Command("rundll32", "url.dll,FileProtocolHandler", u).Start(); err != nil {
			return err
		}
		return nil
	case "darwin":
		if err := exec.Command("open", u).Start(); err != nil {
			return err
		}
		return nil
	default:
		fmt.Printf("\n\nPlease navigate to the following url to complete the OAuth process: '%s'", u)
		return nil
	}
}

// doOauth runs the oauth process
func (G *Google) doOauth() (err error) {
	// stateCode is used as a means to authenticate the oauth response
	stateCode := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	u := internal.OauthConf.AuthCodeURL(stateCode)
	if err := G.openBrowser(u); err != nil {
		return err
	}

	S := callback.NewCallbackServer(stateCode)
	go func() {
		if err := S.Run(); err != nil {
			if err == http.ErrServerClosed {
				// Do nothing
			} else {
				log.Fatal(err)
			}
		}
	}()
	code := <-S.Chan
	if err := S.Stop(ctx); err != nil {
		return err
	}

	token, err := internal.OauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}

	G.client = internal.OauthConf.Client(oauth2.NoContext, token)

	return err
}

// getMediaItem gets all mediaItems from a single page of mediaItems
func (G Google) getMediaItem(req *GetMediaItemsRequest, U *url.URL) (*GetMediaItemsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, U.String(), nil)
	if err != nil {
		return nil, err
	}

	httpResp, err := G.client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	resp := &GetMediaItemsResponse{}

	if err := resp.Unmarshal(httpResp.Body); err != nil {
		return nil, err
	}

	req.PageToken = resp.NextPageToken

	return resp, nil
}

// getMediaItems runs the process of getting all of the listed media items
func (G *Google) getMediaItems(req *GetMediaItemsRequest) error {
	goroLimiter := limiter.NewConcurrencyLimiter(50)
	U, err := url.Parse(fmt.Sprintf("https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=%d", req.PageSize))
	if err != nil {
		return err
	}

	i := 0
	// I don't care for Labels and Goto statements, but they work pretty decently for pagination
LoopStart:
	i = i + 1
	resp, err := G.getMediaItem(req, U)
	if err != nil {
		return err
	}

	for _, v := range resp.MediaItems {
		u, err := url.Parse(v.BaseURL)
		if err != nil {
			return err
		}
		finalU, err := url.Parse(fmt.Sprintf("%s=d", u.String()))
		if err != nil {
			return err
		}
		goroLimiter.Execute(func() {
			G.download(v.Filename, finalU)
		})
	}

	if req.PageToken != "" {
		U, err = url.Parse(fmt.Sprintf("https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=%d&pageToken=%s", req.PageSize, req.PageToken))
		if err != nil {
			return err
		}
		goto LoopStart
	}

	goroLimiter.Wait()
	return nil
}

func (G Google) download(fn string, u *url.URL) {
	fmt.Printf("Downloading %s\n", fn)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	resp, err := G.client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if _, err := os.Stat(fn); err != nil {
		// do nothing
		return
	}

	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Println(err)
		return

	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Println(err)
		return
	}
}
