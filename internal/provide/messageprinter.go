// Package provide is the package that contains all of the providers for the dependency injection container.
package provide

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// MessagePrinter is a function which provides a message.Printer instance.
func MessagePrinter() *message.Printer {
	return message.NewPrinter(language.English)
}
