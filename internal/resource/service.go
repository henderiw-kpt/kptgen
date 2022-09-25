package resource

import (
	"fmt"
	"reflect"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

const (
	ServiceKind = "Service"
)

func (rn *Resource) RenderService(cfg, obj interface{}) (*yaml.RNode, error) {
	svc, ok := cfg.(corev1.Service)
	if !ok {
		return nil, fmt.Errorf("wrong object in renderService: %v", reflect.TypeOf(cfg))
	}

	svc.ObjectMeta.Name = rn.GetServiceName()
	svc.ObjectMeta.Namespace = rn.GetNameSpace()
	if len(svc.ObjectMeta.Labels) == 0 {
		svc.ObjectMeta.Labels = map[string]string{
			rn.GetLabelKey(): rn.GetPackagePodName(),
		}
	} else {
		svc.ObjectMeta.Labels[rn.GetLabelKey()] = rn.GetPackagePodName()
	}
	if len(svc.Spec.Selector) == 0 {
		svc.Spec.Selector = map[string]string{
			rn.GetLabelKey(): rn.GetPackagePodName(),
		}
	} else {
		svc.Spec.Selector[rn.GetLabelKey()] = rn.GetPackagePodName()
	}
	svc.ObjectMeta.Annotations = map[string]string{
		kioutil.PathAnnotation:  rn.GetRelativeFilePath(ServiceKind),
		kioutil.IndexAnnotation: "0",
	}

	x := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       ServiceKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: svc.ObjectMeta,
		Spec:       svc.Spec,
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(ServiceKind), x)
}
