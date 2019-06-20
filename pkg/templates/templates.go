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

package templates

import (
	"bufio"
	"k8s.io/klog"
	"os"
	"path"
	"strings"
	"text/template"

	"partner-code.googlesource.com/gke-terraform-generator/pkg/api"
	"partner-code.googlesource.com/gke-terraform-generator/pkg/terraform"
)

type TerraformTemplate struct {
	FileName   string
	GoTemplate string
}

type GKETemplates struct {
	Templates []*TerraformTemplate
}

var templates = []*TerraformTemplate{
	{
		"main.tf",
		terraform.GKEMainTF,
	},
	{
		"network.tf",
		terraform.GKENetworkTF,
	},
	{
		"outputs.tf",
		terraform.GKEOutputsTF,
	},
	{
		"variables.tf",
		terraform.GKEVariablesTF,
	},
}

func NewGKETemplates() *GKETemplates {
	return &GKETemplates{
		Templates: templates,
	}
}

// CopyTo is used to copy all of the templates in the
// template directory to the given destination
func (gkeTemplates *GKETemplates) CopyTo(dst string, cluster *api.GkeTF) error {
	return gkeTemplates.processTemplates(dst, cluster)
}

func (gkeTemplates *GKETemplates) processTemplates(dst string, cluster *api.GkeTF) error {
	// TODO check if file exists prompt to override
	// TODO refactor to access a bufio.NewWriter interface
	// TODO need to be able to override file writing in unit tests

	for _, t := range templates {
		fileName := path.Join(dst, t.FileName)
		f, err := os.Create(fileName)
		if err != nil {
			return err
		}

		w := bufio.NewWriter(f)

		tmpl, err := template.New(t.FileName).Funcs(
			template.FuncMap{"StringsJoin": strings.Join}).Parse(t.GoTemplate)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(w, cluster); err != nil {
			return err
		}

		if err := w.Flush(); err != nil {
			return err
		}
		klog.Infof("Created terraform file: %s", t.FileName)

	}
	klog.Infof("Finished creating terraform files in: %s", dst)

	return nil
}
