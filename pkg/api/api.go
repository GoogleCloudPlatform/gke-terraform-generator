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

// TODO add nvida card values

// TODO add bastion
// TODO add nat gateway

// GkeTF is the base layer for the API.  It includes a ClusterSpec and other obligatory information.
type GkeTF struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`
	// ClusterSpec include the base information modeling a GKE Cluster.
	Spec ClusterSpec `yaml:"spec" validate:"required,dive"`
}

// ClusterSpec API struct that represents a cluster.
type ClusterSpec struct {
	// The GCP ProjectId that is used to host the GKE cluster.
	ProjectId string `yaml:"projectId" validate:"required"`
	// Create a private GKE cluster.
	Private string `yaml:"private,omitempty" default:"true" validate:"eq=true|eq=false"`
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
	Taints *[]TaintSpec `yaml:"taints" validate:"omitempty,dive"`
	// OauthScopes is a slice of oauth scopes that are applied to a all nodes.  This slice defaults to the base required oauth
	// scopes.
	OauthScopes *[]string `yaml:"oauthScopes" default:"[\"https://www.googleapis.com/auth/trace.append\",\"https://www.googleapis.com/auth/service.management.readonly\",\"https://www.googleapis.com/auth/monitoring\",\"https://www.googleapis.com/auth/devstorage.read_only\",\"https://www.googleapis.com/auth/servicecontrol\"]"`
	// Tags is a slice of tags that are applied to all nodes.
	Tags *[]string `yaml:"tags"`

	// Labels is a map of labels that are applied to all node.  Labels are in the form of key and value strings.
	Labels *map[string]string `yaml:"labels" validate:"omitempty,dive"`
	// NodePools is a slice of NodePoolSpec struts that models a nodepool in GKE.
	NodePools *[]*GkeNodePool `yaml:"nodePools" validate:"required,dive"`
	// Metadata is a map of GCP compute instance metadata that will be applied to all compute instances.
	// This allows you to do things like start scripts.
	Metadata *map[string]string `yaml:"metadata" validate:"omitempty,dive"`

	// TODO test below values

	// Description of the cluster.
	Description *string `yaml:"Description"`

	IpMasqLinkLocal     *string `yaml:"ipMasqLinkLocal"`
	IpMasqRsyncInterval *string `yaml:"ipMasqRsyncInterval"`
	// RFC3339 format
	MaintenanceStartTime *string `yaml:"maintenanceStartTime"`

	IssueClientCertificate *string `yaml:"IssueClientCertificate" default:"false"`

	// MasterAuthorizedNetworksConfig is a slice of the desired configuration options for master authorized networks.
	// Omit the nested cidr_blocks attribute to disallow external access (except the cluster node IPs, which GKE automatically whitelists)
	MasterAuthorizedNetworksConfig *[]MasterAuthorizedNetworksConfigSpec `yaml:"masterAuthorizedNetworksConfig" validate:"omitempty,dive"`

	// StubDomains and their resolvers to forward DNS queries for a certain domain to an external DNS server.
	StubDomains *[]StubDomainsSpec `yaml:"stubDomains" validate:"omitempty,dive"`

	// DatabaseEncryption allows a user to configure etcd database encyptions using a provided KMS key name.
	DatabaseEncryption *DatabaseEncryptionSpec `yaml:"databaseEncryption" validate:"omitempty,dive"`

	// NodeVersion is the default kubernetes version of nodes in the node pools.
	NodeVersion *string `yaml:"nodeVersion"`

	// ServiceAccount is the name of the service account to use for the cluster, or can hold the value "create".
	// This value defaults to the value "create", which creates a new service account for the cluster.
	ServiceAccount *string `yaml:"serviceAccount" default:"create"`

	// TODO check if we have this
	DeployUsingPrivateEndpoint *bool `yaml:"deployUsingPrivateEndpoint"`
}

// GkeNetwork wraps a NetworkSpec.
type GkeNetwork struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`

	// ClusterSpec include the base information modeling a Network for a GKE Cluster.
	Spec NetworkSpec `yaml:"spec" validate:"required,dive"`
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
	// The IP range in CIDR notation to use for the hosted master network
	MasterIPV4CIDRBlock string `yaml:"masterIPV4CIDRBlock" validate:"cidrv4"`

	// TODO test for this - The given master_ipv4_cidr 10.0.0.0/28 overlaps with an existing network 10.0.0.0/24.
	// TODO we make have to make MasterIPV4CIDRBlock required if it is a private cluster.  Need more testing.

}

// GkeNetwork wraps a NodePoolSpec.
type GkeNodePool struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata,omitempty"`

	// ClusterSpec include the base information modeling a NodePool for a GKE Cluster.
	Spec NodePoolSpec `yaml:"spec,omitempty" validate:"required,dive"`
}

// MasterAuthorizedNetworksConfigSpec models the desired configuration options for master authorized networks.
type MasterAuthorizedNetworksConfigSpec struct {
	// CidrBlock stores the CIDR
	CidrBlock *string `yaml:"cidrBlock" validate:"cidrv4"`
	// DisplayName is the display_name in terraform.
	DisplayName *string `yaml:"displayName" validate:"required"`
}

