rm -rf data
mkdir -p data/repos
mkdir -p data/working
echo cloning go...
hg clone https://go.googlecode.com/hg/ data/repos/go
hg clone data/repos/go data/working/go
echo cloning hubbard...
git clone git://github.com/mfoemmel/hubbard.git data/repos/hubbard
(cd data/repos/hubbard && git checkout -b go origin/go)
git clone data/repos/hubbard data/working/hubbard