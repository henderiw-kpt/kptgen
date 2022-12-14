package config

import (
	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	"github.com/henderiw-kpt/kptgen/internal/pkgresource"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Config struct {
	PkgResources pkgresource.Resources
	Selectors    map[string]resid.ResId
}

func (c *Config) Get() map[string]*yaml.RNode {
	//fmt.Println("Config get selectors len", len(c.Selectors))
	//fmt.Println("Config get input selectors", c.Selectors)
	results := make(map[string]*yaml.RNode, len(c.Selectors))
	for _, node := range c.PkgResources.Get() {
		// check local config flags
		if v, ok := node.GetAnnotations()[filters.LocalConfigAnnotation]; ok && v == "true" {
			for kind, selector := range c.Selectors {
				resId := resid.FromRNode(node)
				//fmt.Println("resId", resId, resId.IsSelectedBy(selector))
				if resId.IsSelectedBy(selector) {
					results[kind] = node
				}
			}
		}
	}
	return results
}

func New(r pkgresource.Resources, selectors map[string]string) Config {
	c := Config{
		PkgResources: r,
		Selectors:    make(map[string]resid.ResId, len(selectors)),
	}

	for kind := range selectors {
		if kind == kptv1.KptFileKind {
			c.Selectors[kind] = resid.ResId{
				Gvk: resid.Gvk{
					Group:   kptv1.KptFileGroup,
					Version: kptv1.KptFileVersion,
					Kind:    kptv1.KptFileKind,
				},
			}

		} else {
			c.Selectors[kind] = resid.ResId{
				Gvk: resid.Gvk{
					Group:   kptgenv1alpha1.FnConfigGroup,
					Version: kptgenv1alpha1.FnConfigVersion,
					Kind:    kind,
				},
			}
		}
	}
	return c
}
