package pkgutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
)

func GetPackage(path string, m []string) (*kio.PackageBuffer, error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	var files []string
	if fileInfo.IsDir() {
		// this is a dir
		files, err = fileutil.ReadFiles(path, true, m)
		if err != nil {
			return nil, err
		}
	} else {
		// this is a single file
		files = []string{path}
	}

	inputs := []kio.Reader{}
	for _, path := range files {
		if includeFile(path, m) {
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("cannot read file: %s", path)
			}

			pathSplit := strings.Split(path, "/")
			if len(pathSplit) > 1 {
				path = filepath.Join(pathSplit[1:]...)
			}

			inputs = append(inputs, &kio.ByteReader{
				Reader: strings.NewReader(string(yamlFile)),
				SetAnnotations: map[string]string{
					kioutil.PathAnnotation: path,
				},
				DisableUnwrapping: true,
			})
		}
	}

	var pb kio.PackageBuffer
	err = kio.Pipeline{
		Inputs:  inputs,
		Filters: []kio.Filter{},
		Outputs: []kio.Writer{&pb},
	}.Execute()
	if err != nil {
		return nil, fmt.Errorf("kio error: %s", err)
	}

	return &pb, nil
}

func includeFile(path string, match []string) bool {
	for _, m := range match {
		file := filepath.Base(path)
		if matched, err := filepath.Match(m, file); err == nil && matched {
			return true
		}
	}
	return false
}
