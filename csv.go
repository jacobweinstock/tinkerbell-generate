package main

import (
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// You could also use the standard Stringer interface
func (n *nameservers) String() string {
	return strings.Join(*n, ",")
}

// Convert the CSV string as internal date
func (n *nameservers) UnmarshalCSV(csv string) (err error) {
	*n = strings.Split(csv, "|")
	return nil
}

// You could also use the standard Stringer interface
func (l *labels) String() string {
	return strings.Join(*l, ",")
}

// Convert the CSV string as internal date
func (l *labels) UnmarshalCSV(csv string) (err error) {
	*l = strings.Split(csv, "|")
	return nil
}

// You could also use the standard Stringer interface
func (s *sshPublicKey) String() string {
	return string(*s)
}

// Convert the CSV string as internal date
func (s *sshPublicKey) UnmarshalCSV(csv string) (err error) {
	if csv == "" {
		return nil
	}
	p, err := homedir.Expand(csv)
	if err != nil {
		return err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	*s = sshPublicKey(b)

	return nil
}
