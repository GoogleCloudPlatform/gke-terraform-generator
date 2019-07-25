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
	"strings"

	"github.com/spf13/cobra"

	"os"

	"k8s.io/klog"

	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/api"
	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/files"
	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/templates"
)

var (
	// outDir is the directory that gke-tf write the terraform files to.
	outDir string
	// configFile is the name of the yaml configuration file
	configFile string
	// projectID is the gcp project id.
	projectID string
	// overwriteFile boolean determines whether TF files can be written over
	overwriteFile bool
	// gkeTF is the api cluster representation.
	gkeTF *api.GkeTF
	// tfType is the type of terraform
	tfTypeStr string
	// tfType is the type of terraform
	tfType templates.TFType
)

// NewGenCommand is the entry point for cobra for the gen command.
func NewGenCommand() *cobra.Command {
	defaultDir, err := os.Getwd()
	if err != nil {
		klog.Errorf("Error getting directory: %v", err)
		exitWithError(err)
	}

	genCommand := &cobra.Command{
		Use:   "gen",
		Short: "Generates GKE TF File with given locals",
	}
	// Add root flags so we can get logging flags
	genCommand.Flags().AddFlagSet(RootCMD.Flags())
	// TODO add long help descriptions for these
	genCommand.Flags().StringVarP(&outDir, "directory", "d", defaultDir, "output directory")
	genCommand.Flags().StringVarP(&configFile, "file", "f", "", "config yaml file")
	genCommand.Flags().StringVarP(&projectID, "project-id", "p", "", "gcp project id")
	genCommand.Flags().StringVarP(&tfTypeStr, "tf-type", "t", "Vanilla", "terraform types are CFT or Vanilla")
	genCommand.Flags().BoolVarP(&overwriteFile, "overwrite-file", "o", false, "overwrite file flag")

	if err := cobra.MarkFlagRequired(genCommand.Flags(), "file"); err != nil {
		exitWithError(err)
	}

	genCommand.Run = func(cmd *cobra.Command, args []string) {
		klog.Info(header)
		err := checkCliArgs()
		if err != nil {
			klog.Errorf("Error checking cli arguments: %v", err)
			os.Exit(1)
		}

		gkeTF, err = api.UnmarshalGkeTF(configFile)
		if err != nil {
			klog.Errorf("Error unmarshaling the configuration file: %v", err)
			exitWithError(err)
		}

		klog.Infof("Creating terraform for your GKE cluster %s.", gkeTF.Name)

		// set project id.  This will also override the value if it exists in the
		// YAML file.
		if projectID != "" {
			gkeTF.Spec.ProjectId = projectID
		}

		err = api.SetApiDefaultValues(gkeTF, configFile)
		if err != nil {
			klog.Errorf("Error setting api defaults: %v", err)
			exitWithError(err)
		}

		err = api.ValidateYamlInput(gkeTF)
		if err != nil {
			klog.Errorf("Error validating api values: %v", err)
			exitWithError(err)
		}

		template, err := templates.NewGKETemplates(tfType)
		if err != nil {
			klog.Errorf("Error creating setting up terraform templates: %v", err)
			exitWithError(err)
		}

		// this runs CopyTo on the template created.
		// This func does all the grunt work of processing each go template and writing the
		// terraform results to a file.
		err = template.CopyTo(overwriteFile, outDir, gkeTF)
		if err != nil {
			klog.Errorf("Error creating terraform: %v", err)
			exitWithError(err)
		}

	}
	return genCommand
}

// checkCliArgs in essence checks the cli arguments to ensure that the proper
// flags have been set by the user.
func checkCliArgs() error {

	if outDir == "" {
		return errors.New("--directory option must be set with a directory name")
	}

	if err := files.CreateDirIfNotExist(outDir); err != nil {
		return fmt.Errorf("Error creating directory: %s ... %s", outDir, err.Error())
	}

	test, err := files.IsWritable(outDir)
	if err != nil {
		return fmt.Errorf("Error openning directory: %s ... %s", outDir, err.Error())
	}

	if !test {
		return errors.New("Unable to open directory: " + outDir)
	}

	if configFile == "" {
		return errors.New("--file option must be set with a file name")
	}

	test, err = files.IsFile(configFile)

	if err != nil {
		return fmt.Errorf("Error openning config file: %s ... %s", configFile, err.Error())
	}

	if !test {
		return errors.New("Configuration file is not found: " + configFile)
	}

	switch strings.ToUpper(tfTypeStr) {
	case "CFT":
		tfType = templates.CFT
	case "VANILLA":
		tfType = templates.VANILLA
	default:
		return errors.New("unable to determine terraform type, please set the -t flag with CFT or Vanilla")
	}

	return nil
}
