variable "project_id" {
  description = ""
  default = ""
}

variable "project_name" {
  description = ""
  default = ""
}

variable "region" {
  description = ""
  default = ""
}

variable "zones" {
  description = ""
  default = []
}

variable "service_account_id" {
  description = ""
  default = ""
}

variable "service_account" {
  description = ""
  default = "${var.service_account_id}@${var.project_id}.iam.gserviceaccount.com"
}
