package hubbard

import "fmt"
import "http"
import "io/ioutil"
import "os"
import "strings"

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
		project, sha1 := parseVersions()
		if project != "" && sha1 != "" {
			println("DEBUG *** cmdRetrieve() calling retrieve() for ", project, sha1)
			retrieve(project, sha1)
		}
	}
}

// Retrieve packages for builds on the server.
// Don't do anything if there's not a 'package.hub' file, or
// no [versions] section in 'package.hub'
func srvRetrieve() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", r)
		}
	}()

	println("DEBUG *** In srvRetrieve() ...")
	// Only perform a retrieve if there is a 'package.hub' file.
	if exists("package.hub") {
		println("DEBUG *** srvRetrieve found package.hub file.")
		project, sha1 := parseVersions()
		if project != "" && sha1 != "" {
			println("DEBUG *** srvRetrieve() calling retrieve() for ", project, sha1)
			// Assuming we're already in the correct directory.
			// Was getting a 'File not found' error when downloading via http.Get().
			// So just copy files on the filesystem, instead.
			archiveFile := "/home/jramnani/work/hubbard/data/packages/" + project + "/" + sha1 + ".tar.gz"
			println("DEBUG *** srvRetrieve() unarchiving: ", archiveFile)
			println("Retrieving package: ", archiveFile)
			af, err := os.Open(archiveFile, os.O_RDONLY, 0644)
			if err != nil {
				panic(err)
			}
			unarchive(af)
		}
	}
}

func retrieve(project string, sha1 string) os.Error {
	url := "http://localhost:4788/packages/" + project + "/" + sha1 + ".tar.gz"
	println("Retrieving package: ", url)
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

	err = unarchive(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// Parse the [versions] section of 'package.hub'
func parseVersions() (project, sha1 string) {
	lines, err := readFileLines("package.hub")
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
