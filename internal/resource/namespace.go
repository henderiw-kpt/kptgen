package resource

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	NamespaceKind = "Namespace"
)

func (rn *Resource) RenderNamespace() (*yaml.RNode, error) {
	x := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       NamespaceKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: rn.GetNameSpace(),
			Labels: rn.GetK8sLabels(),
			Annotations: map[string]string{
				kioutil.PathAnnotation:  rn.GetRelativeFilePath(NamespaceKind),
				kioutil.IndexAnnotation: "0",
			},
		},
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(NamespaceKind), x)
}
