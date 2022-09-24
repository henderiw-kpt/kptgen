package pkgconfig

import (
	"fmt"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/resource"
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

func (r *pkgConfig) deployConfig(node *yaml.RNode) error {
	fnCfg := kptgenv1alpha1.Config{}
	if err := sigyaml.Unmarshal([]byte(node.MustString()), &fnCfg); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	rn := &resource.Resource{
		Kind:        kptgenv1alpha1.FnConfigKind,
		PackageName: r.kptFile.GetName(),
		PodName:     fnCfg.Spec.Selector.Name, // TODO trim it
		Name:        node.GetName(),
		Namespace:   r.kptFile.GetNamespace(),
		TargetDir:   r.targetDir,
		SubDir:      node.GetName(),
		//NameKind:     resource.NameKindPackageResource,
		//PathNameKind: resource.NameKindKindResource,
	}

	crdObjects := make([]extv1.CustomResourceDefinition, 0)
	var podNode *yaml.RNode
	for _, node := range r.pb.Nodes {
		switch node.GetKind() {
		case "CustomResourceDefinition":
			crd := extv1.CustomResourceDefinition{}
			if err := sigyaml.Unmarshal([]byte(node.MustString()), &crd); err != nil {
				return err
			}
			crdObjects = append(crdObjects, crd)
		case fnCfg.Spec.Selector.Kind:
			resid := resid.FromRNode(node)
			if resid.IsSelectedBy(fnCfg.Spec.Selector.ResId) {
				if found, err := r.validatePodContainer(fnCfg, node); err != nil {
					return err
				} else {
					if found {
						// TBD what todo if already found, we would expect 1 container that matches
						podNode = node
					}
				}
			}
		}
	}
	if podNode == nil {
		return fmt.Errorf("container pod not found")
	}

	// TODO
	// ClusterRoles -> to add a clusterrole bonding
	// PermissionRequests -> Todo
	// Containers -> Todo

	// render service
	for _, service := range fnCfg.Spec.Services {
		if err := rn.RenderService(service, crdObjects); err != nil {
			return err
		}
	}
	// render webhook
	if fnCfg.Spec.Webhook {
		if err := rn.RenderMutatingWebhook(fnCfg.Spec, crdObjects); err != nil {
			return err
		}
		if err := rn.RenderValidatingWebhook(fnCfg.Spec, crdObjects); err != nil {
			return err
		}
	}
	// render certificate
	if fnCfg.Spec.Certificate.IssuerRef != "" {
		if err := rn.RenderCertificate(fnCfg.Spec, crdObjects); err != nil {
			return err
		}
	}

	// mutate deployment or statefulset
	switch fnCfg.Spec.Selector.Kind {
	case "Deployment":
		x := &appsv1.Deployment{}
		if err := sigyaml.Unmarshal([]byte(podNode.MustString()), &x); err != nil {
			return err
		}
		// update the labels with the service selctor key
		x.Spec.Selector.MatchLabels[rn.GetLabelKey()] = rn.PackageName
		x.Spec.Template.Labels[rn.GetLabelKey()] = rn.PackageName

		found := false
		vol := rn.BuildVolume()
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
				volm := rn.BuildVolumeMount()
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

		// the path must exist since we read the resource from the filesystem

		fp := fileutil.GetFullPath(rn.TargetDir, x.Annotations[kioutil.PathAnnotation])
		return fileutil.UpdateFileFromRObject(fp, x)

	case "StatefulSet":
		x := &appsv1.StatefulSet{}
		if err := sigyaml.Unmarshal([]byte(podNode.MustString()), &x); err != nil {
			return err
		}
		// update the labels with the service selctor key
		x.Spec.Selector.MatchLabels[rn.GetLabelKey()] = rn.PackageName
		x.Spec.Template.Labels[rn.GetLabelKey()] = rn.PackageName

		found := false
		vol := rn.BuildVolume()
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
				volm := rn.BuildVolumeMount()
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

		// the path must exist since we read the resource from the filesystem

		fp := fileutil.GetFullPath(rn.TargetDir, x.Annotations[kioutil.PathAnnotation])
		return fileutil.UpdateFileFromRObject(fp, x)
	}

	return nil
}

func (r *pkgConfig) validatePodContainer(fnCfg kptgenv1alpha1.Config, node *yaml.RNode) (bool, error) {

	switch fnCfg.Spec.Selector.Kind {
	case "Deployment":
		x := &appsv1.Deployment{}
		if err := sigyaml.Unmarshal([]byte(node.MustString()), &x); err != nil {
			return false, fmt.Errorf("deployment marshal Error: %s", err.Error())
		}
		for _, c := range x.Spec.Template.Spec.Containers {
			if c.Name == fnCfg.Spec.Selector.ContainerName {
				return true, nil
			}
		}
	case "StatefulSet":
		x := &appsv1.StatefulSet{}
		if err := sigyaml.Unmarshal([]byte(node.MustString()), &x); err != nil {
			return false, fmt.Errorf("statefulset marshal Error: %s", err.Error())
		}
		for _, c := range x.Spec.Template.Spec.Containers {
			if c.Name == fnCfg.Spec.Selector.ContainerName {
				return true, nil
			}
		}
	}

	return false, nil
}
