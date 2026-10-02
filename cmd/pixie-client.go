package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	pixiecore "github.com/metal-stack/pixie/api"
)

func fetchPixieConfig(pixieURL, machineUUID string) (*pixiecore.V2MetalHammerConfigPayload, error) {
	u, err := url.Parse(pixieURL)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s://%s/config/%s", u.Scheme, u.Host, machineUUID)

	certClient := http.Client{
		Timeout: 5 * time.Second,
	}

	ctx, httpcancel := context.WithTimeout(context.Background(), certClient.Timeout)
	defer httpcancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
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

	var config pixiecore.V2MetalHammerConfigPayload

	if err := json.Unmarshal(js, &config); err != nil {
		return nil, fmt.Errorf("unable to unmarshal pixiecore response: %w", err)
	}

	return &config, nil
}
