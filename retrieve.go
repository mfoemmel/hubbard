package hubbard

import "fmt"
import "http"
import "io/ioutil"
import "os"
import "path"
import "strings"

// Retrieve downloads packages listed in the [versions] section of 'package.hub'
// into the client's project directory.
func cmdRetrieve() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", r)
			os.Exit(1)
		}
	}()

	println("DEBUG *** In cmdRetrieve() ...")
	// Only perform a retrieve if there is a 'package.hub' file.
	if exists("package.hub") {
		println("DEBUG *** cmdRetrieve found package.hub file.")
		project, sha1 := parseVersions("package.hub")
		if project != "" && sha1 != "" {
			println("DEBUG *** cmdRetrieve() calling retrieve() for ", project, sha1)
			retrieve(project, sha1)
		}
	}
}

// Retrieve packages for builds on the server into destination path.
// Don't do anything if there's not a 'package.hub' file, or
// no [versions] section in 'package.hub'
func srvRetrieve(dst string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", r)
		}
	}()

	println("DEBUG *** In srvRetrieve() ...")
	// Only perform a retrieve if there is a 'package.hub' file.
	hubFile := path.Join(dst, "package.hub")
	if exists(hubFile) {
		println("DEBUG *** srvRetrieve found package.hub file.")
		project, sha1 := parseVersions(hubFile)
		retrieveDir := path.Join(dst, "deps", project)
		err := mkdir_p(retrieveDir)
		if err != nil {
			panic(err)
		}
		if project != "" && sha1 != "" {
			// Unpack the package into retrieveDir.
			archiveFile := path.Join(getPackageDir(), project, sha1) + ".tar.gz"
			println("DEBUG *** srvRetrieve() unarchiving: ", archiveFile)
			println("DEBUG *** into: ", retrieveDir)
			println("Retrieving package: ", archiveFile)
			af, err := os.Open(archiveFile, os.O_RDONLY, 0644)
			if err != nil {
				panic(err)
			}
			unarchive(retrieveDir, af)
		}
	}
}

func retrieve(project string, sha1 string) os.Error {
	url := "http://localhost:4788/packages/" + project + "/" + sha1 + ".tar.gz"
	destDir := getDepsDirFor(project)

	println("Retrieving package: ", url)
	println("\tinto: ", destDir)

	err := mkdir_p(destDir)
	if err != nil {
		return err
	}

	resp, _, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		panic(strings.TrimSpace(string(body)))
	}

	err = unarchive(destDir, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Parse the [versions] section of hubFile (usually 'package.hub').
// hubFile is a path to the hubFile.
func parseVersions(hubFile string) (project, sha1 string) {
	lines, err := readFileLines(hubFile)
	if err != nil {
		panic(err)
	}
	// Flag to indicate whether we're in the [versions] section or not
	inVersions := false

	// Loop through lines in package.hub
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") {
			inVersions = line == "[versions]"
		} else if inVersions {
			components := strings.Split(line, "=", 2)
			fields := strings.Split(components[0], "/", 2)
			project = fields[0]
			sha1 = components[1]
		}
	}
	return project, sha1
}
