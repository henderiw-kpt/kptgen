package resource

import (
	"strings"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	ServiceAccountKind = "ServiceAccount"
)

func (rn *Resource) RenderServiceAccount(fc *kptgenv1alpha1.PodSpec) (*yaml.RNode, error) {
	x := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       ServiceAccountKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetServiceAccountName(),
			Namespace: rn.Namespace,
			Annotations: map[string]string{
				kioutil.PathAnnotation:  rn.GetRelativeFilePath(ServiceAccountKind),
				kioutil.IndexAnnotation: "0",
			},
		},
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(ServiceAccountKind), x)

}
