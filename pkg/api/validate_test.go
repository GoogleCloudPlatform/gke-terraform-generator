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

package api

import "testing"

func TestValidate(t *testing.T) {

	configFile := "testdata/example.yaml"
	gkeTF := parseYAML(t, configFile)

	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("failed %v", err)
	}

	if err := ValidateYamlInput(gkeTF); err != nil {
		t.Fatal(err)
	}
}

func TestValidateImageType(t *testing.T) {

	configFile := "testdata/example.yaml"
	gkeTF := parseYAML(t, configFile)

	nodes := *gkeTF.Spec.NodePools
	nodes[0].Spec.ImageType = "foo"

	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("failed %v", err)
	}

	if err := ValidateYamlInput(gkeTF); err == nil {
		t.Fatal("this should have failed")
	}

}

func TestCIDR(t *testing.T) {
	configFile := "testdata/example.yaml"
	gkeTF := parseYAML(t, configFile)

	gkeTF.Spec.Network.Spec.SubnetRange = "bad"
	gkeTF.Spec.Network.Spec.PodSubnetRange = "bad"
	gkeTF.Spec.Network.Spec.ServiceSubnetRange = "bad"

	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("failed %v", err)
	}

	if err := ValidateYamlInput(gkeTF); err == nil {
		t.Fatal("this should have failed")
	}

}

func parseYAML(t *testing.T, configFile string) *GkeTF {
	gkeTF, err := UnmarshalGkeTF(configFile)
	if gkeTF == nil {
		t.Fatal("gkeTf is nil")
	}
	gkeTF.Spec.ProjectId = "my project"
	if err != nil {
		t.Fatal(err)
	}
	return gkeTF
}
