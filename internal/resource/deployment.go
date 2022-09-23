package resource

import (
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DeploymentKind = "Deployment"
)

func (rn *Resource) RenderDeployment(fc *kptgenv1alpha1.PodSpec) error {
	rn.Kind = DeploymentKind

	fc.PodTemplate.Spec.ServiceAccountName = rn.GetPackageResourceName("")
	fc.PodTemplate.ObjectMeta.Name = rn.GetName()
	fc.PodTemplate.ObjectMeta.Namespace = rn.GetNameSpace()
	if len(fc.PodTemplate.ObjectMeta.Labels) == 0 {
		fc.PodTemplate.ObjectMeta.Labels = map[string]string{
			rn.GetLabelKey(): rn.Name,
		}
	} else {
		fc.PodTemplate.ObjectMeta.Labels[rn.GetLabelKey()] = rn.Name
	}

	x := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       DeploymentKind,
			APIVersion: appsv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetName(),
			Namespace: rn.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
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
