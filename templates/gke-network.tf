/* TODO: delete all this when we have the first environment setup
  project_id                 = "${var.project_id}"
  name                       = "${var.project_name}-k8s"
  region                     = "${var.region}"
  zones                      = ["${split(",", var.zones)}"]
  network                    = "gke-vpc-01"
  subnetwork                 = "${var.region}-01"
  ip_range_pods              = "${var.region}-01-gke-01-pods"
  ip_range_services          = "${var.region}-01-gke-01-services"A

*/

resource "google_compute_network" "main" {
  name                    = "${var.project_name}-net"
  auto_create_subnetworks = "false"
}

//TODO - delete these after testing
//ip_range_pods              = "${var.region}-01-gke-01-pods"
//ip_range_services          = "${var.region}-01-gke-01-services"A

resource "google_compute_subnetwork" "main" {
  name          = "${var.region}-01"
  ip_cidr_range = "10.0.0.0/17"
  region        = "${var.region}"
  network       = "${google_compute_network.main.self_link}"

  secondary_ip_range {
    range_name    = "${var.region}-01-gke-01-pods"
    ip_cidr_range = "192.168.0.0/18"
  }

  secondary_ip_range {
    range_name    = "${var.region}-01-gke-01-services"
    ip_cidr_range = "192.168.64.0/18"
  }
}

