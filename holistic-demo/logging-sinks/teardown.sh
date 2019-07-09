#!/usr/bin/env bash

# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# "---------------------------------------------------------"
# "-                                                       -"
# "-  deletes the resouses created                         -"
# "-                                                       -"
# "---------------------------------------------------------"

# Do not set errexit as it makes partial deletes impossible
# If it is set, the script would exit if any statement returns a non-true return value
# that means it won't destory other resources if the flow continues.
set -o nounset
set -o pipefail

# We have to delete the dataset before the Terraform
# Otherwise we will run into the following error
# "google_bigquery_dataset.gke-bigquery-dataset: googleapi:
# Error 400: Dataset pso-helmsman-cicd-infra:gke_logs_dataset is still in use, resourceInUse"
bq rm -r -f "${PROJECT}":"${BQ_LOG_DS}"
