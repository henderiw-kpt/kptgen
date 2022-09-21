package resource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

const (
	DeploymentKind = "Deployment"
)

func RenderDeployment(rn *Resource, fc *kptgenv1alpha1.PodSpec) error {
	rn.Kind = DeploymentKind

	fc.Pod.ServiceAccountName = rn.GetName()

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
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      rn.GetName(),
					Namespace: rn.Namespace,
					Labels: map[string]string{
						rn.GetLabelKey(): rn.Name,
					},
				},
				Spec: fc.Pod,
			},
		},
	}
	return ApplyDeployment(rn, x)
}

func ApplyDeployment(rn *Resource, x *appsv1.Deployment) error {
	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	var fp string
	path, ok := x.Annotations["config.kubernetes.io/path"]
	fmt.Println(path)
	if ok {
		fp = path
		fmt.Println(rn.TargetDir)
		pathSplit := strings.Split(rn.TargetDir, "/")
		if len(pathSplit) > 1 {
			fmt.Println(pathSplit)
			pp := filepath.Join(pathSplit[:(len(pathSplit) - 1)]...)
			fmt.Println(pp)
			fp = filepath.Join(pp, fp)
		}
	}
	if fp == "" {
		fp = rn.GetFilePath()
	}

	fmt.Println(fp)

	if err := fileutil.EnsureDir(DeploymentKind, filepath.Dir(fp), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(fp, []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}
