package resource

import (
	"path/filepath"
	"strings"

	"github.com/yndd/ndd-runtime/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

func (rn *Resource) BuildVolume(certificate bool) corev1.Volume {
	if certificate {
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
	return corev1.Volume{
		Name: rn.GetResourceName(),
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

func (rn *Resource) BuildVolumeMount(certificate bool) corev1.VolumeMount {
	if certificate {
		return corev1.VolumeMount{
			Name:      rn.GetResourceName(),
			MountPath: filepath.Join("tmp", strings.Join([]string{"k8s", rn.Name, "server"}, "-"), CertPathSuffix),
			ReadOnly:  true,
		}
	}

	return corev1.VolumeMount{
		Name:      rn.GetResourceName(),
		MountPath: filepath.Join(rn.Name),
	}
}
