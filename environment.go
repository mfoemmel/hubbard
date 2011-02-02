package hubbard

import "fmt"
import "os"
import "path"
import "strings"

// Environment is a map containing environment variables. VAR:VALUE
type Environment map[string]string

func (e Environment) exists(variable string) bool {
	_, present := e[variable]
	return present
}

// Converts an Environment into a []string which is what os.Environ() outputs,
// and what os.Exec() expects.
func (e Environment) toEnviron() []string {
	newEnv := []string(nil)
	for key, value := range e {
		if key == "" {
			continue
		}
		envEntry := fmt.Sprintf("%s=%s", key, value)
		newEnv = append(newEnv, envEntry)
	}
	return newEnv
}

// Given an environment from os.Environ(), transform it into an Environment.
func newEnvironment(environment []string) Environment {
	newEnv := make(Environment, 100)
	for _, pair := range environment {
		components := strings.Split(pair, "=", 2)
		key := components[0]
		value := components[1]
		newEnv[key] = value
	}
	return newEnv
}

/*
Parses the [environment] section of a hubFile.
hubFile is a path to the hubFile (usually named 'package.hub')

Contents of [environment] looks like this:
[environment]
path.PYTHONPATH=relative/path/to/dir
var.CFLAGS=-myflag
*/
func parseEnvironment(hubFile string, env Environment) (hubEnv Environment) {
	println("DEBUG *** parseEnvironment processing", hubFile)
	// The default is to return the existing environment
	hubEnv = env
	lines, err := readFileLines(hubFile)
	if err != nil {
		panic(err)
	}

	// Flag to indicate we're in the [environment] section or not.
	inEnvironment := false

	// Loop through the lines in package.hub
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") {
			inEnvironment = line == "[environment]"
		} else if inEnvironment {
			if strings.HasPrefix(line, "path.") {
				baseDir := getBaseDir(hubFile)
				if baseDir == "." || baseDir == "./" {
					baseDir, err = os.Getwd()
					if err != nil {
						panic(err)
					}
				}
				// pathExpr is a path expression of VARIABLE=path/to/foo
				pathExpr := strings.Split(line, "path.", 2)[1]
				components := strings.Split(pathExpr, "=", 2)
				envVar, envPath := components[0], components[1]
				envPath = path.Join(baseDir, envPath)
				if hubEnv.exists(envVar) {
					// prepend the new value
					hubEnv[envVar] = fmt.Sprintf("%s:%s", envPath, hubEnv[envVar])
				} else {
					hubEnv[envVar] = envPath
				}
			} else if strings.HasPrefix(line, "var.") {
				// varExpr is a variable expression of VARIABLE=value
				varExpr := strings.Split(line, "var.", 2)[1]
				components := strings.Split(varExpr, "=", 2)
				envVar, envVal := components[0], components[1]
				if hubEnv.exists(envVar) {
					// prepend the new value
					hubEnv[envVar] = fmt.Sprintf("%s %s", envVal, hubEnv[envVar])
				} else {
					hubEnv[envVar] = envVal
				}
			}
		}
	}
	return hubEnv
}
