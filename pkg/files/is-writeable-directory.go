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

	"k8s.io/klog"
)

func isWritableDirectory(path string) (isWritable bool, err error) {
	isWritable = false
	info, err := os.Stat(path)
	if err != nil {
		klog.Errorf("Path doesn't exist: %s", path)
		return
	}

	err = nil
	if !info.IsDir() {
		klog.Errorf("Path isn't a directory: %s", path)
		return
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		klog.Errorf("Write permission bit is not set on this directory for user: %s", path)
		return
	}

	isWritable = true
	return
}
