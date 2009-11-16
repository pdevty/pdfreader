include $(GOROOT)/src/Make.$(GOARCH)

FMT = -spaces -tabwidth=2
ALL = pdtest

all: $(ALL)

%: %.$O
	$(LD) -o $* $*.$O

%.$O: %.go
	$(GC) -I. $*.go

fmt:
	for a in *.go ; do \
	  gofmt $(FMT) $$a >$$a.new && mv $$a $$a~ && mv $$a.new $$a ; \
	done

clean:
	-rm *~

distclean: clean
	-rm $(ALL) *.$O DEADJOE

pdtest: pdfread.$O
