package config

import (
	kptv1 "github.com/GoogleContainerTools/kpt/pkg/api/kptfile/v1"
	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	"github.com/henderiw-nephio/kptgen/internal/util/resid"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Config struct {
	Pb        *kio.PackageBuffer
	Selectors []resid.ResId
}

func (c *Config) Get() []*yaml.RNode {
	results := make([]*yaml.RNode, len(c.Selectors))
	for _, node := range c.Pb.Nodes {
		// check local config flags
		if v, ok := node.GetAnnotations()[filters.LocalConfigAnnotation]; ok && v == "true" {
			for i, selector := range c.Selectors {
				resId := resid.FromRNode(node)
				if resId.IsSelectedBy(selector) {
					results[i] = node
				}
			}
		}
	}
	return results
}

func New(pb *kio.PackageBuffer, selectors map[string]string) Config {
	c := Config{
		Pb:        pb,
		Selectors: make([]resid.ResId, 0, len(selectors)),
	}
	for kind, fileName := range selectors {
		if kind == kptv1.KptFileKind {
			c.Selectors = append(c.Selectors,
				resid.ResId{
					Gvk: resid.Gvk{
						Group:   kptv1.KptFileGroup,
						Version: kptv1.KptFileVersion,
						Kind:    kptv1.KptFileKind}},
			)
		} else {
			c.Selectors = append(c.Selectors,
				resid.ResId{
					Gvk: resid.Gvk{
						Group:   kptgenv1alpha1.FnConfigGroup,
						Version: kptgenv1alpha1.FnConfigVersion,
						Kind:    kind},
					FileName: fileName,
				},
			)
		}
	}
	return c
}
