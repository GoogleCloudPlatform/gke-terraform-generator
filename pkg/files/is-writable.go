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

// +build !windows

package files

import (
	"os"
	"syscall"

	"k8s.io/klog"
)

func IsWritable(path string) (isWritable bool, err error) {
	isWritable = false

	isWritable, err = isWritableDirectory(path)
	if !isWritable || err != nil {
		return
	}

	isWritable = false

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		klog.Error("Unable to get stat")
		return
	}

	err = nil
	if uint32(os.Geteuid()) != stat.Uid {
		isWritable = false
		klog.Error("User doesn't have permission to write to this directory")
		return
	}

	isWritable = true
	return
}
