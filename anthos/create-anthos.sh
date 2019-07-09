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

set -o errexit
set -o nounset
set -o pipefail
# set -x

# gcloud and kubectl are required for this POC
command -v gcloud >/dev/null 2>&1 || { \
 echo >&2 "I require gcloud but it's not installed.  Aborting."; exit 1; }
command -v kubectl >/dev/null 2>&1 || { \
 echo >&2 "I require kubectl but it's not installed.  Aborting."; exit 1; }
command -v nomos >/dev/null 2>&1 || { \
 echo >&2 "I require nomos but it's not installed.  Aborting."; exit 1; }
command -v kustomize >/dev/null 2>&1 || { \
 echo >&2 "I require kustomize but it's not installed.  Aborting."; exit 1; }

# TODO: install nomos

# create cluster
# we can remove the following lines when using terraform
gcloud compute networks create my-vpc-network \
  --subnet-mode=custom

gcloud compute networks subnets create my-cluster-nodes-subnet \
    --network my-vpc-network \
    --region us-central1 \
    --range 10.0.0.0/24 \
    --secondary-range my-cluster-pod-subnet=10.1.0.0/16,my-cluster-service-subnet=10.2.0.0/20

gcloud beta container clusters create anthos-cluster \
  --region us-central1 \
  --network=my-vpc-network \
  --enable-ip-alias \
  --subnetwork=my-cluster-nodes-subnet \
  --cluster-secondary-range-name=my-cluster-pod-subnet \
  --services-secondary-range-name=my-cluster-service-subnet \
  --enable-network-policy \
  --enable-pod-security-policy \
  --num-nodes=1

# TODO: remove above cluster create commands when we have TF

# "---------------------------------------------------------"
# "-                                                       -"
# "-  Installs Anthos Config Manager                       -"
# "-                                                       -"
# "---------------------------------------------------------"

REPO="anthos-demo"
PROJECT=$(gcloud config list --format 'value(core.project)' 2>/dev/null)
ACCOUNT=$(gcloud config list --format 'value(core.account)' 2>/dev/null)

gcloud beta container clusters get-credentials anthos-cluster  \
  --region us-central1

kubectl create clusterrolebinding cluster-admin-binding \
  --clusterrole cluster-admin \
  --user "${ACCOUNT}"

# create a repo
gcloud source repos create "${REPO}"

# clone the repo
gcloud source repos clone "${REPO}"
cd anthos-demo

echo "Initialize repo"
nomos init
# TODO clone the repo and copy part of it, or should we just have it??
# TODO use $ROOT here
cp -R ../anthos-config-example/namespaces/audit namespaces/

cat <<EOT >> .gitignore
anthos-demo-key
anthos-demo-key.pub
EOT

git add .
git commit -m 'Adding namespace'

echo "Adding a commit to repo"
git push

echo "Creating key pair"
ssh-keygen -t rsa -b 4096 \
 -C "${ACCOUNT}" \
 -N '' \
 -f ./anthos-demo-key

echo "Installing nomos operator"
kubectl apply --filename https://storage.googleapis.com/nomos-release/operator-rc/nomos-operator-v0.1.16-rc.12/nomos-operator.yaml

echo "Adding ssh key to GKE"
# add ssh key to GKE
kubectl create secret generic git-creds \
 --namespace=config-management-system \
 --from-file=ssh=./anthos-demo-key

echo ""
echo "Please add the following ssh public key to the Cloud Source Repo"
echo "Add the SSH Key Here:  https://source.cloud.google.com/user/ssh_keys"
echo ""
echo "----------------------------------------------------------------"
echo ""
cat anthos-demo-key.pub
echo ""
echo "----------------------------------------------------------------"

read -r -p "Press enter to continue, once the SSH key is added"

echo "Creating and applying RBAC Bindings for Anthos to have a PSP"
# TODO use $ROOT
kubectl apply -f ../anthos-rbac-psp.yaml

# TODO move this to customize
echo "Creating and applying the config-management CRD"
cat <<EOF | kubectl apply -f -
apiVersion: addons.sigs.k8s.io/v1alpha1
kind: ConfigManagement
metadata:
  name: config-management
  namespace: config-management-system
spec:
  # clusterName is required and must be unique among all managed clusters
  clusterName: anthos-cluster
  git:
    syncRepo: ssh://${ACCOUNT}@source.developers.google.com:2022/p/${PROJECT}/r/${REPO}
    syncBranch: master
    secretType: ssh
    policyDir: "."
EOF

# put in wait for the deployment to rollout
