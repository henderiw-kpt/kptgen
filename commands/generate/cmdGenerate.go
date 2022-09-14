package generate

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	docs "github.com/henderiw-nephio/kptgen/internal/docs/generated/generatedocs"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{
		Ctx: ctx,
	}
	c := &cobra.Command{
		Use:     "generate SOURCE_DIR TARGET_DIR",
		Args:    cobra.MaximumNArgs(2),
		Short:   docs.GenerateShort,
		Long:    docs.GenerateShort + "\n" + docs.GenerateLong,
		Example: docs.GenerateExamples,
		RunE:    r.runE,
	}

	c.Flags().StringVar(&r.Description, "description", "sample description", "short description of the package.")
	c.Flags().StringSliceVar(&r.Keywords, "keywords", []string{}, "list of keywords for the package.")
	c.Flags().StringVar(&r.Site, "site", "", "link to page with information about the package.")
	//cmdutil.FixDocs("kpt", parent, c)
	r.Command = c
	return r
}

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command     *cobra.Command
	Keywords    []string
	Name        string
	Description string
	Site        string
	Ctx         context.Context
}

func (r *Runner) runE(c *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("SOURCE_DIR and TARGET_DIR are required positional arguments; %d provided", len(args))
	}

	sourceDir := args[0]
	targetDir := args[1]

	if err := fileutil.EnsureDir("SOURCE_DIR", sourceDir, false); err != nil {
		return err
	}

	if err := fileutil.EnsureDir("TARGET_DIR", targetDir, true); err != nil {
		return err
	}

	// read only yml, yaml files
	files, err := fileutil.ReadFiles(sourceDir, true)
	if err != nil {
		return err
	}

	inputs := []kio.Reader{}
	for _, path := range files {
		if includeFile(path) && !excludeFile(path) {
			yamlFile, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("cannot read file: %s", path)
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
		return fmt.Errorf("kio error: %s", err)
	}

	for _, node := range pb.Nodes {
		if node.GetKind() != "" {
			fileName := getPath(node)
			if v, ok := node.GetAnnotations()["config.kubernetes.io/index"]; ok {
				if v != "0" {
					split := strings.Split(fileName, ".")
					fileName = strings.Join([]string{split[0] + v, split[1]}, ".")
				}
			}

			fmt.Printf("path: %s\n", getPath(node))
			fmt.Printf("annotations: %v\n", node.GetAnnotations())
			fmt.Printf("node: kind: %s, apiversion: %s, name: %s\n", node.GetKind(), node.GetApiVersion(), node.GetName())

			fullFileName := filepath.Join(targetDir, fileName)
			pathSplit := strings.Split(fullFileName, "/")
			if len(pathSplit) > 1 {
				path := filepath.Join(pathSplit[:(len(pathSplit) - 1)]...)
				if err := fileutil.EnsureDir("TARGET_DIR", path, true); err != nil {
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
			f.Close()
		}
	}

	return nil

}

func getPath(node *yaml.RNode) string {
	ann := node.GetAnnotations()
	if path, ok := ann[kioutil.PathAnnotation]; ok {
		return path
	}
	/*
		ns := node.GetNamespace()
		if ns == "" {
			ns = "non-namespaced"
		}
	*/
	name := node.GetName()
	if name == "" {
		//name = "unnamed" + strconv.Itoa(id)
		return ""
	}
	// TODO: harden for escaping etc.
	//	return path.Join(ns, fmt.Sprintf("%s.yaml", name))
	return fmt.Sprintf("%s.yaml", name)
}

var MatchAll = []string{"*.yaml", "*.yml"}
var ExcludeAll = []string{"patch", "samples"}

func includeFile(path string) bool {
	for _, m := range MatchAll {
		file := filepath.Base(path)
		if matched, err := filepath.Match(m, file); err == nil && matched {
			return true
		}
	}
	return false
}

func excludeFile(path string) bool {
	for _, m := range ExcludeAll {
		if strings.Contains(path, m) {
			return true
		}
	}
	return false
}
