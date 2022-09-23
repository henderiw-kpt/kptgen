package pkgconfig

import (
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
)

func (r *pkgConfig) Deploy() error {
	// first deploy the pods
	pods, ok := r.fc[kptgenv1alpha1.FnPodKind]
	if ok {
		for _, pod := range pods {
			r.deployPod(pod)
		}
	}

	// reinitialize the package to ensure the pod resources
	// can be used to resolve some info like volums and services.
	if err := r.initializePackage(); err != nil {
		return err
	}

	// render the next resources
	for kind, nodes := range r.fc {
		for _, node := range nodes {
			switch kind {
			case kptgenv1alpha1.FnPodKind:
				// do nothing as this was already handled
			default:
				// deploy the specific function
				if err := r.supportedKinds[kind](node); err != nil {
					return err
				}
			}
		}

	}
	return nil
}
