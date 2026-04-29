package env

import "os"

type LookupEnvFunc func(key string) (string, bool)

func LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
