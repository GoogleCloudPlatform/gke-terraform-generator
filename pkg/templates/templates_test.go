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
	"partner-code.googlesource.com/gke-terraform-generator/pkg/api"
	"testing"
)

func TestTemplates(t *testing.T) {

	gkeTF, err := api.UnmarshalGkeTF("../api/testdata/example.yaml")

	if err != nil {
		t.Fatal(err)
	}

	if gkeTF == nil {
		t.Fatal("unable to load file")
	}

	if gkeTF.Name == "" {
		t.Fatal("gkeTF.Name is empty")
	}

	testTemplates := NewGKETemplates()
	err = testTemplates.CopyTo(".", gkeTF)

	if err != nil {
		t.Fatal(err)
	}

	_, err = ioutil.ReadFile("main.tf")
	if err != nil {
		t.Fatal(err)
	}

	/*
	tf := "../../../terraform/terraform"

	_, err = os.Stat(tf)
	if err != nil {
		t.Fatal("terraform doesn't exist")
	}

	cmd := exec.Command("testdata/tf_wrapper.sh", tf)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("terraform init failed with %s\n", err)
	}

	out, err := exec.Command(tf, "init").Output()
	if err != nil {
		t.Logf("terraform plan %s\n", out)
		t.Fatalf("terraform plan failed with %s\n", err)
	}
	 */
}

