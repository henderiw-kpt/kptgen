package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/kio"
)

var MatchAll = []string{"*.yaml", "*.yml"}

// ResolveAbsAndRelPaths returns absolute and relative paths for input path
func ResolveAbsAndRelPaths(path string) (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	var relPath string
	var absPath string
	if filepath.IsAbs(path) {
		// If the provided path is absolute, we find the relative path by
		// comparing it to the current working directory.
		relPath, err = filepath.Rel(cwd, path)
		if err != nil {
			return "", "", err
		}
		absPath = filepath.Clean(path)
	} else {
		// If the provided path is relative, we find the absolute path by
		// combining the current working directory with the relative path.
		relPath = filepath.Clean(path)
		absPath = filepath.Join(cwd, path)
	}

	return absPath, relPath, nil
}

func EnsureDir(originTag, dirName string, create bool) error {
	if _, err := os.Stat(dirName); err != nil {
		if create {
			err = os.MkdirAll(dirName, os.ModePerm)
			if err != nil {
				return fmt.Errorf("%s cannot create directory: %s", originTag, dirName)
			}
		} else {
			return fmt.Errorf("%s does not exist: %s", originTag, dirName)
		}
	}
	if stat, err := os.Stat(dirName); err == nil && !stat.IsDir() {
		return fmt.Errorf("%s %s is not a directory", originTag, dirName)
	}
	return nil
}

func includeFile(path string) bool {
	for _, m := range MatchAll {
		file := filepath.Base(path)
		if matched, err := filepath.Match(m, file); err == nil && matched {
			return true
		}
	}
	return false
}

func ReadFiles(source string, recursive bool) ([]string, error) {
	filePaths := make([]string, 0)
	if recursive {
		err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if includeFile(filepath.Ext(info.Name())) {
				filePaths = append(filePaths, path)
			}
			return nil
		})
		if err != nil {
			return filePaths, err
		}
	} else {
		if includeFile(filepath.Ext(source)) {
			filePaths = append(filePaths, source)
		} else {
			files, err := os.ReadDir(source)
			if err != nil {
				return filePaths, err
			}
			for _, info := range files {
				if includeFile(filepath.Ext(info.Name())) {
					path := filepath.Join(source, info.Name())
					filePaths = append(filePaths, path)
				}
			}
		}
	}
	return filePaths, nil
}

// Code adapted from Porch internal cmdrpkgpull and cmdrpkgpush
func ResourcesToPackageBuffer(resources map[string]string) (*kio.PackageBuffer, error) {
	keys := make([]string, 0, len(resources))
	//fmt.Printf("ResourcesToPackageBuffer: resources: %v\n", len(resources))
	for k := range resources {
		fmt.Printf("ResourcesToPackageBuffer: resources key %s\n", k)
		if !includeFile(k) {
			continue
		}
		//fmt.Printf("ResourcesToPackageBuffer: resources key append %s\n", k)
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Printf("keys: %v\n", keys)

	// Create kio readers
	inputs := []kio.Reader{}
	for _, k := range keys {
		//fmt.Printf("ResourcesToPackageBuffer: key %s\n", k)
		v := resources[k]
		inputs = append(inputs, &kio.ByteReader{
			Reader:            strings.NewReader(v),
			DisableUnwrapping: true,
		})
	}

	var pb kio.PackageBuffer
	err := kio.Pipeline{
		Inputs:  inputs,
		Outputs: []kio.Writer{&pb},
	}.Execute()

	if err != nil {
		return nil, err
	}

	return &pb, nil
}
