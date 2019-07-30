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
	"io/ioutil"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/api"
	"github.com/GoogleCloudPlatform/gke-terraform-generator/pkg/terraform/cft"
)

// TODO implement https://github.com/hashicorp/hcl/blob/master/decoder_test.go

func TestTemplates(t *testing.T) {

	configFile := "../../examples/public-example.yaml"
	gkeTF, err := api.UnmarshalGkeTF(configFile)

	if err != nil {
		t.Fatal(err)
	}

	if gkeTF == nil {
		t.Fatal("unable to load file")
	}

	if gkeTF.Name == "" {
		t.Fatal("gkeTF.Name is empty")
	}

	if gkeTF.Spec.Private == "true" {
		t.Fatal("gkeTF.Spec.Private should be false")
	}

	if err := api.SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("error merging defaults: %v", gkeTF)
	}

	testTemplates, err := NewGKETemplates(CFT)
	if err != nil {
		t.Fatal(err)
	}

	err = testTemplates.CopyTo(false, ".", gkeTF)

	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile("main.tf")
	if err != nil {
		t.Fatal(err)
	}

	s := string(b)

	public := "source = \"terraform-google-modules/kubernetes-engine/google\""
	if !strings.Contains(s, public) {

		t.Log(s)
		t.Log(cft.GKEMainTF)
		t.Fatalf("template does not contain the correct source provider")
	}
}

func TestPrivateTemplate(t *testing.T) {

	configFile := "../../examples/example.yaml"
	gkeTF, err := api.UnmarshalGkeTF(configFile)

	if err != nil {
		t.Fatal(err)
	}

	if gkeTF == nil {
		t.Fatal("unable to load file")
	}

	if gkeTF.Name == "" {
		t.Fatal("gkeTF.Name is empty")
	}

	if gkeTF.Spec.Private == "false" {
		t.Fatal("gkeTF.Spec.Private should be true")
	}

	if err := api.SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("error merging defaults: %v", gkeTF)
	}

	testTemplates, err := NewGKETemplates(CFT)
	if err != nil {
		t.Fatal(err)
	}

	err = testTemplates.CopyTo(true, ".", gkeTF)

	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile("main.tf")
	if err != nil {
		t.Fatal(err)
	}

	s := string(b)

	private := "source = \"terraform-google-modules/kubernetes-engine/google//modules/private-cluster\""
	if !strings.Contains(s, private) {
		t.Log(s)
		t.Log(cft.GKEMainTF)
		t.Fatalf("template does not contain the private source provider")
	}
}

// TODO test all yaml files
func TestFullTemplate(t *testing.T) {

	configFile := "../../examples/full.yaml"
	gkeTF, err := api.UnmarshalGkeTF(configFile)

	if err != nil {
		t.Fatal(err)
	}

	if gkeTF == nil {
		t.Fatal("unable to load file")
	}

	if gkeTF.Name == "" {
		t.Fatal("gkeTF.Name is empty")
	}

	if gkeTF.Spec.Private == "false" {
		t.Fatal("gkeTF.Spec.Private should be true")
	}

	if err := api.SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("error merging defaults: %v", gkeTF)
	}

	testTemplates, err := NewGKETemplates(VANILLA)
	if err != nil {
		t.Fatal(err)
	}

	err = testTemplates.CopyTo(true, ".", gkeTF)

	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile("network.tf")
	if err != nil {
		t.Fatal(err)
	}

	s := string(b)

	bastionzone := "bastion_zone = \"us-east4-c\""
	if !strings.Contains(s, bastionzone) {
		t.Log(s)
		t.Fatalf("template does not contain the defined bastion zone")
	}

	zone := "zone = local.bastion_zone"
	if !strings.Contains(s, zone) {
		t.Log(s)
		t.Fatalf("template does not contain the defined instnace zone for the bastion")
	}
}
