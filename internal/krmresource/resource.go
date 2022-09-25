package krmresource

import (
	"fmt"

	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Resources interface {
	// Add adds a resource to the resource list, the name should be unique e.g. gvknns in k8s context
	Add(rn *yaml.RNode)
	Get() []*yaml.RNode
	Copy() Resources
	IsEqual(arnl []*yaml.RNode) (bool, error)
	Print(prefix ...string)
	Write(targetDir string) error
}

func New() Resources {
	return &resources{
		resources: []*yaml.RNode{},
	}
}

type resources struct {
	resources []*yaml.RNode
}

func (x *resources) Add(rn *yaml.RNode) {
	for i, rnode := range x.resources {
		if rn.GetApiVersion() == rnode.GetApiVersion() &&
			rn.GetKind() == rnode.GetKind() &&
			rn.GetName() == rnode.GetName() {
			x.resources[i] = rn
		}
	}
	x.resources = append(x.resources, rn)
}

func (x *resources) Get() []*yaml.RNode {
	return x.resources
}

func (x *resources) Copy() Resources {
	resources := New()
	for _, r := range x.resources {
		resources.Add(r.Copy())
	}
	return resources
}

// IsEqual validates if the resources are equal or not
func (x *resources) IsEqual(arnl []*yaml.RNode) (bool, error) {
	for _, nr := range x.resources {
		found := false
		for i, ar := range arnl {
			nrStr, err := nr.String()
			if err != nil {
				return false, err
			}
			ar.SetAnnotations(map[string]string{})
			arStr, err := ar.String()
			if err != nil {
				return false, err
			}
			if nrStr == arStr {
				found = true
				arnl = append(arnl[:i], arnl[i+1:]...)
			}
		}
		if !found {
			return false, nil
		}
	}
	// this means some entries should be deleted
	// hence the resource ar enot equal
	if len(arnl) != 0 {
		return false, nil
	}
	return true, nil
}

func (x *resources) Print(prefix ...string) {
	fmt.Println()
	fmt.Println(prefix)
	fmt.Println()
	for _, rn := range x.resources {
		fmt.Printf("Resource apiversion: %s, kind: %s, name: %s\n", rn.GetApiVersion(), rn.GetKind(), rn.GetName())
		for k, v := range rn.GetAnnotations() {
			fmt.Printf("  annotation key: %s, value: %s\n", k, v)
		}
		for k, v := range rn.GetLabels() {
			fmt.Printf("  label key: %s, value: %s\n", k, v)
		}
	}
}

func (x *resources) Write(targetDir string) error {
	for _, rn := range x.resources {
		fp := fileutil.GetFullPath(targetDir, rn.GetAnnotations()[kioutil.PathAnnotation])

		// remove annotations that are used for internal purposes
		delete(rn.GetAnnotations(), kioutil.PathAnnotation)
		delete(rn.GetAnnotations(), kioutil.IndexAnnotation)
		delete(rn.GetAnnotations(), "config.kubernetes.io/index")

		if err := fileutil.WriteFileFromRNode(fp, rn); err != nil {
			return err
		}
	}
	return nil
}
