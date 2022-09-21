package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

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

func includeFile(path string, match []string) bool {
	for _, m := range match {
		file := filepath.Base(path)
		if matched, err := filepath.Match(m, file); err == nil && matched {
			return true
		}
	}
	return false
}

func ReadFiles(source string, recursive bool, match []string) ([]string, error) {
	filePaths := make([]string, 0)
	if recursive {
		err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// check for the file extension
			if includeFile(filepath.Ext(info.Name()), match) {
				filePaths = append(filePaths, path)
			} else {
				// check for exact file match
				if includeFile(info.Name(), match) {
					filePaths = append(filePaths, path)
				}
			}
			return nil
		})
		if err != nil {
			return filePaths, err
		}
	} else {
		if includeFile(filepath.Ext(source), match) {
			filePaths = append(filePaths, source)
		} else {
			files, err := os.ReadDir(source)
			if err != nil {
				return filePaths, err
			}
			for _, info := range files {
				if includeFile(filepath.Ext(info.Name()), match) {
					path := filepath.Join(source, info.Name())
					filePaths = append(filePaths, path)
				}
			}
		}
	}
	return filePaths, nil
}

// Code adapted from Porch internal cmdrpkgpull and cmdrpkgpush
func ResourcesToPackageBuffer(resources map[string]string, match []string) (*kio.PackageBuffer, error) {
	keys := make([]string, 0, len(resources))
	//fmt.Printf("ResourcesToPackageBuffer: resources: %v\n", len(resources))
	for k := range resources {
		//fmt.Printf("ResourcesToPackageBuffer: resources key %s\n", k)
		if !includeFile(k, match) {
			continue
		}
		//fmt.Printf("ResourcesToPackageBuffer: resources key append %s\n", k)
		keys = append(keys, k)
	}
	sort.Strings(keys)
	//fmt.Printf("keys: %v\n", keys)

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

func CreateFile(targetDir string, node *yaml.RNode) error {
	fileName := getPath(node)
	// exclude patch files
	if strings.Contains(fileName, "patch") {
		return nil
	}
	if v, ok := node.GetAnnotations()["config.kubernetes.io/index"]; ok {
		if v != "0" {
			split := strings.Split(fileName, ".")
			fileName = strings.Join([]string{split[0] + v, split[1]}, ".")
		}
	}

	fullFileName := filepath.Join(targetDir, fileName)
	pathSplit := strings.Split(fullFileName, "/")
	if len(pathSplit) > 1 {
		path := filepath.Join(pathSplit[:(len(pathSplit) - 1)]...)
		if err := EnsureDir("TARGET_DIR", path, true); err != nil {
			return err
		}
	}

	f, err := os.Create(fullFileName)
	if err != nil {
		return fmt.Errorf("cannot create file: %s", fullFileName)
	}
	node.SetAnnotations(map[string]string{})
	str, err := node.String()
	if err != nil {
		return fmt.Errorf("cannot stringify node err: %s", err.Error())
	}

	_, err = f.WriteString(str)
	if err != nil {
		return fmt.Errorf("cannot write string to file %s err: %s", fullFileName, err.Error())
	}
	return f.Close()
}

func UpdateFile(targetDir, data string, node *yaml.RNode) error {
	fileName := getPath(node)

	targetDirSplit := strings.Split(targetDir, "/")
	if len(targetDirSplit) > 1 {
		targetDir = filepath.Join(targetDirSplit[:(len(targetDirSplit) - 1)]...)
	} else {
		targetDir = "."
	}
	fullFileName := filepath.Join(targetDir, fileName)
	f, err := os.Create(fullFileName)
	if err != nil {
		return fmt.Errorf("cannot create file: %s", fullFileName)
	}
	_, err = f.WriteString(data)
	if err != nil {
		return fmt.Errorf("cannot write string to file %s err: %s", fullFileName, err.Error())
	}
	return f.Close()

}

func getPath(node *yaml.RNode) string {
	ann := node.GetAnnotations()
	if path, ok := ann[kioutil.PathAnnotation]; ok {
		return path
	}
	name := node.GetName()
	if name == "" {
		//name = "unnamed" + strconv.Itoa(id)
		return ""
	}
	return fmt.Sprintf("%s.yaml", name)
}

func GetResosurcePathFromConfigPath(targetDir, configPath string) string {
	split1 := strings.Split(targetDir, "/")
	split2 := strings.Split(configPath, "/")

	idx := 0
	for i := range split1 {
		if split1[i] != split2[i] {
			break
		}
		idx = i
	}
	if idx > 0 {
		configPath = filepath.Join(split2[(idx):]...)
	}
	return configPath
}
