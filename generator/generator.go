package generator

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// Generate reads in all the template files from the given directory
// and executes them with the given data, then writes the out put to
// the output dir
func Generate(data interface{}, tplDir string, outDir string) error {
	err := filepath.Walk(tplDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext != ".tf" && ext != ".tfvars" {
			return nil
		}

		tplData, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		if outDir == "" {
			outDir = filepath.Dir(p)
		}

		f, err := os.Create(path.Join(outDir, strings.TrimSuffix(info.Name(), ".tpl")))
		if err != nil {
			return err
		}
		defer f.Close()

		if err := ExecTemplate(f, tplData, data); err != nil {
			return err
		}

		return nil
	})
	return err
}

// ExecTemplate is a simple wrapper around go temlates
func ExecTemplate(wr io.Writer, tmpl []byte, data interface{}) error {
	// hash template
	h := sha256.New()
	h.Write(tmpl)
	hs := hex.EncodeToString(h.Sum(nil))

	// parse template
	t, err := template.New(hs).Parse(string(tmpl))
	if err != nil {
		return err
	}

	// process template with data
	if err := t.Execute(wr, data); err != nil {
		return err
	}
	return nil
}
