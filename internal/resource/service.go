package resource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
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
			rn.GetLabelKey(): rn.Name,
		}
	} else {
		svc.ObjectMeta.Labels[rn.GetLabelKey()] = rn.Name
	}
	if len(svc.Spec.Selector) == 0 {
		svc.Spec.Selector = map[string]string{
			rn.GetLabelKey(): rn.Name,
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

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := fileutil.EnsureDir(ServiceKind, filepath.Dir(rn.GetFilePath("")), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(rn.GetFilePath(""), []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}
