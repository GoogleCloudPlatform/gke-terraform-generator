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
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

var configFile = "../../examples/example.yaml"

func TestValidate(t *testing.T) {

	gkeTF := parseYAML(t, configFile)

	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("failed %v", err)
	}

	if err := ValidateYamlInput(gkeTF); err != nil {
		t.Fatal(err)
	}
}

func TestCIDR(t *testing.T) {
	gkeTF := parseYAML(t, configFile)

	gkeTF.Spec.Network.Spec.SubnetRange = "bad"
	gkeTF.Spec.Network.Spec.PodSubnetRange = "bad"
	gkeTF.Spec.Network.Spec.ServiceSubnetRange = "bad"
	gkeTF.Spec.Network.Spec.MasterIPV4CIDRBlock = "bad"

	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("failed %v", err)
	}

	if err := ValidateYamlInput(gkeTF); err == nil {
		t.Fatal("this should have failed")
	}

}

func TestValidation(t *testing.T) {
	gkeTF := parseYAML(t, "../../examples/test-data.yaml")
	if err := SetApiDefaultValues(gkeTF, configFile); err != nil {
		t.Fatalf("failed %v", err)
	}

	if err := ValidateYamlInput(gkeTF); err != nil {
		t.Fatal(err)
	}

	domains := *gkeTF.Spec.StubDomains
	domains[0].DNSServerIPAddresses = []string{"foo", "bar"}
	gkeTF.Spec.StubDomains = &domains

	err := validator.New().Var(domains[0].DNSServerIPAddresses, "required,dive,ipv4")
	if err == nil {
		t.Fatalf("Validation should have failed, as foo is not a ipv4")
	}

	if err := ValidateYamlInput(gkeTF); err == nil {
		t.Fatal("this should have failed, not validating StubDomains correctly")
	}
}

func TestValidateIP(t *testing.T) {
	myEmail := "joeybloggs@gmail.com"

	errs := validator.New().Var(myEmail, "required,email")

	if errs != nil {
		t.Fatalf("err: %v", errs)
	}

	errs = validator.New().Var("foo", "required,ipv4")

	if errs == nil {
		t.Fatal("Validation should have failed, as foo is not a ipv4")
	}

	errs = validator.New().Var("73.95.44.185", "required,ipv4")

	if errs != nil {
		t.Fatalf("err: %v", errs)
	}

	ips := []string{"foo", "127.0.0.1", ""}

	errs = validator.New().Var(ips, "required,dive,ipv4")
	if errs == nil {
		t.Fatal("Validation should have failed, as foo is not a ipv4")
	}

	errs = validator.New().Var(&ips, "required,dive,ipv4")
	if errs == nil {
		t.Fatalf("Validation should have failed, as foo is not a ipv4 %v", errs)
	}

}

func parseYAML(t *testing.T, configFile string) *GkeTF {
	gkeTF, err := UnmarshalGkeTF(configFile)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	gkeTF.Spec.ProjectId = "my project"
	if err != nil {
		t.Fatal(err)
	}
	return gkeTF
}
