package pkgconfig

import (
	"fmt"
	"strings"

	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/util/config"
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	"github.com/henderiw-kpt/kptgen/internal/util/pkgutil"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type PkgConfig interface {
	Deploy() error
}

type pkgConfig struct {
	//fnConfigDir string
	targetDir      string
	pb             *kio.PackageBuffer
	kptFile        *yaml.RNode
	fc             map[string][]*yaml.RNode
	supportedKinds map[string]func(node *yaml.RNode) error
}

func New(args []string, fnConfig string) (PkgConfig, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("TARGET_DIR is required, positional arguments; %d provided", len(args))
	}

	r := &pkgConfig{
		targetDir: args[0],
		fc:        map[string][]*yaml.RNode{},
	}
	r.supportedKinds = map[string]func(node *yaml.RNode) error{
		kptgenv1alpha1.DummyFnConfig:     r.deployNamespace,
		kptgenv1alpha1.FnPodKind:         r.deployPod,
		kptgenv1alpha1.FnClusterRoleKind: r.deployClusterRole,
		kptgenv1alpha1.FnConfigKind:      r.deployConfig,
	}

	if err := fileutil.EnsureDir(r.targetDir, true); err != nil {
		return nil, err
	}

	if err := r.initializeFnConfig(fnConfig); err != nil {
		return nil, err
	}

	// package should be initialized first as namespace
	// takes kptfile as a dummy yaml.RNode
	if err := r.initializePackage(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *pkgConfig) initializeFnConfig(fnConfig string) error {
	// kptgenv1alpha1.DummyFnConfig is used for namespace as it does not have a fnConfig
	if fnConfig == kptgenv1alpha1.DummyFnConfig {
		r.fc[kptgenv1alpha1.DummyFnConfig] = []*yaml.RNode{
			r.kptFile,
		}
		return nil
	}
	/* READ THE CONFIG FILES*/
	// read only yml, yaml files and Kptfile
	match := []string{"*.yaml", "*.yml"}
	fcpb, err := pkgutil.GetPackage(fnConfig, match)
	if err != nil {
		return err
	}
	for _, node := range fcpb.Nodes {
		if node.GetApiVersion() == kptgenv1alpha1.FnConfigAPIVersion {
			if r.isSupportedKind(node.GetKind()) {
				// initialize the list of fn RNodes for this kind
				if _, ok := r.fc[node.GetKind()]; !ok {
					r.fc[node.GetKind()] = make([]*yaml.RNode, 0)
				}
				r.fc[node.GetKind()] = append(r.fc[node.GetKind()], node)
			} else {
				return fmt.Errorf("unexpected fnconfig kind: got: %s, supportedKind: %s", node.GetKind(), r.supportedKindString())
			}
		}
	}
	return nil
}

func (r *pkgConfig) isSupportedKind(kind string) bool {
	if _, ok := r.supportedKinds[kind]; ok {
		return true
	}
	return false
}

func (r *pkgConfig) supportedKindString() string {
	var sb strings.Builder
	for kind := range r.supportedKinds {
		sb.WriteString(kind)
	}
	return sb.String()
}

func (r *pkgConfig) initializePackage() error {
	// read only Kptfile
	match := []string{"*.yaml", "*.yml", "Kptfile"}
	pb, err := pkgutil.GetPackage(r.targetDir, match)
	if err != nil {
		return err
	}
	r.pb = pb

	cfg := config.New(r.pb, map[string]string{
		kptv1.KptFileKind: "",
		//kptgenv1alpha1.FnClusterRoleKind: filepath.Base(r.FnConfigPath),
	})

	//fmt.Println("relative", filepath.Base(r.FnConfigPath))

	selectedNodes := cfg.Get()
	if selectedNodes[kptv1.KptFileKind] == nil {
		return fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	r.kptFile = selectedNodes[kptv1.KptFileKind]

	/*
		if selectedNodes[kptgenv1alpha1.FnClusterRoleKind] == nil {
			return fmt.Errorf("fnConfig must be provided -> add fnConfig file with apiVersion: %s, kind: %s, name: %s", kptgenv1alpha1.FnConfigAPIVersion, kind, r.FnConfigPath)
		}
		r.fnConfig = selectedNodes[kptgenv1alpha1.FnClusterRoleKind]
		//fmt.Println("fn config", r.fnConfig.MustString())

		r.fc = kptgenv1alpha1.ClusterRole{}
		if err := sigyaml.Unmarshal([]byte(r.fnConfig.MustString()), &r.fc); err != nil {
			return fmt.Errorf("fnConfig marshal Error: %s", err.Error())
		}
	*/
	return nil
}
