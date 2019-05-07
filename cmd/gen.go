package cmd

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"partner-code.googlesource.com/gke-terraform-generator/generator"
	"partner-code.googlesource.com/gke-terraform-generator/templates"
)

// genConfig flag vars
var genConfigVars *generator.Vars
var outDir string
var vars []string

// genCMD command to generate GKE Terraform file
var genCMD = &cobra.Command{
	Use:   "gen",
	Short: "Generates GKE TF File with given locals",
	Run: func(cmd *cobra.Command, args []string) {
		// Override any individually set vars
		genConfigVars.Set(strings.Join(vars, ","))
		if err := ioutil.WriteFile(path.Join(outDir, "/vars.tf"), genConfigVars.TFVars(), 0644); err != nil {
			PrintErrln(err)
		}
		if err := ioutil.WriteFile(path.Join(outDir, "/main.tf"), templates.Main, 0644); err != nil {
			PrintErrln(err)
		}
	},
}

// init sets up cli command and flags
func init() {
	var err error
	outDir, err = os.Getwd()
	if err != nil {
		PrintErrln(err)
	}
	genConfigVars = &generator.Vars{}
	RootCMD.AddCommand(genCMD)
	genCMD.Flags().StringVar(&outDir, "out", outDir, "output directory")
	genCMD.Flags().Var(genConfigVars, "vars", "config vars")
	genCMD.Flags().StringSliceVar(&vars, "var", []string{}, "config var (ex. --var=\"private=true\"")
}
