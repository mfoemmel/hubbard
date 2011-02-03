package hubbard

import "strings"

type Version struct {
	project string
	ref     string
	sha1    string
}

// Parse the [versions] section of hubFile (usually 'package.hub').
// hubFile is a path to the hubFile.
func parseVersions(hubFile string) []*Version {
	lines, err := readFileLines(hubFile)
	if err != nil {
		panic(err)
	}

	versions := []*Version(nil)

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
			project := fields[0]
			ref := fields[1]
			sha1 := components[1]
			versions = append(versions, &Version{project, ref, sha1})
		}
	}
	return versions
}
