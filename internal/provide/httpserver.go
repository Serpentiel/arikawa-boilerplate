// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// HTTPServer is a function which provides an http.Server instance.
func HTTPServer(v *viper.Viper) *http.Server {
	const (
		// readTimeoutSeconds is the maximum duration for reading the entire request, including the body.
		readTimeoutSeconds = 1

		// readHeaderTimeoutSeconds is the amount of time allowed to read request headers.
		readHeaderTimeoutSeconds = 2

		// writeTimeoutSeconds is the maximum duration before timing out writes of the response.
		writeTimeoutSeconds = 1

		// idleTimeoutSeconds is the maximum amount of time to wait for the next request when keep-alives are enabled.
		idleTimeoutSeconds = 30
	)

	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", v.GetString("http.host"), v.GetInt("http.port")),
		ReadTimeout:       readTimeoutSeconds * time.Second,
		ReadHeaderTimeout: readHeaderTimeoutSeconds * time.Second,
		WriteTimeout:      writeTimeoutSeconds * time.Second,
		IdleTimeout:       idleTimeoutSeconds * time.Second,
	}
}
