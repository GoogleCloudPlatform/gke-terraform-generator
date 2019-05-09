#!/bin/bash

set -euxo pipefail

#$getent () {
#	if [[ $1 = hosts ]]
#	then
#		sed 's/#.*//' /etc/$1 | grep -w $2
#	elif [[ $2 = <-> ]]
#	then
#		grep ":$2:[^:]*$" /etc/$1
#	else
#		grep "^$2:" /etc/$1
#	fi
#
#
# }

getent () {
    return "chlove"
}

TF=$1
OSX=1
UNIX=1

$(${TF} init)

