package google

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type CallbackServer struct {
	HTTP *http.Server
	Chan chan string
}

func (C *CallbackServer) callbackHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("close").Parse(`<!DOCTYPE html>
	<html>
		<head>
			<title>Google Photos Downloader</title>
		</head>
		<body>
			<p style="margin-left: auto; margin-right: auto; width: 10em;">You may now close this window</p>
		</body>
	</html>`)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
	C.Chan <- r.FormValue("code")
	return
}

func NewCallbackServer(authCode string) *CallbackServer {
	S := new(CallbackServer)
	R := mux.NewRouter()
	R.HandleFunc("/auth/google/callback", S.callbackHandler)

	S.HTTP = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: R,
	}

	S.Chan = make(chan string, 1)

	return S
}

func (C *CallbackServer) Run() error {
	if err := C.HTTP.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (C *CallbackServer) Stop(ctx context.Context) error {
	if err := C.HTTP.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
