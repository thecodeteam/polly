package cli

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/emccode/polly/api/types"
	"github.com/spf13/cobra"
)

func (c *CLI) initVolumeCmdsAndFlags() {
	c.initVolumeCmds()
	c.initVolumeFlags()
}

func (c *CLI) initVolumeCmds() {

	c.volumeCmd = &cobra.Command{
		Use:   "volume",
		Short: "The volume manager",
		Run: func(cmd *cobra.Command, args []string) {
			if isHelpFlags(cmd) {
				cmd.Usage()
			} else {
				c.volumeGetCmd.Run(c.volumeGetCmd, args)
			}
		},
	}
	c.c.AddCommand(c.volumeCmd)

	c.volumeGetCmd = &cobra.Command{
		Use:     "get",
		Short:   "Get one or more volumes",
		Aliases: []string{"ls", "list"},
		Run: func(cmd *cobra.Command, args []string) {

			var av []*types.Volume
			var err error

			if c.volumeID != "" {
				v, err := c.pc.VolumeInspect(c.volumeID)
				if err != nil {
					log.Fatal(err)
				}

				out, err := c.marshalOutput(&v)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(out)
			} else if c.all {
				av, err = c.pc.VolumesAll()
				if err != nil {
					log.Fatal(err)
				}
			} else {
				av, err = c.pc.Volumes()
				if err != nil {
					log.Fatal(err)
				}
			}

			if len(av) > 0 {
				out, err := c.marshalOutput(&av)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(out)
			}
		},
	}
	c.volumeCmd.AddCommand(c.volumeGetCmd)

	c.volumeOfferCmd = &cobra.Command{
		Use:   "offer",
		Short: "Offer a volume to schedulers",
		Run: func(cmd *cobra.Command, args []string) {
			av, err := c.pc.VolumeOffer(c.volumeID, c.schedulers)
			if err != nil {
				log.Fatal(err)
			}

			out, err := c.marshalOutput(&av)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)
		},
	}
	c.volumeCmd.AddCommand(c.volumeOfferCmd)

	c.volumeOfferRevokeCmd = &cobra.Command{
		Use:   "revoke",
		Short: "Revoke an offer of a volume to schedulers",
		Run: func(cmd *cobra.Command, args []string) {
			av, err := c.pc.VolumeOfferRevoke(c.volumeID, c.schedulers)
			if err != nil {
				log.Fatal(err)
			}

			out, err := c.marshalOutput(&av)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)
		},
	}
	c.volumeCmd.AddCommand(c.volumeOfferRevokeCmd)

	c.volumeLabelCmd = &cobra.Command{
		Use:   "label",
		Short: "Create labels on a volume",
		Run: func(cmd *cobra.Command, args []string) {
			av, err := c.pc.VolumeLabel(c.volumeID, c.labels)
			if err != nil {
				log.Fatal(err)
			}

			out, err := c.marshalOutput(&av)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)
		},
	}
	c.volumeCmd.AddCommand(c.volumeLabelCmd)

	c.volumeLabelRemoveCmd = &cobra.Command{
		Use:     "labelremove",
		Short:   "Remove labels from a volume",
		Aliases: []string{"lr"},
		Run: func(cmd *cobra.Command, args []string) {
			av, err := c.pc.VolumeLabelsRemove(c.volumeID, c.labels)
			if err != nil {
				log.Fatal(err)
			}

			out, err := c.marshalOutput(&av)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)
		},
	}
	c.volumeCmd.AddCommand(c.volumeLabelRemoveCmd)

	c.volumeCreateCmd = &cobra.Command{
		Use:     "create",
		Short:   "Creates a volume",
		Aliases: []string{"new"},
		Run: func(cmd *cobra.Command, args []string) {
			av, err := c.pc.VolumeCreate(c.serviceName, c.name, c.volumeType,
				c.size, c.IOPS, c.availabilityZone, c.schedulers, c.labels,
				nil)
			if err != nil {
				log.Fatal(err)
			}

			out, err := c.marshalOutput(&av)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(out)
		},
	}
	c.volumeCmd.AddCommand(c.volumeCreateCmd)

	c.volumeRemoveCmd = &cobra.Command{
		Use:     "remove",
		Short:   "Removes a volume",
		Aliases: []string{"rm", "delete"},
		Run: func(cmd *cobra.Command, args []string) {
			err := c.pc.VolumeRemove(c.volumeID)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	c.volumeCmd.AddCommand(c.volumeRemoveCmd)

}

func (c *CLI) initVolumeFlags() {
	c.volumeGetCmd.Flags().BoolVar(&c.all, "all", false, "all")
	c.volumeGetCmd.Flags().StringVar(&c.volumeID, "volumeid", "", "volumeid")
	c.volumeOfferCmd.Flags().StringVar(&c.volumeID, "volumeid", "", "volumeid")
	c.volumeOfferCmd.Flags().StringSliceVar(&c.schedulers, "scheduler", []string{""}, "scheduler")
	c.volumeOfferRevokeCmd.Flags().StringVar(&c.volumeID, "volumeid", "", "volumeid")
	c.volumeOfferRevokeCmd.Flags().StringSliceVar(&c.schedulers, "scheduler", []string{""}, "scheduler")
	c.volumeLabelCmd.Flags().StringVar(&c.volumeID, "volumeid", "", "volumeid")
	c.volumeLabelCmd.Flags().StringSliceVar(&c.labels, "label", []string{""}, "label")
	c.volumeLabelRemoveCmd.Flags().StringVar(&c.volumeID, "volumeid", "", "volumeid")
	c.volumeLabelRemoveCmd.Flags().StringSliceVar(&c.labels, "label", []string{""}, "label")
	c.volumeCreateCmd.Flags().StringVar(&c.name, "name", "", "name")
	c.volumeCreateCmd.Flags().StringVar(&c.serviceName, "servicename", "", "servicename")
	c.volumeCreateCmd.Flags().StringVar(&c.volumeType, "type", "", "type")
	c.volumeCreateCmd.Flags().StringVar(&c.volumeID, "volumeid", "", "volumeid")
	c.volumeCreateCmd.Flags().Int64Var(&c.IOPS, "iops", 0, "IOPS")
	c.volumeCreateCmd.Flags().Int64Var(&c.size, "size", 0, "size")
	c.volumeCreateCmd.Flags().StringVar(&c.availabilityZone, "availabilityzone", "", "availabilityzone")
	c.volumeCreateCmd.Flags().StringSliceVar(&c.labels, "label", []string{""}, "label")
	c.volumeCreateCmd.Flags().StringSliceVar(&c.schedulers, "scheduler", []string{""}, "scheduler")
	c.volumeRemoveCmd.Flags().StringVar(&c.volumeID, "volumeid", "", "volumeid")

	c.addOutputFormatFlag(c.volumeCmd.Flags())
	c.addOutputFormatFlag(c.volumeGetCmd.Flags())
	c.addOutputFormatFlag(c.volumeOfferCmd.Flags())
}
