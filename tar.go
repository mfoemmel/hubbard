package hubbard

import "archive/tar"
import "compress/gzip"
import "fmt"
import "io"
import "os"
import "path"

// Creates a tarball of src to target.
func archive(path string, target string) (err os.Error) {
	defer func() os.Error {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", r)
			return err
		}
		return nil
	}()

	out, err := os.Open(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	zout, err := gzip.NewWriter(out)
	if err != nil {
		panic(err)
	}
	defer zout.Close()

	archive := tar.NewWriter(zout)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	err = copyToArchive(path, "", archive)
	if err != nil {
		panic(err)
	}
	return nil
}

// Unpacks tarball into prefix.
// Prefix is the destination for the unarchived files.
func unarchive(prefix string, tarball io.Reader) os.Error {
	gunzip, err := gzip.NewReader(tarball)
	if err != nil {
		return err
	}
	defer gunzip.Close()

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

		name := path.Join(prefix, header.Name)
		mode := uint32(header.Mode)

		// What are we unpacking? A file or a directory?
		// TODO: Handle hard links and symlinks.
		switch header.Typeflag {
		case '5': // Directory
			err = os.MkdirAll(name, mode)
			if err != nil {
				return err
			}
		case '0', 0: // File. '0' or ASCII NULL
			err = copyToFile(name, tr, mode)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy data from a tar.Reader into a destination file.
// "name" is a path to the destination file.
func copyToFile(name string, tr *tar.Reader, mode uint32) os.Error {
	dst, err := os.Open(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, tr)
	if err != nil {
		return err
	}
	return nil
}

func copyToArchive(basedir string, path string, archive *tar.Writer) os.Error {
	// Set up directories to exclude from archives.
	sourceCodeDirs := []string{".hg", ".git", ".svn", ".bzr"}
	depsDir := "deps"
	excludes := append(sourceCodeDirs, depsDir)

	if !fileExists(basedir + path) {
		return os.NewError("file not found")
	}

	if isDirectory(basedir + path) {
		if path != "" {
			header, err := createTarHeader(basedir, path+"/")
			header.Typeflag = tar.TypeDir
			header.Size = 0
			if err != nil {
				panic(err)
			}
			err = archive.WriteHeader(header)
			if err != nil {
				panic(err)
			}
		}

		children, err := list(basedir + path)
		if err != nil {
			panic(err)
		}
		for _, child := range children {
			// Don't archive excluded items.
			if contains(excludes, child) {
				continue
			}
			copyToArchive(basedir, path+"/"+child, archive)
		}
	} else {
		sourceFile, err := openReader(basedir + "/" + path)
		if err != nil {
			panic(err)
		}

		header, err := createTarHeader(basedir, path)
		if err != nil {
			panic(err)
		}
		err = archive.WriteHeader(header)
		if err != nil {
			panic(err)
		}

		io.Copy(archive, sourceFile)
	}
	return nil
}

func createTarHeader(basedir string, path string) (*tar.Header, os.Error) {
	info, err := os.Stat(basedir + "/" + path)
	if err != nil {
		return nil, err
	}
	header := new(tar.Header)
	header.Name = path[1:]
	header.Size = info.Size
	header.Mode = int64(info.Mode)
	header.Ctime = info.Ctime_ns / (1000 * 1000 * 1000)
	header.Mtime = info.Mtime_ns / (1000 * 1000 * 1000)
	header.Atime = info.Atime_ns / (1000 * 1000 * 1000)
	return header, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if pathErr, ok := err.(*os.PathError); ok {
		if pathErr.Error == os.ENOENT {
			return false
		}
	}
	panic(err)
}

func mkdir(path string) os.Error {
	return os.Mkdir(path, 0755)
}

func mkdir_p(path string) os.Error {
	return os.MkdirAll(path, 0755)
}

func list(path string) ([]string, os.Error) {
	dir, err := os.Open(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return dir.Readdirnames(-1)
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if pathErr, ok := err.(*os.PathError); ok && pathErr.Error == os.ENOENT {
			return false
		} else {
			panic(err)
		}
	}
	return info.IsDirectory()
}

func openReader(path string) (io.ReadCloser, os.Error) {
	return os.Open(path, os.O_RDONLY, 0)
}

func openWriter(path string) (io.WriteCloser, os.Error) {
	return os.Open(path, os.O_CREAT|os.O_WRONLY, 0666)
}

// Does slice s contain item?
func contains(s []string, item string) bool {
	for _, x := range s {
		if x == item {
			return true
		}
	}
	return false
}
