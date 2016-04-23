package polly

import (
	"log"

	gofig "github.com/akutz/gofig"
	core "github.com/emccode/polly/core"
)

//New init the lib
func New() (p *core.PollyCore, err error) {
	config := gofig.New()

	myErr := config.ReadConfigFile("/etc/polly/config.yml")
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}

	pollyCore, myErr := core.NewWithConfig(config)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}
	log.Print("PollyStore Type: ", pollyCore.PollyStore.StoreType())

	return pollyCore, nil
}

//NewWithConfigFile init the lib
func NewWithConfigFile(path string) (p *core.PollyCore, err error) {
	config := gofig.New()

	myErr := config.ReadConfigFile(path)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}

	pollyCore, myErr := core.NewWithConfig(config)
	if myErr != nil {
		log.Fatal("Fatal: ", myErr)
		return nil, myErr
	}
	log.Print("PollyStore Type: ", pollyCore.PollyStore.StoreType())

	return pollyCore, nil
}
