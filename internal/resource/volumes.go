package resource

import (
	"path/filepath"
	"strings"

	"github.com/yndd/ndd-runtime/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

func (rn *Resource) BuildVolume() corev1.Volume {
	return corev1.Volume{
		Name: rn.GetResourceName(),
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  rn.GetCertificateName(),
				DefaultMode: utils.Int32Ptr(420),
			},
		},
	}
}

func (rn *Resource) BuildVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      rn.GetResourceName(),
		MountPath: filepath.Join("tmp", strings.Join([]string{"k8s", WebhookSuffix, "server"}, "-"), CertPathSuffix),
		ReadOnly:  true,
	}
}
