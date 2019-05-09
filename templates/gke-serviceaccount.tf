
// Create a default service account resource
resource "google_service_account" "gke_sa" {
  account_id   = "${var.service_account_id}"
  display_name = "${var.service_account_id}"
}

// using roles give it the needed set of permissiont to function
// Compute Viewer
resource "google_project_iam_binding" "gke_sa_computeView" {
  role = "roles/compute.viewer"

  members = [
    "serviceAccount:${google_service_account.gke_sa.email}",
  ]
}

//resource ="google_project_iam_binding" "gke_sa_networkAdmin" {
resource "google_project_iam_binding" "gke_sa_networkAdmin" {
  role = "roles/compute.networkAdmin"

  members = [
    "serviceAccount:${google_service_account.gke_sa.email}",
  ]
}

// Container Cluster Administrator 
resource "google_project_iam_binding" "gke_sa_clusterAdmin" {
  role = "roles/container.clusterAdmin"

  members = [
    "serviceAccount:${google_service_account.gke_sa.email}",
  ]
}

// Container Developer
resource "google_project_iam_binding" "gke_sa_containerDeveloper" {
  role = "roles/container.developer"

  members = [
    "serviceAccount:${google_service_account.gke_sa.email}",
  ]
}

// Service Account Amdministrator
resource "google_project_iam_binding" "gke_sa_saAccontAdmin" {
  role = "roles/iam.serviceAccountAdmin"

  members = [
    "serviceAccount:${google_service_account.gke_sa.email}",
  ]
}

// Service Account User
resource "google_project_iam_binding" "gke_sa_saUser" {
  role = "roles/iam.serviceAccountUser"

  members = [
    "serviceAccount:${google_service_account.gke_sa.email}",
  ]
}

// Project Iam Admin - (only required if service_account is set to create)
resource "google_project_iam_binding" "gke_sa_iamAdmin" {
  role = "roles/resourcemanager.projectIamAdmin"

  members = [
    "serviceAccount:${google_service_account.gke_sa.email}",
  ]
}

