package file

import (
	"fmt"
	"strings"
)

func ParseFlags(args []string) (map[string]string, error) {
	flags := make(map[string]string)

	for i := 1; i < len(args); i++ {
		flag := args[i]

		if !strings.HasPrefix(flag, "--") {
			return nil, fmt.Errorf("expected flag, got %q", flag)
		}

		if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
			return nil, fmt.Errorf("missing value for %s", flag)
		}

		if _, exists := flags[flag]; exists {
			return nil, fmt.Errorf("duplicate flag %s", flag)
		}

		flags[flag] = args[i+1]
		i++
	}

	return flags, nil
}

func ValidateFlags(flags map[string]string, allowed ...string) error {
	for flag := range flags {
		valid := false

		for _, name := range allowed {
			if flag == name {
				valid = true
				break
			}
		}

		if !valid {
			return fmt.Errorf("unknown flag %s", flag)
		}
	}

	return nil
}
