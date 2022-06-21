package cli

import (
	"gopkg.in/yaml.v3"
	"os"
)

func ParseProjectArgs() ([]*ProjectArgs, error) {
	var args []*ProjectArgs

	data, e := os.ReadFile(os.Args[1])
	if e != nil {
		return nil, e
	}

	if e := yaml.Unmarshal(data, &args); e != nil {
		return nil, e
	}

	return args, nil
}
