package resid

import (
	"strings"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Gvk struct {
	Group   string `json:"group,omitempty" yaml:"group,omitempty"`
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	Kind    string `json:"kind,omitempty" yaml:"kind,omitempty"`
}

func NewGvk(g, v, k string) Gvk {
	result := Gvk{Group: g, Version: v, Kind: k}
	return result
}

func GvkFromNode(r *yaml.RNode) Gvk {
	g, v := ParseGroupVersion(r.GetApiVersion())
	return NewGvk(g, v, r.GetKind())
}

// ParseGroupVersion parses a KRM metadata apiVersion field.
func ParseGroupVersion(apiVersion string) (group, version string) {
	if i := strings.Index(apiVersion, "/"); i > -1 {
		return apiVersion[:i], apiVersion[i+1:]
	}
	return "", apiVersion
}

func (x Gvk) IsSelected(selector *Gvk) bool {
	if selector == nil {
		return true
	}
	//fmt.Println(selector.Group, x.Group)
	if len(selector.Group) > 0 {
		if x.Group != selector.Group {
			return false
		}
	}
	//fmt.Println(selector.Version, x.Version)
	if len(selector.Version) > 0 {
		if x.Version != selector.Version {
			return false
		}
	}
	//fmt.Println(selector.Kind, x.Kind)
	if len(selector.Kind) > 0 {
		if x.Kind != selector.Kind {
			return false
		}
	}
	return true
}
