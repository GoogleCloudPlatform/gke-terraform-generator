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
	"fmt"
	"io"

	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand is the entry point for cobra for the version command.
func NewVersionCommand(out io.Writer) (cmd *cobra.Command) {
	versionCMD := &cobra.Command{
		Use:   "version",
		Short: "Print gke-tf version",
	}
	versionCMD.Run = func(cmd *cobra.Command, args []string) {
		s := fmt.Sprintf("gke-tf version: %s", version.Version)
		_, err := fmt.Fprintf(out, "%s\n", s)
		if err != nil {
			exitWithError(err)
		}
	}
	return versionCMD
}
