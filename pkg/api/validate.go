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
	"k8s.io/klog"
)

// ValidateYamlInput checks the values that the user passes in via the yaml file.
func ValidateYamlInput(gkeTF *GkeTF) error {

	validate := validator.New()

	if err := validate.Struct(gkeTF.Spec); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok || validationErrors == nil {
			return err
		}
		klog.Errorf("error validating gke tf struct: %v", validationErrors)
		return validationErrors
	}

	// TODO we can remove this duplicated code with a refactor, for now it is ok

	if gkeTF.Spec.DatabaseEncryption != nil {
		if err := validate.Struct(gkeTF.Spec.DatabaseEncryption); err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if !ok || validationErrors == nil {
				return err
			}
			klog.Errorf("error validating gke tf databaseEncryption struct: %v", validationErrors)
			return validationErrors
		}
	}

	if gkeTF.Spec.MasterAuthorizedNetworksConfig != nil {
		if err := validate.Struct(gkeTF.Spec.MasterAuthorizedNetworksConfig); err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if !ok || validationErrors == nil {
				return err
			}
			klog.Errorf("error validating gke tf masterAuthorizedNetworksConfig struct: %v", validationErrors)
			return validationErrors
		}
	}

	if gkeTF.Spec.StubDomains != nil {
		if err := validate.Struct(gkeTF.Spec.StubDomains); err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if !ok || validationErrors == nil {
				return err
			}
			klog.Errorf("error validating gke tf stubDomains struct: %v", validationErrors)
			return validationErrors
		}
	}

	return nil
}
