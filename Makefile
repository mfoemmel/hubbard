include $(GOROOT)/src/Make.inc

TARG=hubbard
GOFILES=\
	build.go\
	exec.go\
	hg.go\
	html.go\
	http.go\
	hubbard.go\
	repo.go\

hub: hub.$(O)
	$(LD) -L _obj -o hub hub.$(O)

hub.$(O): main.go _obj/hubbard.a
	$(GC) -I_obj -o $@ main.go

$(GOBIN)/hub: hub
	cp hub $@

INSTALLFILES+=$(GOBIN)/hub
CLEANFILES+=$(GOBIN)/hub

include $(GOROOT)/src/Make.pkg