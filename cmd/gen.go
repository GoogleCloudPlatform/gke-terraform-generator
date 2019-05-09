package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"partner-code.googlesource.com/gke-terraform-generator/generator"
)

// flag vars
var outDir string
var dataFile string

// genCMD command to generate GKE Terraform file
var genCMD = &cobra.Command{
	Use:   "gen",
	Short: "Generates GKE TF File with given locals",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(dataFile)
		checkErr(err)

		var context interface{}
		err = yaml.Unmarshal(data, &context)

		// get module dir path for templates
		_, filename, _, ok := runtime.Caller(0)
		if !ok {
			checkErr(errors.New("error fetching templates"))
		}
		tplDir := path.Join(path.Dir(filename), "../templates")

		checkErr(generator.Generate(context, tplDir, outDir))
	},
}

// init sets up cli command and flags
func init() {
	RootCMD.AddCommand(genCMD)
	genCMD.Flags().StringVar(&outDir, "out", "./", "output directory")
	genCMD.Flags().StringVar(&dataFile, "data", "./data.yml", "data yaml file path")
}

// checkErr is a simple helper that checks errors and exits if err is not nil
func checkErr(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
