package resource

import (
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	StatefullSetKind = "StatefulSet"
)

func (rn *Resource) RenderProviderStatefulSet(fc *kptgenv1alpha1.PodSpec) error {
	rn.Kind = StatefullSetKind

	fc.PodTemplate.Spec.ServiceAccountName = rn.GetServiceAccountName()
	fc.PodTemplate.ObjectMeta.Name = rn.GetResourceName()
	fc.PodTemplate.ObjectMeta.Namespace = rn.GetNameSpace()
	if len(fc.PodTemplate.ObjectMeta.Labels) == 0 {
		fc.PodTemplate.ObjectMeta.Labels = map[string]string{
			rn.GetLabelKey(): rn.Name,
		}
	} else {
		fc.PodTemplate.ObjectMeta.Labels[rn.GetLabelKey()] = rn.Name
	}

	x := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       StatefullSetKind,
			APIVersion: appsv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetResourceName(),
			Namespace: rn.GetNameSpace(),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: fc.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					rn.GetLabelKey(): rn.Name,
				},
			},
			Template: fc.PodTemplate,
		},
	}

	return fileutil.CreateFileFromRObject(rn.GetFilePath(""), x)
}
