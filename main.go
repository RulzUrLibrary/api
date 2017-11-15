package main

import (
	"github.com/rulzurlibrary/api/app"
)

func main() {
	rulz := app.New(app.NewDefaultInitializer())
	// Start application
	rulz.Logger.Fatal(rulz.Start())
}
