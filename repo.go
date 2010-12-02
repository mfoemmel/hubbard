package hubbard

import "path"

type commit struct {
	timestamp string
	sha1      string
	author    string
	comment   string
	tags      []string
}

type repo interface {
	log() <-chan *commit
}

func findRepo(name string) repo {
	dir := path.Join("data", "repos", name)

	if exists(path.Join(dir, ".hg")) {
		return &hgRepo{dir}
	} else if exists(path.Join(dir, ".git")) {
		return &gitRepo{dir}
	}

	panic("Unknown repository type")
}
