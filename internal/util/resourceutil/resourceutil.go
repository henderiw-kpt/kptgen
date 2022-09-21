package resourceutil

import (
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/kustomize/kyaml/kio"
	sigyaml "sigs.k8s.io/yaml"
)

func GetCRDs(pb *kio.PackageBuffer) ([]extv1.CustomResourceDefinition, error) {
	crdObjects := make([]extv1.CustomResourceDefinition, 0)
	for _, node := range pb.Nodes {
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
