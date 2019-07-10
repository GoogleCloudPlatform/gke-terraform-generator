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

import (
	"testing"
)

// TODO add more tests

func TestDefaults(t *testing.T) {

	var configFile = "../../examples/example.yaml"
	gkeTF, err := UnmarshalGkeTF(configFile)

	if err != nil {
		t.Fatal(err)
	}

	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("error merging defaults: %v", gkeTF)
	}

	if gkeTF.Spec.Private == "false" {
		t.Fatal("gkeTF.Spec.Private is not set to true and it should be")
	}

	t.Logf("*gkeTF.Spec.Addons.ClusterAutoscaling: %v", *gkeTF.Spec.Addons.ClusterAutoscaling)
	if *gkeTF.Spec.Addons.ClusterAutoscaling != false {
		t.Fatal("*gkeTF.Spec.Addons.ClusterAutoscaling is not false")
	}

	if gkeTF.Spec.OauthScopes == nil {
		t.Fatal("*gkeTF.Spec.OauthScopes is nil")
	}

}

func TestDefaultsSmall(t *testing.T) {

	var configFile = "../../examples/min-example.yaml"
	gkeTF, err := UnmarshalGkeTF(configFile)

	if err != nil {
		t.Fatal(err)
	}

	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("error merging defaults: %v", gkeTF)
	}

	if gkeTF.Spec.OauthScopes == nil {
		t.Fatal("*gkeTF.Spec.OauthScopes is nil")
	}

	if len(*gkeTF.Spec.OauthScopes) == 0 {
		t.Fatal("*gkeTF.Spec.OauthScopes is empty")
	}

}
