package hubbard

import "fmt"
import "os"
import "path"

import "testing"

func TestNewEnv(t *testing.T) {
	testEnv := []string{"FOO=bar"}
	myEnv := newEnvironment(testEnv)
	fmt.Println("DEBUG myEnv is: ", myEnv)
	if env := newEnvironment(testEnv); env == nil {
		t.Errorf("Failed to create a new Environment!")
	}
}

func TestExists(t *testing.T) {
	testEnv := newEnvironment([]string{"FOO=bar"})
	if testEnv.exists("FOO") != true {
		t.Errorf("Failed to test whether an environment variable exists!  Expected true, got false.")
	}
}

func TestToEnviron(t *testing.T) {
	testEnv := newEnvironment([]string{"FOO=bar"})
	origEnv := testEnv.toEnviron()
	if origEnv[0] != "FOO=bar" {
		t.Errorf("Failed to convert Environment to an os.Environ() compatible type. Expected ['FOO=bar'], got: ", origEnv)
	}
}

func TestParseEnvironment(t *testing.T) {
	testCwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// expectedFoo is the expanded value of env var path.FOO
	expectedFoo := path.Join(testCwd, "bar")
	emptyEnv := make(Environment, 0)
	testFile := "./package.hub.test"
	testEnv := parseEnvironment(testFile, emptyEnv)
	if testEnv.exists("FOO") != true {
		t.Errorf("Failed parsing [environment] section of package.hub.test!\n\tExpected true, got false.")
	}
	if testEnv["FOO"] != expectedFoo {
		t.Errorf("Failed to compute path for path.FOO!\n\tExpected: %s, got %v", expectedFoo, testEnv["FOO"])
	}
	if testEnv["SPAM"] != "eggs" {
		t.Errorf("Failed to set environment variable correctly!\n\tExpected SPAM=eggs, got SPAM=%v", testEnv["SPAM"])
	}
}

func TestPrependEnvironmentVariable(t *testing.T) {
	testCwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// testFoo is the expanded value of env var path.FOO
	testFoo := path.Join(testCwd, "bar")

	populatedEnv := Environment{"FOO": "/path/to/a/dir", "SPAM": "-bacon"}

	expectedFoo := fmt.Sprintf("%s:%s", testFoo, populatedEnv["FOO"])
	expectedSpam := "eggs -bacon"

	testFile := "./package.hub.test"
	testEnv := parseEnvironment(testFile, populatedEnv)

	if testEnv.exists("FOO") != true {
		t.Errorf("Failed parsing [environment] section of package.hub.test!\n\tExpected true, got false.")
	}
	if testEnv["FOO"] != expectedFoo {
		t.Errorf("Failed to compute path for path.FOO!\n\tExpected: %s, got %v", expectedFoo, testEnv["FOO"])
	}
	if testEnv["SPAM"] != expectedSpam {
		t.Errorf("Failed to compute path for var.SPAM!\n\tExpected: %s, got %v", expectedSpam, testEnv["SPAM"])
	}
}
