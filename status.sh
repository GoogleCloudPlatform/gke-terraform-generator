#!/bin/bash

git_describe=$(git describe --tags --always)
echo "GKE_TF_VERSION ${git_describe}"
