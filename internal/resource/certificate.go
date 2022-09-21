package resource

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	certv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	certmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	kptgenv1alpha1 "github.com/henderiw-nephio/kptgen/api/v1alpha1"
	"github.com/henderiw-nephio/kptgen/internal/util/fileutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/printers"
)

func RenderCertificate(rn *Resource, cfg, obj interface{}) error {
	rn.Kind = certv1.CertificateKind

	info, ok := cfg.(*kptgenv1alpha1.WebhookSpec)
	if !ok {
		return fmt.Errorf("wrong object in rendercertificate: %v", reflect.TypeOf(cfg))
	}
	x := &certv1.Certificate{
		TypeMeta: metav1.TypeMeta{
			Kind:       certv1.CertificateKind,
			APIVersion: certv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rn.GetCertificateName(),
			Namespace: rn.Namespace,
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
	if err := p.PrintObj(x, b); err != nil {
		return err
	}

	if err := fileutil.EnsureDir(certv1.CertificateKind, filepath.Dir(rn.GetFilePath()), true); err != nil {
		return err
	}

	if err := ioutil.WriteFile(rn.GetFilePath(), []byte(b.String()), 0644); err != nil {
		return err
	}
	return nil
}
