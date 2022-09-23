package resource

import (
	"fmt"
	"reflect"

	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ServiceKind = "Service"
)

func (rn *Resource) RenderService(cfg, obj interface{}) error {
	rn.Kind = ServiceKind

	svc, ok := cfg.(corev1.Service)
	if !ok {
		return fmt.Errorf("wrong object in renderService: %v", reflect.TypeOf(cfg))
	}

	svc.ObjectMeta.Name = rn.GetServiceName()
	svc.ObjectMeta.Namespace = rn.GetNameSpace()
	if len(svc.ObjectMeta.Labels) == 0 {
		svc.ObjectMeta.Labels = map[string]string{
			rn.GetLabelKey(): rn.PackageName,
		}
	} else {
		svc.ObjectMeta.Labels[rn.GetLabelKey()] = rn.Name
	}
	if len(svc.Spec.Selector) == 0 {
		svc.Spec.Selector = map[string]string{
			rn.GetLabelKey(): rn.PackageName,
		}
	} else {
		svc.Spec.Selector[rn.GetLabelKey()] = rn.Name
	}

	x := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       ServiceKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: svc.ObjectMeta,
		Spec:       svc.Spec,
	}

	return fileutil.CreateFileFromRObject(rn.GetFilePath(""), x)
}
