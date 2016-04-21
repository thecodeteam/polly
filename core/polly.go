package core

import (
	"log"

	gofig "github.com/akutz/gofig"
	store "github.com/emccode/polly/pkg/store"
)

//PollyCore this represents the "core" functionality for Polly
type PollyCore struct {
	PollyStore store.IPollyStore
}

//NewWithConfig This initializes new instance of this library
func NewWithConfig(config gofig.Config) (pc *PollyCore, err error) {
	ps, err := store.NewWithConfig(config)
	if err != nil {
		log.Fatal("Fatal initializing PollyStore: ", err)
		return nil, err
	}
	log.Print("StoreType: ", ps.StoreType())

	myPc := new(PollyCore)
	myPc.PollyStore = ps

	//TODO

	return myPc, nil
}
