// +build !windows

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

package files

import (
	"os"
	"syscall"

	"k8s.io/klog"
)

// IsWritable checks to see if a user has permissions to write a file
// to a directory.
func IsWritable(path string) (isWritable bool, err error) {
	isWritable = false

	isWritable, err = isWritableDirectory(path)
	if !isWritable || err != nil {
		return
	}

	isWritable = false

	var stat syscall.Stat_t
	if err = syscall.Stat(path, &stat); err != nil {
		klog.Errorf("Error to getting stat: %v", err)
		return
	}

	err = nil
	if uint32(os.Geteuid()) != stat.Uid {
		isWritable = false
		klog.Errorf("User doesn't have permission to write to this directory: %s", path)
		return
	}

	return true, nil
}
