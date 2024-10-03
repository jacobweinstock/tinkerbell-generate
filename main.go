package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	"gopkg.in/yaml.v3"
)

type nameservers []string
type labels []string
type sshPublicKey string
type ymls []yml

type yml struct {
	data []byte
	name string
}

type record struct {
	BMCIP        string       `csv:"bmc_ip,omitempty"`
	BMCPassword  string       `csv:"bmc_password,omitempty"`
	BMCUsername  string       `csv:"bmc_username,omitempty"`
	Disk         string       `csv:"disk,omitempty"`
	Gateway      string       `csv:"gateway,omitempty"`
	Hostname     string       `csv:"hostname,omitempty"`
	IPAddress    string       `csv:"ip_address,omitempty"`
	Labels       labels       `csv:"labels,omitempty"`
	Mac          string       `csv:"mac,omitempty"`
	Nameservers  nameservers  `csv:"nameservers,omitempty"`
	Namespace    string       `csv:"namespace,omitempty"`
	Netmask      string       `csv:"netmask,omitempty"`
	SSHPublicKey sshPublicKey `csv:"ssh_pub_key,omitempty"`
	Vendor       string       `csv:"vendor,omitempty"`
	Template     string       `csv:"template,omitempty"`
}

type config struct {
	csvPath         string
	disableBMC      bool
	disableHardware bool
	disableWorkflow bool
	outputLocation  string
	namespace       string
	sshPublicKey    string
	template        string
}

func main() {
	fs := flag.NewFlagSet("tinkerbell-generate", flag.ExitOnError)
	cfg := config{}
	fs.StringVar(&cfg.csvPath, "csv", "hardware.csv", "path to csv file")
	fs.StringVar(&cfg.outputLocation, "location", "output", "location to write yaml files")
	fs.BoolVar(&cfg.disableHardware, "disable-hardware", false, "disable generating hardware.yaml")
	fs.BoolVar(&cfg.disableBMC, "disable-bmc", false, "disable generating bmc-machine.yaml, bmc-secret.yaml")
	fs.BoolVar(&cfg.disableWorkflow, "disable-workflow", false, "disable generating workflow.yaml")
	fs.StringVar(&cfg.namespace, "namespace", "tink-system", "namespace to use in the generated yaml files")
	fs.StringVar(&cfg.sshPublicKey, "ssh-public-key", "", "path to ssh public key")
	fs.StringVar(&cfg.template, "template", "cleanup-flow", "template to use in the generated yaml files")
	fs.Parse(os.Args[1:])

	f, err := os.Open(cfg.csvPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	records := []record{}
	if err := gocsv.UnmarshalFile(f, &records); err != nil {
		log.Fatal(err)
	}

	for _, record := range records {
		// 1. create a directory for each record
		// 2. create templates as bytes
		// 3. write byte templates to files on disk in the directory
		loc := filepath.Join(cfg.outputLocation, record.Hostname)
		if err := os.MkdirAll(loc, 0755); err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		if record.Namespace == "" {
			record.Namespace = cfg.namespace
		}
		if record.SSHPublicKey == "" && cfg.sshPublicKey != "" {
			if err := record.SSHPublicKey.UnmarshalCSV(cfg.sshPublicKey); err != nil {
				log.Fatal(err)
			}
		}
		record.Template = cfg.template
		y := createYamls(cfg.disableHardware, cfg.disableBMC, cfg.disableWorkflow, record)
		for _, yaml := range y {
			if err := os.WriteFile(filepath.Join(loc, yaml.name), yaml.data, 0644); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func createYamls(dh, db, dw bool, r record) ymls {
	ymls := ymls{}
	if !dh {
		ymls = append(ymls, yml{
			name: "hardware.yaml",
			data: marshal(r.hardware()),
		})
	}
	if !db {
		ymls = append(ymls, yml{
			name: "bmc-machine.yaml",
			data: marshal(r.bmcMachine()),
		})
		ymls = append(ymls, yml{
			name: "bmc-secret.yaml",
			data: marshal(r.bmcSecret()),
		})
	}
	if !dw {
		ymls = append(ymls, yml{
			name: "workflow.yaml",
			data: marshal(r.workflow()),
		})
	}
	return ymls
}

func marshal(h any) []byte {
	b, err := Marshal(&h)
	if err != nil {
		return []byte{}
	}

	return b
}

// Marshal the object into JSON then convert
// JSON to YAML and returns the YAML.
func Marshal(o interface{}) ([]byte, error) {
	j, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("error marshaling into JSON: %v", err)
	}

	y, err := JSONToYAML(j)
	if err != nil {
		return nil, fmt.Errorf("error converting JSON to YAML: %v", err)
	}

	return y, nil
}

// JSONToYAML Converts JSON to YAML.
func JSONToYAML(j []byte) ([]byte, error) {
	// Convert the JSON to an object.
	var jsonObj interface{}
	// We are using yaml.Unmarshal here (instead of json.Unmarshal) because the
	// Go JSON library doesn't try to pick the right number type (int, float,
	// etc.) when unmarshalling to interface{}, it just picks float64
	// universally. go-yaml does go through the effort of picking the right
	// number type, so we can preserve number type throughout this process.
	err := yaml.Unmarshal(j, &jsonObj)
	if err != nil {
		return nil, err
	}

	// Marshal this object into YAML.
	return yaml.Marshal(jsonObj)
}
