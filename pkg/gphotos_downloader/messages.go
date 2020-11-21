package gp2app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type GetMediaItemsRequest struct {
	PageSize  int    `json:"pageSize"`
	PageToken string `json:"pageToken"`
}

type MediaItem struct {
	ID            string `json:"id"`
	Description   string `json:"description"`
	ProductURL    string `json:"productUrl"`
	BaseURL       string `json:"baseUrl"`
	MIMEType      string `json:"mimeType"`
	Filename      string `json:"filename"`
	MediaMetadata struct {
		Width        string `json:"width"`
		Height       string `json:"height"`
		CreationTime string `json:"creationTime"`
		Photo        struct {
			CameraMake      string  `json:"cameraMake"`
			CambraModel     string  `json:"cameraModel"`
			FocalLength     float32 `json:"focalLength"`
			ApertureFNumber float32 `json:"apertureFNumber"`
			ISOEquivalent   int     `json:"isoEquivalent"`
			ExposureTime    string  `json:"exposureTime"`
		} `json:"photo"`
		Video struct {
			CameraMake  string  `json:"cameraMake"`
			CameraModel string  `json:"cameraModel"`
			FPS         float32 `json:"fps"`
			Status      string  `json:"status"`
		} `json:"video"`
	} `json:"mediaMetadata"`
	ContributorInfo struct {
		ProfilePictureBaseURL string `json:"profilePictureBaseUrl"`
		DisplayName           string `json:"displayName"`
	} `json:"contributorInfo"`
}

type GetMediaItemsResponse struct {
	MediaItems    []MediaItem `json:"mediaItems"`
	NextPageToken string      `json:"nextPageToken"`
}

func (G *GetMediaItemsResponse) Unmarshal(httpBody io.ReadCloser) error {
	b, err := ioutil.ReadAll(httpBody)
	if err != nil {
		return fmt.Errorf("unable to read GetMediaItems response due to error: %+v", err)
	}

	if err := json.Unmarshal(b, G); err != nil {
		return fmt.Errorf("unable to unmarshal the GetMediaItems response due to error: %+v", err)
	}

	return nil
}
