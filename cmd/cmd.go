package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "undefined"

var RootCMD = &cobra.Command{
	Use:   "gketf [command]",
	Short: "CLI Interface for GKE TF util",
}

func Execute() {
	if err := RootCMD.Execute(); err != nil {
		PrintErrln(err.Error())
		os.Exit(-1)
	}
}

// PrintErrln ...
func PrintErrln(msg interface{}) (int, error) {
	return fmt.Fprintln(os.Stderr, msg)
}
