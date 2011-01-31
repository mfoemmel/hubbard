include $(GOROOT)/src/Make.inc

TARG=hubbard
GOFILES=\
	build.go\
	exec.go\
	git.go\
	hg.go\
	io.go\
	html.go\
	http.go\
	hubbard.go\
	project.go\
	repo.go\
	resolve.go\
	retrieve.go\
    tar.go\

CLEANFILES+=hub

hub: hub.$(O)
	$(LD) -L _obj -o hub hub.$(O)

hub.$(O): main.go _obj/hubbard.a
	$(GC) -I_obj -o $@ main.go

$(GOBIN)/hub: hub
	cp hub $@

INSTALLFILES+=$(GOBIN)/hub
CLEANFILES+=$(GOBIN)/hub

include $(GOROOT)/src/Make.pkg