type DatabaseEncryptionSpec struct {
	// State can be two different values "ENCRYPTED", "DECRYPTED".
	State *string `yaml:"state" validate:"required,eq=ENCRYPTED|eq=DECRYPTED"`
	// Keyname is the name of the KMS key.
	KeyName *string `yaml:"keyName" validate:"required"`
}

type StubDomainsSpec struct {
	TypeMeta             `yaml:",inline"`
	ObjectMeta           `yaml:"metadata,omitempty"`
	DNSServerIPAddresses []string `yaml:"dnsServerIPAddresses" validate:"required,dive,ipv4"`
}

// NodePoolSpec API struct that represents a GKE Nodepool.
type NodePoolSpec struct {

	// MinCount of the nodepool.
	// This value defaults to 1 and must be less than MaxCount
	MinCount int16 `yaml:"minCount" default:"1" validate:"ltefield=MaxCount"`
	// MaxCount of the nodepool.
	// This value defaults to 1.
	MaxCount int16 `yaml:"maxCount" default:"1"`
	// MachineType of the nodepool, which defaults to a n1-standard-1. See
	// https://cloud.google.com/compute/docs/machine-types for more information about
	// GCP machine types.
	MachineType string `yaml:"machineType" default:"n1-standard-1"`
	// AutoRepair enables GKE's node auto-repair feature that helps keeping the nodes in your
	// cluster in a healthy, running state.
	// See https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-repair.
	// This feature defaults to true, by default, and is enabled.
	AutoRepair *bool `yaml:"autoRepair,omitempty" default:"true"`
	// AutoUpgrade enables Node auto-upgrades help you keep the nodes in your cluster up to date
	// with the cluster master version when your master is updated on your behalf.
	// See https://cloud.google.com/kubernetes-engine/docs/how-to/node-auto-upgrades.
	// This feature defaults to false, and not enabled.
	AutoUpgrade *bool `yaml:"autoUpgrade,omitempty" default:"false"`
	// Preemptible causes the nodepool create with Preemptible VMs which are Google
	// Compute Engine VM instances that last a maximum of 24 hours and
	// provide no availability guarantees.
	// See https://cloud.google.com/kubernetes-engine/docs/how-to/preemptible-vms.
	// Preemptible defaults to false.
	Preemptible *bool `yaml:"preemptible" default:"false"`
	// Version is the GKE version for the nodepool.. This value defaults to 'latest'.
	// This value will override the version in the parent struct.
	Version *string `yaml:"version,omitempty"`
	// DiskSizeGB is the node disk size.
	// This value defaults to 100.
	DiskSizeGB int `yaml:"diskSizeGB" default:"100"` // TODO validate that this is a positive number
	// DiskType is the node disk type.
	// Values can be pd-ssd or pd-standard, and it defaults to pd-ssd.
	DiskType string `yaml:"diskType" default:"pd-ssd" validate:"eq=pd-ssd|eq=pd-standard"`
	// ImageType is the node operating system.
	// See https://cloud.google.com/kubernetes-engine/docs/concepts/node-images.
	// Values can be cos, cos_containerd or ubuntu, and it defaults to cos.
	ImageType string `yaml:"imageType" default:"cos" validate:"eq=cos|eq=ubuntu|eq=cos_containerd"`

	// InitialNodeCount is the number of nodes created at inception of the cluster.
	// This value defaults to 1.
	InitialNodeCount int16 `yaml:"initialNodeCount" validate:"required,ltefield=MaxCount" default:"1"`

	// Tags slice containing node network tags for this specific nodepool.
	// See https://cloud.google.com/vpc/docs/add-remove-network-tags.
	Tags *[]string `yaml:"tags"` // validate:"omitempty,dive"` // this will cause a panic with validator
	// OauthScopes is a slice of Oauth Scope URLs that are applied to
	// the GCP instances in a nodepool.
	// See https://developers.google.com/identity/protocols/googlescopes
	OauthScopes *[]string `yaml:"oauthScopes"` // validate:"omitempty,dive,url"`
	// Taints are a slice of TaintSpec structs.
	// See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/ and
	// https://cloud.google.com/kubernetes-engine/docs/how-to/node-taints.
	Taints *[]TaintSpec `yaml:"taints"` // validate:"omitempty,dive"`
	// Labels is a map of GCP instance labels.
	// See https://cloud.google.com/compute/docs/labeling-resources.
	Labels   *map[string]string `yaml:"labels"`
	Metadata *map[string]string `yaml:"metadata"`

	AcceleratorType *string `yaml:"acceleratorType,omitempty"`
	ServiceAccount  *string `yaml:"serviceAccount"`
}

// Todo I need to look at the API to figure out how to validate TaintSpec

