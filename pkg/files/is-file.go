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

// IsFile checks to see if a file exists on a defined path.
func IsFile(path string) (isFile bool, err error) {
	isFile = false
	info, err := os.Stat(path)
	if err != nil {
		klog.Errorf("File doesn't exist: %s", path)
		return
	}

	if info.IsDir() {
		klog.Error("Path is a directory was expecting a file")
		return
	}

	return true, nil
}
