package resourceutil

import (
	"github.com/henderiw-kpt/kptgen/internal/krmresource"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	sigyaml "sigs.k8s.io/yaml"
)

func GetCRDs(r krmresource.Resources) ([]extv1.CustomResourceDefinition, error) {
	crdObjects := make([]extv1.CustomResourceDefinition, 0)
	for _, node := range r.Get() {
		if node.GetKind() == "CustomResourceDefinition" {
			crd := extv1.CustomResourceDefinition{}
			if err := sigyaml.Unmarshal([]byte(node.MustString()), &crd); err != nil {
				return crdObjects, err
			}
			crdObjects = append(crdObjects, crd)
		}
	}
	return crdObjects, nil
}
