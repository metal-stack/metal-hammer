package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	pixiecore "github.com/metal-stack/pixie/api"
)

func fetchConfig(pixieURL string) (*pixiecore.MetalConfig, error) {
	certClient := http.Client{
		Timeout: 5 * time.Second,
	}
	ctx, httpcancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer httpcancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pixieURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := certClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	js, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var metalConfig pixiecore.MetalConfig
	if err := json.Unmarshal(js, &metalConfig); err != nil {
		return nil, fmt.Errorf("unable to unmarshal grpcConfig:%w", err)
	}
	return &metalConfig, nil
}
