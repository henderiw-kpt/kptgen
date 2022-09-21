package resource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	ServiceKind = "Service"
)

func RenderService(rn *Resource, cfg, obj interface{}) error {
	rn.Kind = ServiceKind

	info, ok := cfg.(*kptgenv1alpha1.WebhookSpec)
	if !ok {
		return fmt.Errorf("wrong object in rendercertificate: %v", reflect.TypeOf(cfg))
	}
	x := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       ServiceKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetServiceName(),
			Namespace: rn.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				rn.GetLabelKey(): rn.GetServiceName(),
			},
			Ports: []corev1.ServicePort{
				{
					Name:       WebhookSuffix,
					Port:       info.Service.Port,
					TargetPort: intstr.FromInt(int(info.Service.TargetPort)),
					Protocol:   corev1.Protocol("TCP"),
				},
			},
		},
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := fileutil.EnsureDir(ServiceKind, filepath.Dir(rn.GetFilePath()), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(rn.GetFilePath(), []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}