// TaintSpec models a Kubernetes Node Taint.
//
// See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
// and
// https://cloud.google.com/kubernetes-engine/docs/how-to/node-taints.
//
// For example:
//
//	tolerations:
//	  - key: "key"
//	  operator: "Equal"
//	  value: "value"
//	  effect: "NoSchedule"
//
type TaintSpec struct {

	// TODO need validation

	// Key is the key value in a taint.
	Key string `yaml:"key"`
	// Value is the value field in a taint.
	Value string `yaml:"value"`
	// Effect is the effect field in a taint.
	Effect string `yaml:"effect"`
}

// AddonsSpec is struct that contains multiple bool flags that are used to denote which addons are to be installed.
// See the following URLS for more information about the various addons.
//
// - https://cloud.google.com/istio/docs/istio-on-gke/overview
//
// - https://cloud.google.com/run/docs/gke/setup
//
// - https://cloud.google.com/monitoring/kubernetes-engine/
//
// - https://cloud.google.com/kubernetes-engine/docs/how-to/network-policy
//
// - https://cloud.google.com/kubernetes-engine/docs/how-to/scaling-apps
//
// - https://cloud.google.com/kubernetes-engine/docs/concepts/verticalpodautoscaler
//
// - https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-autoscaler
//
// - https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies
//
// - https://cloud.google.com/binary-authorization/docs/
type AddonsSpec struct {
	// Istio installs managed istio on the cluster.
	// Default value for Istio is false.
	// See: https://cloud.google.com/istio/docs/istio-on-gke/overview
	Istio *bool `yaml:"istio,omitempty" default:"false"`
	// Cloudrun installs managed Cloudrun on the cluster.
	// Default value for Cloudrun is false.
	//
	// See  https://cloud.google.com/run/docs/gke/setup.
	Cloudrun *bool `yaml:"cloudrun,omitempty" default:"false"`
	// Logging enables stack driver logging for the cluster.
	// Default value for Logging is true.
	// Automatically send logs from the cluster to the Google Cloud Logging
	// API.
	// include logging.googleapis.com, logging.googleapis.com/kubernetes
	Logging *string `yaml:"logging,omitempty" default:"logging.googleapis.com/kubernetes"` // TODO validate
	// Monitoring enables stack driver logging for the cluster.
	// Default value for Monitoring is true.
	//
	// See https://cloud.google.com/monitoring/kubernetes-engine/.
	//
	// Automatically send metrics from pods in the cluster to the Google Cloud
	// Monitoring API. VM metrics will be collected by Google Compute Engine
	// regardless of this setting.
	// monitoring.googleapis.com, monitoring.googleapis.com/kubernetes
	Monitoring *string `yaml:"monitoring,omitempty" default:"monitoring.googleapis.com/kubernetes"` // TODO validate
	// NetworkPolicy enables network policy for the cluster.
	// Enable network policy enforcement for this cluster.
	// Default value for NetworkPolicy is true.
	//
	// See https://cloud.google.com/kubernetes-engine/docs/how-to/network-policy.
	NetworkPolicy *bool `yaml:"networkPolicy,omitempty" default:"true"`
	// HPA enables horizontal pod autoscaling for the cluster.
	// Default value for HPA is false.
	//
	// See https://cloud.google.com/kubernetes-engine/docs/how-to/scaling-apps.
	HPA *bool `yaml:"hpa,omitempty" default:"true"`
	// VPA enables vertical pod autoscaling for the cluster.
	// Default value for HPA is false.
	//
	// See https://cloud.google.com/kubernetes-engine/docs/concepts/verticalpodautoscaler.
	VPA *bool `yaml:"vpa,omitempty" default:"false"`
	// Autoscaling enables cluster nodepool autoscaling.
	// Default value Autoscaling is true.
	//
	// See https://cloud.google.com/kubernetes-engine/docs/concepts/cluster-autoscaler
	Autoscaling *bool `yaml:"autoscaling,omitempty" default:"true"`
	// BinaryAuth enables binary authorization for the cluster.
	// Default value BinaryAuth is true.
	//
	// See https://cloud.google.com/binary-authorization/docs/.
	BinaryAuth *bool `yaml:"binaryAuth,omitempty" default:"true"`

	// HTTPLoadBalancing enables HTTP Load Balancing for the cluster.
	// Default value HTTPLoadBalancing is true.
	// TODO figure out what this actually is, cannot find it in gcloud
	HTTPLoadBalancing *bool `yaml:"httpLoadBalancing,omitempty" default:"true"`
	// PodSecurityPolicy enables Pod Security Policy for the Cluster.
	// Default value PodSecurityPolicy is false.
	// Enables the pod security policy admission controller for the cluster.
	// The pod security policy admission controller adds fine-grained pod
	// create and update authorization controls through the PodSecurityPolicy
	//
	// API objects. For more information, see
	//
	// https://cloud.google.com/kubernetes-engine/docs/how-to/pod-security-policies.
	PodSecurityPolicy *bool `yaml:"podSecurityPolicy,omitempty" default:"false"`
}

// TypeMeta is metadata that all resources must have, which includes all objects
// users must create.
type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `yaml:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	APIVersion string `yaml:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

// ObjectMeta is metadata that all resources must have, which includes all objects
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
