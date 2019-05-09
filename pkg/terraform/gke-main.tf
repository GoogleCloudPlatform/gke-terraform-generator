/**
 * Copyright 2018 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

provider "google" {
  version = "2.7.0"
  project = "${var.project_id}"
  region  = "${var.region}"
}

provider "google-beta" {
  version = "2.7.0"
  project = "${var.project_id}"
  region  = "${var.region}"
}

// TODO: - setup add capability to use remote state

module "gke" {
{{- if .Spec.Private }}
  source = "terraform-google-modules/kubernetes-engine/google"
{{- else}}
  source = "terraform-google-modules/kubernetes-engine/google//modules/private-cluster"
  enable_private_endpoint    = "true"
  enable_private_nodes       = "true"
  // TODO make this configurable
  master_ipv4_cidr_block     = "10.0.0.0/28"
{{- end}}

  project_id = "${var.project_id}"
  name       = "${var.cluster_name}"
  region   = "${var.region}"
{{- if .Spec.Zones }}
  zones   = "${var.zones}" // FIXME we may need to convert a list to a string here
{{- end}}
  regional   = {{.Spec.Regional}}
  kubernetes_version    = "{{.Spec.Version}}"

  network           = "${module.gke-network.network_name}"
  subnetwork        = "${module.gke-network.subnets_names[0]}"
  ip_range_pods     = "${module.gke-network.network_name}-${var.cluster_name}-pod-range"
  ip_range_services = "${module.gke-network.network_name}-${var.cluster_name}-service-range"

  /* dashboard is being depricated, so do not install it */
  kubernetes_dashboard        = "false"
  http_load_balancing         = "{{.Spec.Addons.HTTPLoadBalancing}}"
  network_policy              = "{{.Spec.Addons.NetworkPolicies}}"
  horizontal_pod_autoscaling  = "{{.Spec.Addons.HPA}}"
  // TODO add psp, binary auth, cloudrun, istio
  // TODO double check I am not missing anything else
  // enable_binary_authorization = "{{.Spec.Addons.BinaryAuth}}"
  // istio = "{{.Spec.Addons.Istio}}"
  // cloudrun = "{{.Spec.Addons.Cloudrun}}"

  service_account          = "create"
  remove_default_node_pool = "{{.Spec.RemoveDefaultNodePool}}"

  disable_legacy_metadata_endpoints = "true"

  // TODO need version of gke cfp module

  // TODO capability to build empty nodepool
  node_pools = [
{{- range .Spec.NodePools}}
    {
      name               = "{{.Name}}"
      machine_type       = "{{.Spec.MachineType}}"
      {{- if .Spec.AcceleratorType}}
      accelerator_type  = "{{.Spec.AcceleratorType}}"
      {{ end }}
      min_count          = {{.Spec.MinCount}}
      max_count          = {{.Spec.MaxCount}}
      disk_size_gb       = {{.Spec.DiskSizeGB}}
      disk_type          = "{{.Spec.DiskType}}"
      image_type         = "{{.Spec.ImageType}}"
      auto_repair        = {{.Spec.AutoRepair}}
      auto_upgrade       = {{.Spec.AutoUpgrade}}
      preemptible        = {{.Spec.Preemptible}}
      initial_node_count = {{.Spec.InitialNodeCount}}
    },{{end}}
  ]

  node_pools_oauth_scopes = {
    all = [
      {{- if .Spec.OauthScopes}}
        {{- range .Spec.OauthScopes}}
        "{{.}}",{{end}}{{end}}
       ]
  {{- range .Spec.NodePools}}

    {{.Name}} = [
    {{- if .Spec.OauthScopes}}
      {{- range .Spec.OauthScopes}}
      {{.}},{{end}}{{end}}
    ]{{end}}
  }

  node_pools_labels = {

    all = {
      {{if .Spec.Labels}}
        {{range $key, $value := .Spec.Labels}}
          {{ $key }} = "{{ $value }}"
        {{end}}
      {{end}}
    }

    {{range .Spec.NodePools}}

    {{.ObjectMeta.Name}} = {
      {{if .Spec.Labels}}
        {{ range $key, $value := .Spec.Labels }}
        {{ $key }} = "{{ $value }}"
        {{end}}
      {{end}}
    }{{end}}
  }

  node_pools_metadata = {
    all = {
    {{- if .Spec.Metadata}}
      {{- range $key, $value := .Metadata }}
        $key = "$value"
      {{- end}}
    {{end -}}
    }
    {{range .Spec.NodePools}}
    {{.Name}} = {
    {{- if .Spec.Metadata}}
      {{- range $key, $value := .Spec.Metadata }}
      $key = "$value"
      {{- end}}
      {{end -}}
    }
    {{end}}
  }

  node_pools_tags = {
    all = [
    {{- if .Spec.Tags}}
      {{- range .Spec.Tags}}
      "{{ .}}",
      {{- end}}
    {{end -}}
    ]
  {{range .Spec.NodePools}}
    {{.ObjectMeta.Name}} = [
    {{- if .Spec.Tags}}
      {{- range .Spec.Tags}}
      "{{.}}",
      {{- end}}
    {{end -}}
    ]
  {{end}}
  }

  node_pools_taints = {
    all = [
    {{- if .Spec.Taints }}
      {{- range .Spec.Taints }}
        {
          key = "{{.Key}}"
          value = "{{.Value}}"
          effect = "{{.Effect}}"
        },
      {{- end}}
     {{end -}}
    ]
    {{ range .Spec.NodePools}}
    {{.ObjectMeta.Name}} = [
      {{- if .Spec.Taints}}
      {{- range .Spec.Taints}}
      {
        key    = "{{.Key}}"
        value  = "{{.Value}}"
        effect = "{{.Effect}}"
      },
      {{- end}}
      {{end -}}
    ]
    {{end}}
  }
}
