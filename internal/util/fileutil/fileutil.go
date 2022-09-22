package fileutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
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

func ResourcesToPackageBuffer(resources map[string]string, match []string) (*kio.PackageBuffer, error) {
	keys := make([]string, 0, len(resources))
	for k := range resources {
		if !includeFile(k, match) {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Create kio readers
	inputs := []kio.Reader{}
	for _, k := range keys {
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

func CreateFileFromRObject(originTag, fp string, x runtime.Object) error {
	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := EnsureDir(originTag, filepath.Dir(fp), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fp, []byte(b.String()), 0644); err != nil {
		return err
	}
	fmt.Printf("creating resource ... : %s %s\n", x.GetObjectKind().GroupVersionKind(), fp)
	return nil
}

func UpdateFileFromRObject(originTag, fp string, x runtime.Object) error {
	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := EnsureDir(originTag, filepath.Dir(fp), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fp, []byte(b.String()), 0644); err != nil {
		return err
	}
	fmt.Printf("updating resource ... : %s %s\n", x.GetObjectKind().GroupVersionKind(), fp)
	return nil
}

func CreateFileFromRNode(targetDir string, node *yaml.RNode) error {
	fileName := GetPathFromRNode(node)
	// exclude patch files
	// NEED TO FIND A BETTER SOLUTION
	if strings.Contains(fileName, "patch") {
		return nil
	}
	// if multiple yaml objects exists we split them accross multiple files
	if v, ok := node.GetAnnotations()["config.kubernetes.io/index"]; ok {
		if v != "0" {
			split := strings.Split(fileName, ".")
			fileName = strings.Join([]string{split[0] + v, split[1]}, ".")
		}
	}

	fp := filepath.Join(targetDir, fileName)
	if err := EnsureDir("TARGET_DIR", filepath.Dir(fp), true); err != nil {
		return err
	}
	// clear annotations
	node.SetAnnotations(map[string]string{})

	fmt.Printf("creating resource ... : %s %s\n", node.GetApiVersion(), fp)
	return ioutil.WriteFile(fp, []byte(node.MustString()), 0644)
}

func UpdateFileFromRNode(targetDir, data string, node *yaml.RNode) error {
	fileName := GetPathFromRNode(node)
	fp := GetFullPath(targetDir, fileName)

	fmt.Printf("updating resource ... : %s %s %s %s\n", node.GetApiVersion(), node.GetKind(), node.GetName(), fp)
	return ioutil.WriteFile(fp, []byte(data), 0644)
}

func GetPathFromRNode(node *yaml.RNode) string {
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

// GetRelativePath returns the path which kioutil would return
// the relative path in the package, used for comparisons
func GetRelativePath(targetDir, fp string) string {
	split1 := strings.Split(targetDir, "/")
	split2 := strings.Split(fp, "/")

	idx := 0
	for i := range split1 {
		if split1[i] != split2[i] {
			break
		}
		idx = i
	}
	// relative path
	rp := fp
	if idx > 0 {
		rp = filepath.Join(split2[(idx):]...)
	}
	return rp
}

// GetFullPath returns the path from where the binary is run
// used to create/update files in the filesystem
func GetFullPath(targetDir, rp string) string {
	// if the targetDir has multiple subdirectories we need to augment the path
	pathSplit := strings.Split(targetDir, "/")
	if len(pathSplit) > 1 {
		pp := filepath.Join(pathSplit[:(len(pathSplit) - 1)]...)
		rp = filepath.Join(pp, rp)
	}
	return rp
}
