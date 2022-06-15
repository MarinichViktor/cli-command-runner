package cli

import (
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

func ParseProjectArgs() ([]*ProjectArgs, error) {
	l, _ := strconv.Atoi(os.Args[2])
	args := make([]*ProjectArgs, l)

	data, e := os.ReadFile(os.Args[1])
	if e != nil {
		return nil, e
	}

	if e := yaml.Unmarshal(data, &args); e != nil {
		return nil, e
	}

	return args, nil
}
