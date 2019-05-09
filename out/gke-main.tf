/*
TODO: - enable debug setup of the state document, its creating havoc
// gs://terraform-gke-tf/state/dev/gke-tf/
terraform {
  backend "gcs" {
    project = "${var.project}"
    bucket = "terraform-gke-tf"
    prefix = "state/dev/gke-tf"
  }
}
*/

provider "google" {
  version = "2.3.0"
  project = "${var.project_name}"
  region  = "${var.region}"
}

provider "google-beta" {
  version = "2.3.0"
  project = "${var.project_name}"
  region  = "${var.region}"
}

// useful for creating unique collision free (mostly) network names
resource  "random_string" "suffix" {
  length  = 4
  special = false
  upper   = false
}

/*
TODO: - debug setup of the state document, its creating havoc
// gs://terraform-gke-tf/state/dev/gke-tf/
terraform {
  backend "gcs" {
    project = "cli-gke-terraform"
    bucket = "terraform-gke-tf"
    prefix = "state/dev/gke-tf"
  }
}*/

module "gke" {
  source = "terraform-google-modules/kubernetes-engine/google"

  /* TODO - for the gke-cli project we will waant to enable private clusters, that means then that even gke-main.tf itself becomes a template
{% if private_cluster %}//modules/private-cluster{% endif %}"
*/

  project_id                 = "${var.project_id}"
  name                       = "${var.project_name}-k8s"
  region                     = "${var.region}"
  zones                      = "${var.zones}"
  network                    = "${google_compute_network.main.name}"
  subnetwork                 = "${google_compute_subnetwork.main.name}"
  ip_range_pods              = "${google_compute_subnetwork.main.secondary_ip_range.0.range_name}"
  ip_range_services          = "${google_compute_subnetwork.main.secondary_ip_range.1.range_name}"
  http_load_balancing        = false
  horizontal_pod_autoscaling = true
  kubernetes_dashboard       = true
  network_policy             = true

  /*
 TODO - when creating a private cluster add these values
  enable_private_endpoint    = true
  enable_private_nodes       = true
  master_ipv4_cidr_block     = "10.0.0.0/28"
*/

  node_pools = [
    {
      name               = "default-node-pool"
      machine_type       = "n1-standard-2"
      min_count          = 1
      max_count          = 5 // TODO: we have on opinion let's enable it via command line
      disk_size_gb       = 100
      disk_type          = "pd-standard"
      image_type         = "COS"
      auto_repair        = true
      auto_upgrade       = true
      service_account    = "${var.service_account}"
      preemptible        = false
      initial_node_count = 3 // TODO: we have an opinion let's enable it via command line
    },
  ]
  node_pools_oauth_scopes = {
    all = []

    default-node-pool = [
      "https://www.googleapis.com/auth/cloud-platform",
    ]
  }
  node_pools_labels = {
    all = {}

    default-node-pool = {
      default-node-pool = "true"
    }
  }
  node_pools_metadata = {
    all = {}

    default-node-pool = {
      node-pool-metadata-custom-value = "my-node-pool"
    }
  }
  node_pools_taints = {
    all = []

    default-node-pool = [
      {
        key    = "default-node-pool"
        value  = "true"
        effect = "PREFER_NO_SCHEDULE"
      },
    ]
  }
  node_pools_tags = {
    all = []

    default-node-pool = [
      "default-node-pool",
    ]
  }
}
