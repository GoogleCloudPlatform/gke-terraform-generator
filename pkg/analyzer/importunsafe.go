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

package analyzer

import (
	"fmt"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "importunsafe",
	Doc:  "reports imports of package unsafe",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		// mousetrap is imported by windows and uses unsafe
		if f.Name.Name == "mousetrap" {
			continue
		}
		fmt.Printf("name: %s", f.Name.Name)
		for _, imp := range f.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err == nil && path == "unsafe" {
				msg := fmt.Sprintf("package unsafe must not be imported pkg: %s", f.Name.Name)
				pass.Reportf(imp.Pos(), msg)
			}
		}
	}
	return nil, nil
}
