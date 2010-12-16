package hubbard

import "archive/tar"
import "compress/gzip"
import "http"
import "io/ioutil"
import "os"
import "strings"

func cmdRetrieve() {
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
			project := fields[0]
      sha1 := components[1]
			err := retrieve(project, sha1)
			if err != nil {
				panic(err)
			}
		}
	}
}

func retrieve(project string, sha1 string) os.Error {
  resp, _, err := http.Get("http://localhost:4788/packages/" + project + "/" + sha1 + ".tar.gz")
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
  gunzip, err := gzip.NewReader(resp.Body)
  if err != nil {
    return err
  }
  tarrer := tar.NewReader(gunzip)
  for {
    buf := make([]byte, 1024)
    for {
      count, err := tarrer.Read(buf)
      if err == os.EOF && count == 0 {
        break
      } else if err != nil {
        return err
      }
    }
    header, err := tarrer.Next()
    if err != nil {
      return err
    } else if header == nil {
      break
    }
  }
	return nil
}
