package cli

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	dkv "github.com/docker/libkv/store"
	"github.com/emccode/polly/core/store"
	"github.com/spf13/cobra"
)

func (c *CLI) initStoreCmdsAndFlags() {
	c.initStoreCmds()
	c.initStoreFlags()
}

func (c *CLI) initStoreCmds() {

	c.storeCmd = &cobra.Command{
		Use:   "store",
		Short: "The store manager",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}
	c.c.AddCommand(c.storeCmd)

	c.storeEraseCmd = &cobra.Command{
		Use:   "erase",
		Short: "Erase the persistent store",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.p.Store.EraseStore(); err != nil {
				log.Fatal(err)
			}
		},
	}
	c.storeCmd.AddCommand(c.storeEraseCmd)

	c.storeGetCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get volumes keys from store",
		Aliases: []string{"ls", "list"},
		Run: func(cmd *cobra.Command, args []string) {
			vols, err := c.p.Store.GetVolumeIds()
			if err != nil {
				log.Fatal(err)
			}

			var kvl []*dkv.KVPair
			for _, vol := range vols {
				key, _ := c.p.Store.GenerateObjectKey(store.VolumeAdminLabelsType, vol)
				kv, _ := c.p.Store.List(key)

				kvl = append(kvl, kv...)

				key, _ = c.p.Store.GenerateObjectKey(store.VolumeType, vol)
				kv, _ = c.p.Store.List(key)

				kvl = append(kvl, kv...)

				key, _ = c.p.Store.GenerateObjectKey(store.VolumeInternalLabelsType, vol)
				kvl, _ = c.p.Store.List(key)

				kvl = append(kvl, kv...)
			}

			kvmap := make(map[string]string)
			for _, kv := range kvl {
				kvmap[kv.Key] = string(kv.Value)
			}

			out, err := c.marshalOutput(&kvmap)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)

		},
	}
	c.storeCmd.AddCommand(c.storeGetCmd)

}

func (c *CLI) initStoreFlags() {
}
