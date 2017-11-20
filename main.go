package main

import (
	"github.com/RulzUrLibrary/api/app"
)

func main() {
	rulz := app.New(app.NewDefaultInitializer())
	// Start application
	rulz.Logger.Fatal(rulz.Start())
}
