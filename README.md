# tinkerbell-generate

This tool will take a csv file of data and generate Tinkerbell custom resource objects.
The following objects are generated:
hardware, worker, workflow, and template.

## CSV Format

The CSV file should have the following columns:

| Column Name | Description | Data Type | Example |
| ----------- | ----------- | --------- | ------- |
| hostname | The hostname of the machine | string | `machine-01` |
| bmc_ip | The IP address of the BMC | string | `10.20.30.40` |
| bmc_username | The username for the BMC | string | `admin` |
| bmc_password | The password for the BMC | string | `password` |
| mac | The MAC address of the primary NIC | string | `00:11:22:33:44:55` |
| ip_address | The IP address of the primary NIC | string | `172.16.20.3` |
| netmask | The netmask of the primary NIC | string | `255.255.255.252` |
| gateway | The gateway of the primary NIC | string | `172.16.20.1` |
| nameservers | The nameservers of the primary NIC | pipe delimited string | `8.8.8.8\|1.1.1.1` |
| labels | The labels to apply to the machine | pipe delimited string | `label1\|label2` |
| disk | The disk to use for the machine | string | `/dev/sda` |
| vendor | The vendor of the machine | string | `Supermicro` |
| namespace | The namespace to use in the CR objects | string | `default` |
| ssh_pub_key | The SSH public key to use for the machine | string | `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDQ...` |

## Usage

```shell
$ tinkerbell-generate -h

Usage of tinkerbell-generate:
  -csv string
        path to csv file (default "hardware.csv")
  -disable-bmc
        disable bmc-machine.yaml, bmc-secret.yaml, bmc-job.yaml
  -disable-hardware
        disable hardware.yaml
  -disable-workflow
        disable workflow.yaml
  -location string
        location to write yaml files (default "output")
  -namespace string
        namespace to use in the generated yaml files (default "tink-system")
```

## Building

```shell
make build
```
