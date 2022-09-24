package resource

import (
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NamespaceKind = "Namespace"
)

func (rn *Resource) RenderNamespace() error {
	rn.Kind = NamespaceKind

	x := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       NamespaceKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: rn.GetNameSpace(),
		},
	}

	return fileutil.CreateFileFromRObject(rn.GetFilePath(""), x)
}
