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

	// todo check for git/hg/etc
	return &hgRepo{dir}
}