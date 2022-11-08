package resource

import (
	"fmt"
	"reflect"
	"strings"

	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	kptgenv1alpha1 "github.com/henderiw-kpt/kptgen/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
	"sigs.k8s.io/kustomize/kyaml/kio/kioutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func (rn *Resource) RenderCertificate(cfg, obj interface{}) (*yaml.RNode, error) {
	info, ok := cfg.(*kptgenv1alpha1.ConfigSpec)
	if !ok {
		return nil, fmt.Errorf("wrong object in rendercertificate: %v", reflect.TypeOf(cfg))
	}
	x := &certv1.Certificate{
		TypeMeta: metav1.TypeMeta{
			Kind:       certv1.CertificateKind,
			APIVersion: certv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetCertificateName(),
			Namespace: rn.Namespace,
			Labels:    rn.GetK8sLabels(),
			Annotations: map[string]string{
				kioutil.PathAnnotation:  rn.GetRelativeFilePath(certv1.CertificateKind),
				kioutil.IndexAnnotation: "0",
			},
		},
		Spec: certv1.CertificateSpec{
			DNSNames: []string{
				rn.GetDnsName(),
				rn.GetDnsName("cluster", "local"),
			},
			IssuerRef: certmetav1.ObjectReference{
				Kind: certv1.IssuerKind,
				Name: info.Certificate.IssuerRef,
			},
			SecretName: rn.GetCertificateName(),
		},
	}

	b := new(strings.Builder)
	p := printers.YAMLPrinter{}
	p.PrintObj(x, b)
	return yaml.Parse(b.String())

	//return fileutil.CreateFileFromRObject(rn.GetFilePath(certv1.CertificateKind), x)
}
