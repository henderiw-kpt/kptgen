package pkgconfig

import (
	"fmt"
	"strings"

	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/pkgresource"
	"github.com/henderiw-kpt/kptgen/internal/util/config"
	"github.com/henderiw-kpt/kptgen/internal/util/fileutil"
	"github.com/henderiw-kpt/kptgen/internal/util/pkgutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type PkgConfig interface {
	Deploy() error
}

type pkgConfig struct {
	targetDir      string
	pkgResources   pkgresource.Resources
	kptFile        *yaml.RNode
	fc             map[string][]*yaml.RNode
	supportedKinds map[string]func(node *yaml.RNode) error
}

func New(args []string, fnConfig string) (PkgConfig, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("TARGET_DIR is required, positional arguments; %d provided", len(args))
	}

	r := &pkgConfig{
		targetDir:    args[0],
		pkgResources: pkgresource.New(),
		fc:           map[string][]*yaml.RNode{},
	}
	// list of supported kind and methods to implement them
	r.supportedKinds = map[string]func(node *yaml.RNode) error{
		kptgenv1alpha1.DummyFnConfig:     r.deployNamespace,
		kptgenv1alpha1.FnPodKind:         r.deployPod,
		kptgenv1alpha1.FnClusterRoleKind: r.deployClusterRole,
		kptgenv1alpha1.FnConfigKind:      r.deployConfig,
	}

	if err := fileutil.EnsureDir(r.targetDir, true); err != nil {
		return nil, err
	}

	// initializes the fnConfigs that will be used to render the resources
	if err := r.initializeFnConfig(fnConfig); err != nil {
		return nil, err
	}

	// read the yaml files and kptfile as an initialization of the package
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
	// read only yml, yaml files
	match := []string{"*.yaml", "*.yml"}
	fcpb, err := pkgutil.GetPackage(fnConfig, match)
	if err != nil {
		return err
	}
	for _, node := range fcpb.Nodes {
		//fmt.Println("initializeFnConfig", node.GetApiVersion(), node.GetKind())
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

// validates if the fnCfongig kind is supported
func (r *pkgConfig) isSupportedKind(kind string) bool {
	if _, ok := r.supportedKinds[kind]; ok {
		return true
	}
	return false
}

// supportedKindString concatenates the supported kinds in a string
// used for printing
func (r *pkgConfig) supportedKindString() string {
	var sb strings.Builder
	for kind := range r.supportedKinds {
		sb.WriteString(kind)
	}
	return sb.String()
}

// initializePackage reads all yaml files and kptFile and initialzes the resources
func (r *pkgConfig) initializePackage() error {
	// read only Kptfile
	match := []string{"*.yaml", "*.yml", "Kptfile"}
	pb, err := pkgutil.GetPackage(r.targetDir, match)
	if err != nil {
		return err
	}

	// initialize resources
	for _, node := range pb.Nodes {
		r.pkgResources.Add(node)
	}

	cfg := config.New(r.pkgResources, map[string]string{
		kptv1.KptFileKind: "",
		//kptgenv1alpha1.FnClusterRoleKind: filepath.Base(r.FnConfigPath),
	})

	selectedNodes := cfg.Get()
	if selectedNodes[kptv1.KptFileKind] == nil {
		return fmt.Errorf("kptFile must be provided -> run kpt pkg init <DIR>")
	}
	r.kptFile = selectedNodes[kptv1.KptFileKind]

	return nil
}
