package resource

import (
	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ServiceAccountKind = "ServiceAccount"
)

func (rn *Resource) RenderServiceAccount(fc *kptgenv1alpha1.PodSpec) error {
	rn.Kind = ServiceAccountKind

	x := &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       ServiceAccountKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetControllerName(""),
			Namespace: rn.Namespace,
		},
	}

	return fileutil.CreateFileFromRObject(ServiceAccountKind, rn.GetFilePath(""), x)

}
