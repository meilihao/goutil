package http

import (
	"crypto/tls"
	"net/http"
)

var (
	SecureClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
)
