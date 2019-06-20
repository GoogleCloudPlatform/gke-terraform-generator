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

/*

	Package api implements the api used by gke-tf.
	It allows for a user to define a GKE Cluster in yaml.

	An example YAML document:

		kind: gke-cluster
		metadata:
	  	  name: "test-cluster"
		spec:
		  zones:
		    - us-west1-c
		    - us-west1-b
		  private: false
		  region: "us-west1"
		  regional: false
		  addons:
		    istio: true
		    logging: true
		    monitoring: true
		    networkPolicy: true
		    podSecurityPolicy: false
		    hpa: true
		    vpa: false
		    autoscaling: false
		    binaryAuth: true
		    httpLoadBalancing: true
		  network:
		    metadata:
		      name: my-network
		    spec:
		      subnetName: my-subnet
		      subnetRange: "10.0.0.0/24"
		      podSubnetRange: "10.1.0.0/16"
		      serviceSubnetRange: "10.2.0.0/20"
		  version: 1.13.5-gke.10
		  nodePools:
		    - metadata:
		    name: my-node-pool
		    spec:
		      initialNodeCount: 1
		      machineType: n1-standard-1
		      diskSizeGB: 50
		    - metadata:
		      name: my-other-nodepool
		    spec:
		      initialNodeCount: 1
		      machineType: n1-standard-1
		      diskSizeGB: 50
		      diskType: pd-standard
		      tags:
		        - red
		        - white
		  tags:
		    - blue
		    - green

*/
package api

// TODO more
