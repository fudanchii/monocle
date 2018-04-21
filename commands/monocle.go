package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var MonocleCmd = &cobra.Command{
	Use:   "monocle",
	Short: "monocle runs app, build, and test",
	Long: `Monocle can be used to run, build, and test app.

Monocle
Nurahmadie <nurahmadie@gmail.com>
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("this print something.")
		return nil
	},
}

func Execute() {
	if c, err := MonocleCmd.ExecuteC(); err != nil {
		c.Println(err.Error())
		c.Println(c.UsageString())
		os.Exit(-1)
	}
}
