package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/rajabhishekmaurya/ecom/internal"
)

var (
	app = flag.String("app", "unknown", "Used to specify target app to run at startup.")
)

func main() {
	flag.Parse()
	if !flag.Parsed() {
		log.Fatal("failed to parse command line flags")
	}

	runner, err := internal.New(app)
	if err != nil {
		log.Fatal(err)
	}

	if err := runner.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
