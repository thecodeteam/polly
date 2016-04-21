package main

import (
	"log"

	gofig "github.com/akutz/gofig"
	core "github.com/emccode/polly/core"
)

var pollyCore *core.PollyCore

func main() {
	config := gofig.New()

	if err := config.ReadConfigFile("/etc/polly/config.yml"); err != nil {
		panic(err)
	}

	var err error
	pollyCore, err = core.NewWithConfig(config)
	if err != nil {
		log.Fatal("Fatal: ", err)
	}
	log.Print("PollyStore Type: ", pollyCore.PollyStore.StoreType())

}
