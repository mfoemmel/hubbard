package hubbard

import "testing"

func TestVersions(t *testing.T) {
	versions := parseVersions("package.hub.test")
	version := versions[0]

	expectedLen := 2
	expectedProject := "project1"
	expectedRef := "1.0"
	expectedSha1 := "abc123"

	// parseVersions should return a slice of *Versions.
	if len(versions) != expectedLen {
		t.Errorf("Did not parse the number of versions correctly!\n\tExpected: %v, got: %v", expectedLen, len(versions))
	}
	if version.project != expectedProject {
		t.Errorf("Did not parse the project name correctly!\n\tExpected: %s, got: %s", expectedProject, version.project)
	}
	if version.ref != expectedRef {
		t.Errorf("Did not parse the ref name correctly!\n\tExpected: %s, got: %s", expectedRef, version.ref)
	}
	if version.sha1 != expectedSha1 {
		t.Errorf("Did not parse the sha1 name correctly!\n\tExpected: %s, got: %s", expectedSha1, version.sha1)
	}
}
