// Package asset is the package that contains all of the assets for the application.
package asset

import _ "embed" // embed is imported as a blank import to let us embed the assets.

// ExampleConfig is the example config for the application.
//
//go:embed .arikawa-boilerplate.example.yaml
var ExampleConfig []byte
