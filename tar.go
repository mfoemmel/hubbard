package hubbard

import "archive/tar"
import "compress/gzip"
import "io"
import "os"

func archive(path string, target string) os.Error {
	out, err := os.Open(target, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0666)
  if err != nil {
    panic(err)
  }
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
    return err
  }
  return nil
}

func copyToArchive(basedir string, path string, archive *tar.Writer) os.Error {
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
      if child == ".hg" || child == ".git" || child == ".svn" {
        continue
      }
			copyToArchive(basedir, path + "/" + child, archive)
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
  header.Uid = info.Uid
  header.Gid = info.Gid
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
	return os.Mkdir(path, 0777)
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