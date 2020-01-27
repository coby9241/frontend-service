package utils

import "os"

// EnvPair contains the key value to be set for the test
type EnvPair struct {
	Key, Value string
}

// SetTestEnv sets the env variables for testing background ticket in Flex.
func SetTestEnv(environ []EnvPair) func() {
	for _, v := range environ {
		old := os.Getenv(v.Key)
		os.Setenv(v.Key, v.Value)
		v.Value = old
	}
	return func() { // Restore old environment after the test completes.
		for _, v := range environ {
			if v.Value == "" {
				os.Unsetenv(v.Key)
				continue
			}
			os.Setenv(v.Key, v.Value)
		}
	}
}
