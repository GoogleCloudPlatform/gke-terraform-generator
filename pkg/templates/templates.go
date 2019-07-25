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

// TODO package godoc
// TODO godocs in general

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"k8s.io/klog"

	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/api"
	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/terraform/cft"
	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/terraform/vanilla"
)

type TFType int

const (
	CFT     TFType = 0
	VANILLA TFType = 1
)

func (templateType TFType) String() string {
	names := [...]string{
		"CFT",
		"VANILLA",
	}
	if templateType < CFT || templateType > VANILLA {
		return "Unknown"
	}
	return names[templateType]
}

type TerraformTemplate struct {
	FileName   string
	GoTemplate string
}

type GKETemplates struct {
	Templates []*TerraformTemplate
}

func NewGKETemplates(tfType TFType) (*GKETemplates, error) {
	switch tfType {
	case CFT:
		return &GKETemplates{
			[]*TerraformTemplate{
				{
					"main.tf",
					cft.GKEMainTF,
				},
				{
					"network.tf",
					cft.GKENetworkTF,
				},
				{
					"outputs.tf",
					cft.GKEOutputsTF,
				},
				{
					"variables.tf",
					cft.GKEVariablesTF,
				},
			},
		}, nil
	case VANILLA:
		return &GKETemplates{
			[]*TerraformTemplate{
				{
					"main.tf",
					vanilla.GKEMainTF,
				},
				{
					"network.tf",
					vanilla.GKENetworkTF,
				},
				{
					"outputs.tf",
					vanilla.GKEOutputsTF,
				},
				{
					"variables.tf",
					vanilla.GKEVariablesTF,
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unable to find terraform type: %s", tfType)
	}
}

// CopyTo is used to copy all of the templates in the
// template directory to the given destination
func (gkeTemplates *GKETemplates) CopyTo(allowOverwrite bool, dst string, cluster *api.GkeTF) error {
	return gkeTemplates.processTemplates(allowOverwrite, dst, cluster)
}

func (gkeTemplates *GKETemplates) processTemplates(allowOverwrite bool, dst string, cluster *api.GkeTF) error {
	// TODO refactor to access a bufio.NewWriter interface
	// TODO need to be able to override file writing in unit tests

	for _, t := range gkeTemplates.Templates {
		fileName := path.Join(dst, t.FileName)
		if !allowOverwrite {
			if f, err := os.Open(fileName); f != nil && err == nil {
				return fmt.Errorf("file already exists and overwrites not allowed file: %s", fileName)
			}
		}
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
