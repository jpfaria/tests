// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023 Datadog, Inc.

// This file is pulled from datadog-agent/pkg/util/http (Only removing agent SSL config and unused funcs)

package httputils

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func createTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 30 * time.Second,
			// Enables TCP keep-alives to detect broken connections
			KeepAlive: 30 * time.Second,
			// Disable RFC 6555 Fast Fallback ("Happy Eyeballs")
			FallbackDelay: -1 * time.Nanosecond,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 5,
		// This parameter is set to avoid connections sitting idle in the pool indefinitely
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy:                 http.ProxyFromEnvironment,
	}
}

// Get is a high level helper to query a URL and return its body as a string
func Get(ctx context.Context, URL string, headers map[string]string, timeout time.Duration) (string, error) {
	client := http.Client{
		Transport: createTransport(),
		Timeout:   timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return "", err
	}

	for header, value := range headers {
		req.Header.Add(header, value)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	return parseResponse(res, "GET", URL)
}

func parseResponse(res *http.Response, method string, URL string) (string, error) {
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code %d trying to %s %s", res.StatusCode, method, URL)
	}

	defer res.Body.Close()
	all, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response from %s: %s", URL, err)
	}

	return string(all), nil
}
