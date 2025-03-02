package imageutils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io"
	"net/http"
)

// DownloadImage fetches an image from a URL and returns it as a Base64-encoded string
func DownloadImageAsBase64(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read image bytes
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Encode as Base64
	return base64.StdEncoding.EncodeToString(imgData), nil
}

// DecodeBase64Image takes a Base64 string and returns an image.Image
func DecodeBase64Image(encoded string) (image.Image, error) {
	// Decode Base64
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// Convert bytes to image
	img, err := jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return img, nil
}
