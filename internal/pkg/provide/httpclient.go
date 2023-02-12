// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"net/http"
	"time"
)

// HTTPClient is a function which provides an http.Client instance.
func HTTPClient() *http.Client {
	const timeoutSeconds = 10

	return &http.Client{
		Timeout: timeoutSeconds * time.Second,
	}
}
