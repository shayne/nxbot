package nxapi

import (
	"bytes"
	"fmt"
	"net/http"
)

// CameraInfo holds the unique ID and Name of the camera
type CameraInfo struct {
	ID   string `json:"cameraId"`
	Name string `json:"cameraName"`
}

// GetCameras performs an Nx API call to retrieve camera info
func (a *API) GetCameras() ([]CameraInfo, error) {
	var infos []CameraInfo
	err := a.GETRequest("getCameraUserAttributesList", &infos)
	if err != nil {
		return nil, fmt.Errorf("GetCameras failed: %v", err)
	}
	return infos, nil
}

// GetSnapshot returns raw image data from the Nx API for the given camera ID
func (a *API) GetSnapshot(id string) (*bytes.Buffer, error) {
	req, err := a.newAPIRequest("GET", "cameraThumbnail")
	if err != nil {
		return nil, fmt.Errorf("GetSnapshot failed: %v", err)
	}
	q := req.URL.Query()
	q.Add("cameraId", id)
	req.URL.RawQuery = q.Encode()

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetSnapshot failed: %v", err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetSnapshot failed: %v", err)
	}
	resp.Body.Close()

	return &buf, nil
}
