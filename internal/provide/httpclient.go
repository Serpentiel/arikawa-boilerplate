// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"net/http"
	"time"
)

// HTTPClient is a function which provides an http.Client instance.
func HTTPClient() *http.Client {
	// timeoutSeconds is the number of seconds that the HTTP client will wait before timing out.
	const timeoutSeconds = 10

	return &http.Client{
		Timeout: timeoutSeconds * time.Second,
	}
}
