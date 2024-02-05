package main

import (
	"fmt"

	rufio "github.com/tinkerbell/rufio/api/v1alpha1"
	"github.com/tinkerbell/tink/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var userData = `#cloud-config

package_update: true

users:
  - name: tink
	sudo: ['ALL=(ALL) NOPASSWD:ALL']
	shell: /bin/bash
	plain_text_passwd: 'tink'
	lock_passwd: false
	ssh_authorized_keys:
	  - %v
packages:
  - openssl
runcmd:
  - sed -i 's/^PasswordAuthentication no/PasswordAuthentication yes/g' /etc/ssh/sshd_config
  - systemctl enable ssh.service
  - systemctl start ssh.service
  - systemctl disable apparmor
  - systemctl disable snapd
  - rm -f /etc/hostname
`

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func (r record) hardware() v1alpha1.Hardware {
	return v1alpha1.Hardware{
		TypeMeta: v1.TypeMeta{
			Kind:       "Hardware",
			APIVersion: "tinkerbell.org/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      r.Hostname,
			Namespace: r.Namespace,
			Labels:    map[string]string{},
		},
		Spec: v1alpha1.HardwareSpec{
			BMCRef: &corev1.TypedLocalObjectReference{
				APIGroup: stringPtr("bmc.tinkerbell.org"),
				Kind:     "Machine",
				Name:     fmt.Sprintf("bmc-%s", r.Hostname),
			},
			Interfaces: []v1alpha1.Interface{
				{
					Netboot: &v1alpha1.Netboot{
						AllowPXE:      boolPtr(true),
						AllowWorkflow: boolPtr(true),
					},
					DHCP: &v1alpha1.DHCP{
						MAC:         r.Mac,
						Hostname:    r.Hostname,
						LeaseTime:   4294967294,
						NameServers: r.Nameservers,
						Arch:        "x86_64",
						UEFI:        true,
						IP: &v1alpha1.IP{
							Address: r.IPAddress,
							Netmask: r.Netmask,
							Gateway: r.Gateway,
							Family:  4,
						},
					},
				},
			},
			Metadata: &v1alpha1.HardwareMetadata{
				Instance: &v1alpha1.MetadataInstance{
					AllowPxe:  true,
					AlwaysPxe: true,
					Hostname:  r.Hostname,
					ID:        r.Mac,
					Ips: []*v1alpha1.MetadataInstanceIP{
						{
							Address: r.IPAddress,
							Family:  4,
							Gateway: r.Gateway,
							Netmask: r.Netmask,
							Public:  true,
						},
					},
				},
				Facility: &v1alpha1.MetadataFacility{
					FacilityCode: "onprem",
					PlanSlug:     "c2.medium.x86",
				},
			},
			TinkVersion: 0,
			Disks: []v1alpha1.Disk{
				{Device: r.Disk},
			},
			UserData: stringPtr(fmt.Sprintf(userData, r.SSHPublicKey)),
		},
	}
}

func (r record) bmcMachine() rufio.Machine {
	return rufio.Machine{
		TypeMeta: v1.TypeMeta{
			Kind:       "Machine",
			APIVersion: "bmc.tinkerbell.org/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("bmc-%s", r.Hostname),
			Namespace: r.Namespace,
		},
		Spec: rufio.MachineSpec{
			Connection: rufio.Connection{
				AuthSecretRef: corev1.SecretReference{
					Name:      fmt.Sprintf("bmc-%s-creds", r.Hostname),
					Namespace: r.Namespace,
				},
				Host:        r.BMCIP,
				InsecureTLS: true,
			},
		},
	}
}

func (r record) bmcSecret() corev1.Secret {
	return corev1.Secret{
		TypeMeta: v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("bmc-%s-creds", r.Hostname),
			Namespace: r.Namespace,
		},
		Type: "kubernetes.io/basic-auth",
		Data: map[string][]byte{
			"username": []byte(r.BMCUsername),
			"password": []byte(r.BMCPassword),
		},
	}
}

func (r record) bmcJob() rufio.Job {
	return rufio.Job{
		TypeMeta: v1.TypeMeta{
			Kind:       "Job",
			APIVersion: "bmc.tinkerbell.org/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("netboot-%s", r.Hostname),
			Namespace: r.Namespace,
		},
		Spec: rufio.JobSpec{
			MachineRef: rufio.MachineRef{
				Name:      fmt.Sprintf("bmc-%s", r.Hostname),
				Namespace: r.Namespace,
			},
			Tasks: []rufio.Action{
				{
					PowerAction: rufio.PowerHardOff.Ptr(),
				},
				{
					OneTimeBootDeviceAction: &rufio.OneTimeBootDeviceAction{
						Devices: []rufio.BootDevice{
							rufio.PXE,
						},
						EFIBoot: true,
					},
				},
				{
					PowerAction: rufio.PowerOn.Ptr(),
				},
			},
		},
	}
}

func (r record) workflow() v1alpha1.Workflow {
	return v1alpha1.Workflow{
		TypeMeta: v1.TypeMeta{
			Kind:       "Workflow",
			APIVersion: "tinkerbell.org/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      r.Hostname + "-cleanup",
			Namespace: r.Namespace,
		},
		Spec: v1alpha1.WorkflowSpec{
			TemplateRef: "cleanup-flow",
			HardwareRef: r.Hostname,
			HardwareMap: map[string]string{
				"machine":       r.Mac,
				"admin_agent":   "admin-node1",
				"hardware_name": r.Hostname,
			},
		},
	}
}
