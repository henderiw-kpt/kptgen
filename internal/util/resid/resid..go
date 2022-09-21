package resid

import "sigs.k8s.io/kustomize/kyaml/yaml"

type ResId struct {
	Gvk      `json:",inline,omitempty" yaml:",inline,omitempty"`
	FileName string `json:"name,omitempty" yaml:"name,omitempty"`
}

func NewResId(k Gvk, f string) ResId {
	return ResId{Gvk: k, FileName: f}
}

// FromRNode returns the ResId for the RNode
func FromRNode(rn *yaml.RNode) ResId {
	group, version := ParseGroupVersion(rn.GetApiVersion())
	return NewResId(
		Gvk{Group: group, Version: version, Kind: rn.GetKind()},
		rn.GetAnnotations()["internal.config.kubernetes.io/path"],
	)
}

func (id *ResId) IsSelectedBy(selector ResId) bool {
	return (selector.FileName == "" || selector.FileName == id.FileName) &&
		id.Gvk.IsSelected(&selector.Gvk)
}
