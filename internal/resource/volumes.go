package resource

import (
	"path/filepath"
	"strings"

	"github.com/yndd/ndd-runtime/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

func BuildVolume(rn *Resource) corev1.Volume {
	return corev1.Volume{
		Name: rn.GetName(),
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName:  rn.GetCertificateName(),
				DefaultMode: utils.Int32Ptr(420),
			},
		},
	}
}

func BuildVolumeMount(rn *Resource) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      rn.GetName(),
		MountPath: filepath.Join("tmp", strings.Join([]string{"k8s", WebhookSuffix, "server"}, "-"), CertPathSuffix),
		ReadOnly:  true,
	}
}
