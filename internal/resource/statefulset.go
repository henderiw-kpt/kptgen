package resource

import (
	"strings"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

const (
	StatefullSetKind = "StatefulSet"
)

func (rn *Resource) RenderProviderStatefulSet(fc *kptgenv1alpha1.PodSpec) (*yaml.RNode, error) {
	fc.PodTemplate.Spec.ServiceAccountName = rn.GetServiceAccountName()
	fc.PodTemplate.ObjectMeta.Name = rn.GetResourceName()
	fc.PodTemplate.ObjectMeta.Namespace = rn.GetNameSpace()
	if len(fc.PodTemplate.ObjectMeta.Labels) == 0 {
		fc.PodTemplate.ObjectMeta.Labels = map[string]string{
			rn.GetLabelKey(): rn.GetPackagePodName(),
		}
	} else {
		fc.PodTemplate.ObjectMeta.Labels[rn.GetLabelKey()] = rn.GetPackagePodName()
	}

	for k, v := range rn.GetK8sLabels() {
		fc.PodTemplate.ObjectMeta.Labels[k] = v
	}

	x := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       StatefullSetKind,
			APIVersion: appsv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetResourceName(),
			Namespace: rn.GetNameSpace(),
			Annotations: map[string]string{
				kioutil.PathAnnotation:  rn.GetRelativeFilePath(StatefullSetKind),
				kioutil.IndexAnnotation: "0",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: fc.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					rn.GetLabelKey(): rn.GetPackagePodName(),
				},
			},
			Template: fc.PodTemplate,
		},
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(StatefullSetKind), x)
}

func (rn *Resource) UpdateStatefulSet(fnCfgName string, fnCfg kptgenv1alpha1.Config, node *yaml.RNode) (*yaml.RNode, error) {
	x := &appsv1.StatefulSet{}
	if err := sigyaml.Unmarshal([]byte(node.MustString()), &x); err != nil {
		return nil, err
	}

	if len(fnCfg.Spec.Services) != 0 {
		// update the labels with the service selctor key
		x.Spec.Selector.MatchLabels[rn.GetLabelKey()] = rn.GetPackagePodName()
		x.Spec.Template.Labels[rn.GetLabelKey()] = rn.GetPackagePodName()
	}

	/*
		if fnCfgName == "grpc" {
			for i, c := range x.Spec.Template.Spec.Containers {
				if c.Name == "controller" {
					found := false
					newEnv := []corev1.EnvVar{}
					for _, env := range c.Env {
						if env.Name == "GRPC_CERT_SECRET_NAME" {
							found = true
							env.Value = rn.GetCertificateName()
							newEnv = append(newEnv, env)
						} else {
							newEnv = append(newEnv, env)
						}
					}
					if !found {
						newEnv = append(newEnv, corev1.EnvVar{
							Name:  "GRPC_CERT_SECRET_NAME",
							Value: rn.GetCertificateName(),
						})
					}
					x.Spec.Template.Spec.Containers[i].Env = newEnv
				}
			}
		}
	*/

	if fnCfg.Spec.Webhook || fnCfg.Spec.Volume || fnCfg.Spec.Certificate.IssuerRef != "" {
		found := false
		vol := rn.BuildVolume(fnCfg.Spec.Certificate.IssuerRef != "")
		for _, volume := range x.Spec.Template.Spec.Volumes {
			if volume.Name == vol.Name {
				found = true
				volume = vol
			}
		}
		if !found {
			x.Spec.Template.Spec.Volumes = append(x.Spec.Template.Spec.Volumes, vol)
		}
		for i, c := range x.Spec.Template.Spec.Containers {
			if c.Name == fnCfg.Spec.Selector.ContainerName {
				found := false
				volm := rn.BuildVolumeMount(fnCfg.Spec.Certificate.IssuerRef != "")
				for _, volumeMount := range c.VolumeMounts {
					if volumeMount.Name == volm.Name {
						found = true
						volumeMount = volm
					}
				}
				if !found {
					if len(c.VolumeMounts) == 0 {
						x.Spec.Template.Spec.Containers[i].VolumeMounts = make([]corev1.VolumeMount, 0, 1)
					}
					x.Spec.Template.Spec.Containers[i].VolumeMounts = append(x.Spec.Template.Spec.Containers[i].VolumeMounts, volm)
				}
			}
		}
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	// the path must exist since we read the resource from the filesystem

	//fp := fileutil.GetFullPath(rn.TargetDir, x.Annotations[kioutil.PathAnnotation])
	//return fileutil.UpdateFileFromRObject(fp, x)
}
