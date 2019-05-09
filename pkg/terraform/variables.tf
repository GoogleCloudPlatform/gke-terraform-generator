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

variable "cluster_name" {
  description = ""
  default = "{{.ObjectMeta.Name}}"
}

variable "project_id" {
  description = ""
  default = "{{.Spec.ProjectId}}"
}

variable "region" {
  description = ""
  default = "{{.Spec.Region}}"
}

{{- if .Spec.Zones }}
variable "zones" {
  description = ""
  // TODO fix bug when we have a single zone
  // TODO fix bug when we do not have zones
  default = ["{{StringsJoin .Spec.Zones "\",\""}}"]
}
{{- end}}
