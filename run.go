package hubbard

import "fmt"
import "os"
import "path"

func cmdRun() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("DEBUG *** os.Args: ", os.Args)
	// Assumes os.Args[1] is "run"
	cmd := findExe(os.Args[2])
	cmdArgs := os.Args[2:]
	cmdEnv := newEnvironment(os.Environ())
	// Process dependencies.
	depHubFiles := getProjectDeps("./package.hub")
	for _, hubFile := range depHubFiles {
		cmdEnv = parseEnvironment(hubFile, cmdEnv)
	}
	err := os.Exec(cmd, cmdArgs, cmdEnv.toEnviron())
	if err != nil {
		panic(err)
	}
}

// Find project dependencies listed in the [versions] section of 'package.hub'.
// Returns a list containing paths to package.hub files for each project dependency.
func getProjectDeps(hubFile string) []string {
	projectHubFiles := []string(nil)
	versions := parseVersions(hubFile)
	for _, version := range versions {
		projectHubFile := path.Join(getDepsDirFor(version.project), "package.hub")
		projectHubFiles = append(projectHubFiles, projectHubFile)
	}
	return projectHubFiles
}
