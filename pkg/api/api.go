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
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// TODO go through the options and ensure we have them all
// TODO we are missing maintenance_start_time for instance
// master_ipv4_cidr_block

type GkeTF struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       ClusterSpec `yaml:"spec" validate:"required"`
}

// ClusterSpec API struct that represents a cluster.
type ClusterSpec struct {
	// The Name of the GKE Cluster
	// Name string `yaml:"name" validate:"required"`
	// The GCP ProjectId that is used to host the GKE cluster.
	ProjectId string `yaml:"projectId" validate:"required"`
	// Create a private GKE cluster.
	Private *bool `yaml:"private,omitempty" validate:"required" default:"true"`
	// Region is the GCP Region that is used for the network and if a regional cluster is created
	// it is used for that as well.
	Region string `yaml:"region" validate:"required"` // TODO validate that if it is not a regional cluster, then we need zones
	// Addons is struct that contains multiple bool flags that are used to deno which addons are to be installed.
	Addons *AddonsSpec `yaml:"addons" validate:"required"`
	// Network is a NetworkSpec struct that contains the details about the Network that will be created for the GKE Cluster.
	Network *GkeNetwork `yaml:"network" validate:"required,dive"`
	// Version is the base version for the cluster. This value defaults to 'latest'.
	// This value will be used for the GKE nodepools as well, unless a nodepool has a version.
	Version string `yaml:"version" default:"latest" validate:"required"`
	// Regional denotes if the GKE cluster will be created as a regional cluster.
	Regional *bool `yaml:"regional,omitempty" validate:"required" default:"true"`

	// RemoveDefaultNodePool enables the removal of the default GKE nodepool, which is the best practice.
	RemoveDefaultNodePool *bool `yaml:"removeDefaultNodePool,omitempty" default:"true"`

	// Zones are the GCP zones that the cluster runs inside of.
	Zones *[]string `yaml:"zones"`
	// Taints are a slice TaintSpec that model Kubernetes taints that are applied to all nodes.
	Taints *[]TaintSpec `yaml:"taints"`
	// OauthScopes is a slice of oauth scopes that are applied to a all nodes.  This slice defaults to the base required oauth
	// scopes.
	OauthScopes *[]string `yaml:"oauthScopes" default:"[\"https://www.googleapis.com/auth/trace.append\",\"https://www.googleapis.com/auth/service.management.readonly\",\"https://www.googleapis.com/auth/monitoring\",\"https://www.googleapis.com/auth/devstorage.read_only\",\"https://www.googleapis.com/auth/servicecontrol\"]"`
	// Tags is a slice of tags that are applied to all nodes.
	Tags *[]string `yaml:"tags"`

	// Labels is a map of labels that are applied to all node.  Labels are in the form of key and value strings.
	Labels *map[string]string `yaml:"labels"`
	// NodePools is a slice of NodePoolSpec struts that models a nodepool in GKE.
	NodePools *[]*GkeNodePool `yaml:"nodePools" validate:"required,dive"`
	// Metadata is a map of GCP compute instance metadata that will be applied to all compute instances.
	// This allows you to do things like start scripts.
	Metadata *map[string]string `yaml:"metadata"`
}

type GkeNetwork struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`

	Spec NetworkSpec `yaml:"spec" validate:"dive"`
}

// NetworkSpec API struct represents a network and subnet that is used for a GKE cluster.
type NetworkSpec struct {
	// SubnetName is the name of the GCP Subnet that is created.
	SubnetName string `yaml:"subnetName" validate:"required"`
	// SubnetRange is the base range for the GKE Nodes.
	SubnetRange string `yaml:"subnetRange" validate:"required,cidrv4"` // TODO calculate these based on nodes
	// PodSubnetRange is the ip aliased range used for GKE pods.
	PodSubnetRange string `yaml:"podSubnetRange" validate:"required,cidrv4"`
	// ServiceSubnetRange is the service subnet that is aliased.
	ServiceSubnetRange string `yaml:"serviceSubnetRange" validate:"required,cidrv4"`
}

type GkeNodePool struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	Spec       NodePoolSpec `yaml:"spec,omitempty" validate:"required,dive"`
}

type NodePoolSpec struct {
	MinCount    int16   `yaml:"minCount" default:"1" validate:"ltefield=MaxCount"`
	MaxCount    int16   `yaml:"maxCount" default:"1"`
	MachineType string  `yaml:"machineType" validate:"required"  default:"n1-standard-1"`
	AutoRepair  *bool   `yaml:"autoRepair,omitempty" default:"false"`
	AutoUpgrade *bool   `yaml:"autoUpgrade,omitempty" default:"false"`
	Preemptible *bool   `yaml:"preemptible" default:"false"`
	Version     *string `yaml:"version,omitempty"`
	DiskSizeGB  int     `yaml:"diskSizeGB" default:"100"`
	DiskType    string  `yaml:"diskType" default:"pd-ssd" validate:"eq=pd-ssd|eq=pd-standard"`
	ImageType   string  `yaml:"imageType" default:"cos" validate:"eq=cos|eq=ubuntu"`

	InitialNodeCount int16 `yaml:"initialNodeCount" validate:"required,ltefield=MaxCount" default:"1"`

	Tags        *[]string          `yaml:"tags" validate:"-"`
	OauthScopes *[]string          `yaml:"oauthScopes" validate:"omitempty,url"`
	Taints      *[]TaintSpec       `yaml:"taints"`
	Labels      *map[string]string `yaml:"labels"`
	Metadata    *map[string]string `yaml:"metadata"`

	AcceleratorType *string `yaml:"acceleratorType,omitempty"`
}

type TaintSpec struct {
	Key    string `yaml:"key" validate:"required"`
	Value  string `yaml:"value" validate:"required"`
	Effect string `yaml:"effect" validate:"required"`
}

type AddonsSpec struct {
	Istio              *bool `yaml:"istio,omitempty" default:"false"`
	Cloudrun           *bool `yaml:"cloudrun,omitempty" default:"false"`
	Logging            *bool `yaml:"logging,omitempty" default:"true"`
	Monitoring         *bool `yaml:"monitoring,omitempty" default:"true"`
	NetworkPolicies    *bool `yaml:"networkPolicies,omitempty" default:"true"`
	PodServicePolicies *bool `yaml:"podServicePolicies,omitempty" default:"true"`
	HPA                *bool `yaml:"hpa,omitempty" default:"true"`
	VPA                *bool `yaml:"vpa,omitempty" default:"false"`
	Autoscaling        *bool `yaml:"autoscaling,omitempty" default:"true"`
	BinaryAuth         *bool `yaml:"binaryAuth,omitempty" default:"false"`
	HTTPLoadBalancing  *bool `yaml:"httpLoadBalancing,omitempty" default:"true"`
}

type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	APIVersion string `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// users must create.
type ObjectMeta struct {
	// Name must be unique within a namespace. Is required when creating resources, although
	// some resources may allow a client to request the generation of an appropriate name
	// automatically. Name is primarily intended for creation idempotence and configuration
	// definition.
	Name string `yaml:"name,omitempty" protobuf:"bytes,1,opt,name=name" valdiate:"required"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `yaml:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `yaml:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
}

// unmarshalConfigurationFile reads the YAML file, configFile, and executes
// yaml.UnmarshalStrict, loading the value into spec.
func UnmarshalGkeTF(configFile string) (gkeTf *GkeTF, err error) {
	gkeTf = &GkeTF{}
	yamlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(yamlFile, gkeTf)
	if err != nil {
		return nil, err
	}

	return gkeTf, nil
}
