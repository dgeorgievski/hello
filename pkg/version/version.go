package version

import "fmt"

var VERSION = "1.1.8"
var COMMIT = "none"
var DATE = "none"

func GetVersion() map[string]string {
	return map[string]string{
		"VERSION": VERSION,
		"COMMIT":  COMMIT,
		"DATE":    DATE,
	}
}

func GetPrettyVersion() string {
	helloVer := `hello version
  version: %s
  commit: %s
	`
	return fmt.Sprintf(helloVer, VERSION, COMMIT)
}
