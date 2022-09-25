package pkgconfig

import (
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
)

func (r *pkgConfig) Deploy() error {
	// print resources before rendering
	//r.resources.Print()
	// first deploy the fnConfig kind pods if they are required to be rendered
	pods, ok := r.fc[kptgenv1alpha1.FnPodKind]
	if ok {
		for _, pod := range pods {
			r.deployPod(pod)
		}
	}

	// render the next kinds
	for kind, nodes := range r.fc {
		for _, node := range nodes {
			switch kind {
			case kptgenv1alpha1.FnPodKind:
				// do nothing as this was already handled
			default:
				// deploy the specific kind function/method
				if err := r.supportedKinds[kind](node); err != nil {
					return err
				}
			}
		}
	}

	// print resources after rendering
	//r.resources.Print()
	return r.resources.Write(r.targetDir)
}
