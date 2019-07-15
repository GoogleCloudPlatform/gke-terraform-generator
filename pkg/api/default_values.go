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
	"github.com/creasty/defaults"
	"github.com/imdario/mergo"
	"k8s.io/klog"
)

// SetApiDefaultValues creates default ClusterSpec and its children and merges the default values into
// the clusterSpec values parsed from the user provided YAML.
// This func allows for gen to provide default values as set in api.ClusterSpec definition.
func SetApiDefaultValues(gkeTF *GkeTF, configFile string) error {

	// The defaults library does not support bool pointers.
	// So we have to make a copy and then copy the original value back.
	original, err := UnmarshalGkeTF(configFile)
	if err != nil {
		klog.Errorf("error creating copy of gke api: %v", err)
		return err
	}

	defaultSpec := &GkeTF{
		Spec: ClusterSpec{
			Addons: &AddonsSpec{},
		},
	}

	defaultNodePool := &GkeNodePool{
		Spec: NodePoolSpec{},
	}

	if err := defaults.Set(defaultSpec); err != nil {
		klog.Errorf("error setting defaults for defaultSepc: %v", err)
		return err
	}
	if err := defaults.Set(&defaultNodePool.Spec); err != nil {
		klog.Errorf("error setting defaults for defaultNodePool.Spec: %v", err)
		return err
	}
	if err := defaults.Set(defaultSpec.Spec.Addons); err != nil {
		klog.Errorf("error setting defaults for defaultSpec.Spec.Addons: %v", err)
		return err
	}

	if gkeTF.Spec.Addons == nil {
		gkeTF.Spec.Addons = &AddonsSpec{}
	}

	if err := mergo.Merge(&gkeTF.Spec, &defaultSpec.Spec); err != nil {
		klog.Errorf("error gkeTF: %v", err)
		return err
	}

	originalNodePools := *original.Spec.NodePools
	for i, nodePool := range *gkeTF.Spec.NodePools {
		if err := mergo.Merge(&nodePool.Spec, &defaultNodePool.Spec); err != nil {
			klog.Errorf("error merging nodePoolSpec: %v", err)
			return err
		}
		originalNodePool := originalNodePools[i]
		if originalNodePool.Spec.AutoRepair != nil {
			nodePool.Spec.AutoRepair = originalNodePool.Spec.AutoRepair
		}

		if originalNodePool.Spec.AutoUpgrade != nil {
			nodePool.Spec.AutoUpgrade = originalNodePool.Spec.AutoUpgrade
		}

	}

	// Go through and reset values overwritten by defaults
	if original.Spec.RemoveDefaultNodePool != nil {
		*gkeTF.Spec.RemoveDefaultNodePool = *original.Spec.RemoveDefaultNodePool
	}
	if original.Spec.Addons != nil {
		if original.Spec.Addons.VPA != nil {
			*gkeTF.Spec.Addons.VPA = *original.Spec.Addons.VPA
		}
		if original.Spec.Addons.PodSecurityPolicy != nil {
			*gkeTF.Spec.Addons.PodSecurityPolicy = *original.Spec.Addons.PodSecurityPolicy
		}
		if original.Spec.Addons.HTTPLoadBalancing != nil {
			*gkeTF.Spec.Addons.HTTPLoadBalancing = *original.Spec.Addons.HTTPLoadBalancing
		}
		if original.Spec.Addons.BinaryAuth != nil {
			*gkeTF.Spec.Addons.BinaryAuth = *original.Spec.Addons.BinaryAuth
		}
		if original.Spec.Addons.ClusterAutoscaling != nil {
			*gkeTF.Spec.Addons.ClusterAutoscaling = *original.Spec.Addons.ClusterAutoscaling
		}
	}

	return nil
}
