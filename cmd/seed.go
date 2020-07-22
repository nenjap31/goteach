package cmd

import (
	"fmt"
	"goteach/config"
	"goteach/database/seeds"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.Flags().BoolVarP(&list, "list", "l", false, "list available seeds")
	seedCmd.SetUsageTemplate(seedUsage)
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run database seeder",
	Run:   seedHandler,
}

var seedUsage = `Usage:
  goteach seed [flags] [seederName [seederName]]

Example:
  goteach seed seedOne seedTwo

Flags:
  -h, --help   help for seed
  -l, --list   list available seeds
`

var list bool

func seedHandler(cmd *cobra.Command, args []string) {
	db := config.DB

	if list {
		fmt.Println("Available seeds:")
		for name := range seeds.SeederList {
			fmt.Println("- " + name)
		}
	} else if len(args) == 0 {
		for name, seeder := range seeds.SeederList {
			fmt.Println("Running " + name)
			seeder.Run(db)
		}
	} else {
		for _, name := range args {
			if seeder, ok := seeds.SeederList[name]; ok {
				seeder.Run(db)
			} else {
				fmt.Println(name + " not exist. run \"goteach seed -l\" to see available seeds.")
			}
		}
	}
}
