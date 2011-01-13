package hubbard

import "archive/tar"
import "compress/gzip"
import "http"
import "io"
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
  gunzip, err := gzip.NewReader(resp.Body)
  if err != nil {
    return err
  }
  tr := tar.NewReader(gunzip)
  for {
    header, err := tr.Next()
    if err == os.EOF {
      // end of tar archive
      break
    } else if err != nil {
      return err
    } else if header == nil {
      // end of tar archive
      break
    }

    mode := uint32(header.Mode)

    // What are we unpacking? A file or a directory?
    // TODO: Handle hard links and symlinks.
    // Directory
    if header.Typeflag == '5' {
      err = os.MkdirAll(header.Name, mode)
      if err != nil {
        return err
      }
    }
    // File
    // '0' or ASCII NULL
    if header.Typeflag == '0' || header.Typeflag == 0 {
      dst, err := os.Open(header.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
      if err != nil {
        return err
      }
      _, err = io.Copy(dst, tr)
      if err != nil {
        return err
      }
    }
  }
	return nil
}
