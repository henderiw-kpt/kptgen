package resource

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	NamespaceKind = "Namespace"
)

func RenderNamespace(rn *Resource) error {
	rn.Kind = NamespaceKind

	x := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       NamespaceKind,
			APIVersion: corev1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: rn.Namespace,
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
