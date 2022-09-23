package webhook

import (
	"context"

	docs "github.com/henderiw-kpt/kptgen/internal/docs/generated/applydocs"
	"github.com/henderiw-kpt/kptgen/internal/pkgconfig"
	"github.com/spf13/cobra"
)

// NewRunner returns a command runner.
func NewRunner(ctx context.Context, parent string) *Runner {
	r := &Runner{
		Ctx: ctx,
	}
	c := &cobra.Command{
		Use:     "webhook TARGET_DIR [flags]",
		Args:    cobra.MaximumNArgs(1),
		Short:   docs.WebhookShort,
		Long:    docs.WebhookShort + "\n" + docs.WebhookLong,
		Example: docs.WebhookExamples,
		RunE:    r.runE,
	}

	r.Command = c
	r.Command.Flags().StringVar(
		&r.FnConfigPath, "fn-config", "", "path to the function config file")
	return r
}

func NewCommand(ctx context.Context, parent string) *cobra.Command {
	return NewRunner(ctx, parent).Command
}

type Runner struct {
	Command      *cobra.Command
	FnConfigPath string
	//TargetDir    string
	Ctx context.Context
	// dynamic input
	//pb      *kio.PackageBuffer
	//kptFile *yaml.RNode
	//fnConfig *yaml.RNode
	//fc kptgenv1alpha1.Webhook
	pkgCfg pkgconfig.PkgConfig
}

/*
type objInfo struct {
	renderFn func(interface{}, interface{}) error
}
*/

func (r *Runner) runE(c *cobra.Command, args []string) error {
	var err error
	r.pkgCfg, err = pkgconfig.New(args, r.FnConfigPath)
	if err != nil {
		return err
	}

	if err := r.pkgCfg.Deploy(); err != nil {
		return err
	}
	/*
		if err := r.validate(args, kptgenv1alpha1.FnWebhookKind); err != nil {
			return err
		}

		rn := &resource.Resource{
			Operation:      resource.WebhookSuffix,
			ControllerName: r.kptFile.GetName(),
			Name:           r.kptFile.GetName(),
			Namespace:      r.kptFile.GetNamespace(),
			TargetDir:      r.TargetDir,
			SubDir:         resource.WebhookDir,
			NameKind:       resource.NameKindController,
			PathNameKind:   resource.NameKindKind,
		}

		matchResources := map[string]*objInfo{
			"Service": {
				renderFn: rn.RenderService,
			},
			"Certificate": {
				renderFn: rn.RenderCertificate,
			},
			"MutatingWebhookConfiguration": {
				renderFn: rn.RenderMutatingWebhook,
			},
			"ValidatingWebhookConfiguration": {
				renderFn: rn.RenderValidatingWebhook,
			},
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
			case r.fc.Spec.Selector.Kind:
				if found, err := r.validatePodContainer(node); err != nil {
					return err
				} else {
					if found {
						// TBD what todo if already found, we would expect 1 container that matches
						podNode = node
					}
				}
			}
		}
		if podNode == nil {
			return fmt.Errorf("container pod not found")
		}

		for kind, objInfo := range matchResources {
			if kind == "Service" {
				for _, service := range r.fc.Spec.Services {
					if err := objInfo.renderFn(service, crdObjects); err != nil {
						return err
					}
				}
			} else {
				if err := objInfo.renderFn(r.fc.Spec, crdObjects); err != nil {
					return err
				}
			}

		}

		switch r.fc.Spec.Selector.Kind {
		case "Deployment":
			x := &appsv1.Deployment{}
			if err := sigyaml.Unmarshal([]byte(podNode.MustString()), &x); err != nil {
				return err
			}
			// update the labels with the service selctor key
			x.Spec.Selector.MatchLabels[rn.GetLabelKey()] = rn.ControllerName
			x.Spec.Template.Labels[rn.GetLabelKey()] = rn.ControllerName

			found := false
			vol := resource.BuildVolume(rn)
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
				if c.Name == r.fc.Spec.Selector.ContainerName {
					found := false
					volm := resource.BuildVolumeMount(rn)
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
			return fileutil.UpdateFileFromRObject(resource.DeploymentKind, fp, x)

		case "StatefulSet":
		}
	*/

	return nil
}

/*
func (r *Runner) validatePodContainer(node *yaml.RNode) (bool, error) {
	switch r.fc.Spec.Selector.Kind {
	case "Deployment":
		x := &appsv1.Deployment{}
		if err := sigyaml.Unmarshal([]byte(node.MustString()), &x); err != nil {
			return false, fmt.Errorf("deployment marshal Error: %s", err.Error())
		}
		for _, c := range x.Spec.Template.Spec.Containers {
			if c.Name == r.fc.Spec.Selector.ContainerName {
				return true, nil
			}
		}
	case "StatefulSet":
		x := &appsv1.StatefulSet{}
		if err := sigyaml.Unmarshal([]byte(node.MustString()), &x); err != nil {
			return false, fmt.Errorf("statefulset marshal Error: %s", err.Error())
		}
		for _, c := range x.Spec.Template.Spec.Containers {
			if c.Name == r.fc.Spec.Selector.ContainerName {
				return true, nil
			}
		}
	}

	return false, nil
}
*/

/*
func (r *Runner) validate(args []string, kind string) error {
	if len(args) < 1 {
		return fmt.Errorf("TARGET_DIR is required, positional arguments; %d provided", len(args))
	}

	r.TargetDir = args[0]

	if err := fileutil.EnsureDir("TARGET_DIR", r.TargetDir, true); err != nil {
		return err
	}

	if r.FnConfigPath == "" {
		return fmt.Errorf("a fn-config must be provided")
	}

	b, err := fileutil.ReadFile(r.FnConfigPath)
	if err != nil {
		return err
	}

	r.fc = kptgenv1alpha1.Webhook{}
	if err := sigyaml.Unmarshal(b, &r.fc); err != nil {
		return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
	}

	// read only Kptfile
	match := []string{"*.yaml", "*.yml", "Kptfile"}
	pb, err := pkgutil.GetPackage(r.TargetDir, match)
	if err != nil {
		return err
	}
	r.pb = pb

	cfg := config.New(r.pb, map[string]string{
		kptv1.KptFileKind: "",
		//kptgenv1alpha1.FnWebhookKind: filepath.Base(r.FnConfigPath),
	})

	//fmt.Println("relative", filepath.Base(r.FnConfigPath))

	selectedNodes := cfg.Get()
	if selectedNodes[kptv1.KptFileKind] == nil {
		return fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	r.kptFile = selectedNodes[kptv1.KptFileKind]


	return nil
}
*/
