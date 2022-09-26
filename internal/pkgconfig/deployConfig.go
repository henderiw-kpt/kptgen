package pkgconfig

import (
	"fmt"
	"strings"

	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/resource"
	appsv1 "k8s.io/api/apps/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	sigyaml "sigs.k8s.io/yaml"
)

func (r *pkgConfig) deployConfig(node *yaml.RNode) error {
	// marshal the fnConfig
	fnCfg := kptgenv1alpha1.Config{}
	if err := sigyaml.Unmarshal([]byte(node.MustString()), &fnCfg); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	rn := &resource.Resource{
		Kind:             kptgenv1alpha1.FnConfigKind,
		PackageName:      r.kptFile.GetName(),
		PodName:          fnCfg.Spec.Selector.Name,
		Name:             node.GetName(),
		Namespace:        r.kptFile.GetNamespace(),
		TargetDir:        r.targetDir,
		SubDir:           node.GetName(),
		ResourceNameKind: resource.NameKindFull,
	}

	// get the crds in the package as they will be used for webhooks if they are
	// required in the fnConfig
	// also get the selected yaml.RNode based on the selector, which will be used
	// later for updating volumes and labels
	crdObjects := make([]extv1.CustomResourceDefinition, 0)
	var selectedNode *yaml.RNode
	for _, node := range r.pkgResources.Get() {
		switch node.GetKind() {
		case "CustomResourceDefinition":
			crd := extv1.CustomResourceDefinition{}
			if err := sigyaml.Unmarshal([]byte(node.MustString()), &crd); err != nil {
				return err
			}
			crdObjects = append(crdObjects, crd)
		case fnCfg.Spec.Selector.Kind:
			resid := resid.FromRNode(node)
			fnCfg.Spec.Selector.Name = strings.Join([]string{r.kptFile.GetName(), fnCfg.Spec.Selector.Name}, "-")
			if resid.IsSelectedBy(fnCfg.Spec.Selector.ResId) {
				if found, err := r.validatePodContainer(fnCfg, node); err != nil {
					return err
				} else {
					if found {
						// TBD what todo if already found, we would expect 1 container that matches
						selectedNode = node
					}
				}
			}
		}
	}
	// if the selected node is not found we stop
	if selectedNode == nil {
		return fmt.Errorf("container pod not found")
	}

	// TODO
	// ClusterRoles -> to add a clusterrole bonding
	// PermissionRequests -> Todo
	// Containers -> Todo

	// render service
	for _, service := range fnCfg.Spec.Services {
		node, err := rn.RenderService(service, crdObjects)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
	}
	// render webhook
	if fnCfg.Spec.Webhook {
		node, err := rn.RenderMutatingWebhook(fnCfg.Spec, crdObjects)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
		node, err = rn.RenderValidatingWebhook(fnCfg.Spec, crdObjects)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
	}
	// render certificate
	if fnCfg.Spec.Certificate.IssuerRef != "" {
		node, err := rn.RenderCertificate(fnCfg.Spec, crdObjects)
		if err != nil {
			return err
		}
		r.pkgResources.Add(node)
	}

	// mutate deployment or statefulset
	switch fnCfg.Spec.Selector.Kind {
	case "Deployment":
		updateNode, err := rn.UpdateDeployment(node.GetName(), fnCfg, selectedNode)
		if err != nil {
			return err
		}
		r.pkgResources.Add(updateNode)

	case "StatefulSet":
		updateNode, err := rn.UpdateStatefulSet(node.GetName(), fnCfg, selectedNode)
		if err != nil {
			return err
		}
		r.pkgResources.Add(updateNode)
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
