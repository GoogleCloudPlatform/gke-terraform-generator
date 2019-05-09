/*
Copyright 2018 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"os"

	"k8s.io/klog"

	"partner-code.googlesource.com/gke-terraform-generator/pkg/api"
	"partner-code.googlesource.com/gke-terraform-generator/pkg/files"
	"partner-code.googlesource.com/gke-terraform-generator/pkg/templates"
)

var (
	// outDir is the directory that gke-tf write the terraform files to.
	outDir string
	// configFile is the name of the yaml configuration file
	configFile string
	// projectId is the gcp project id.
	projectId string
	// gkeTF is the api cluster representation.
	gkeTF *api.GkeTF
	// genCommend command to generate GKE Terraform file
	genCommand = &cobra.Command{
		Use:   "gen",
		Short: "Generates GKE TF File with given locals",
		// TODO refactor into a function
		Run: NewGen,
	}
)

func init() {
	// TODO on the fence about this, should we have the user specify the directory
	defaultDir, err := os.Getwd()
	if err != nil {
		klog.Errorf("Error getting directory: %v", err)
		exitWithError(err)
	}
	gkeTF = &api.GkeTF{
		Spec: api.ClusterSpec{
			Zones:     new([]string),
			NodePools: &[]*api.GkeNodePool{},
			Tags:      new([]string),
		},
	}

	RootCMD.AddCommand(genCommand)
	// Add root flags so we can get logging flags
	genCommand.Flags().AddFlagSet(RootCMD.Flags())
	genCommand.Flags().StringVarP(&outDir, "directory", "d", defaultDir, "output directory")
	genCommand.Flags().StringVarP(&configFile, "file", "f", "", "config yaml file")
	genCommand.Flags().StringVarP(&projectId, "project-id", "p", "", "gcp project id")
	if err := cobra.MarkFlagRequired(genCommand.Flags(), "file"); err != nil {
		exitWithError(err)
	}
}

// NewGen is the entry point for cobra for the gen command.
func NewGen(cmd *cobra.Command, args []string) {
	klog.Info(header)
	err := checkCliArgs()
	if err != nil {
		klog.Errorf("Error checking cli arguments: %v", err)
		os.Exit(1)
	}

	gkeTF, err = api.UnmarshalGkeTF(configFile)
	if err != nil {
		klog.Errorf("Error umarshalling yaml: %v", err)
		os.Exit(1)
	}

	klog.Infof("Creating terraform for your GKE cluster %s.", gkeTF.Name)

	// set project id.  This will also override the value if it exists in the
	// YAML file.
	if projectId != "" {
		gkeTF.Spec.ProjectId = projectId
	}

	err = api.SetApiDefaultValues(gkeTF, configFile)
	if err != nil {
		klog.Errorf("Error setting api defaults: %v", err)
		os.Exit(1)
	}

	err = api.ValidateYamlInput(gkeTF)
	if err != nil {
		klog.Errorf("Error validating api values: %v", err)
		os.Exit(1)
	}

	// this creates a NewGKETemplates struct and runs CopyTo.
	// This func does all the grunt work of processing each go template and writing the
	// terraform results to a file.
	err = templates.NewGKETemplates().CopyTo(outDir, gkeTF)
	if err != nil {
		klog.Errorf("Error creating terraform: %v", err)
		os.Exit(1)
	}
}

// checkCliArgs in essence checks the cli arguments to ensure that the proper
// flags have been set by the user.
func checkCliArgs() error {

	if outDir == "" {
		return errors.New("--directory option must be set with a directory name")
	}

	test, err := files.IsWritable(outDir)
	if err != nil {
		return errors.New(fmt.Sprintf("Error openning directory: %s ... %s", outDir, err.Error()))
	}

	if !test {
		return errors.New("Unable to open directory: " + outDir)
	}

	if configFile == "" {
		return errors.New("--file option must be set with a file name")
	}

	test, err = files.IsFile(configFile)

	if err != nil {
		return errors.New(fmt.Sprintf("Error openning config file: %s ... %s", configFile, err.Error()))
	}

	if !test {
		return errors.New("Configuration file is not found: " + configFile)
	}

	return nil
}
