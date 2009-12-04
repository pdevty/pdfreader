include $(GOROOT)/src/Make.$(GOARCH)

FMT = gofmt -spaces -tabwidth=2
ALLGO = graf.go
ALL = $(ALLGO) pdtosvg pdtest
PIGGY = *.$O DEADJOE

all: $(ALL)

%: %.$O
	$(LD) -o $* $*.$O

%.$O: %.go
	$(GC) -I. $*.go

%.go: %.in
	perl $*.in | $(FMT) >$*.go

depend: $(ALLGO)
	./mkdepend *.go <Makefile >mkf.new && \
	mv -f Makefile Makefile~ && \
	mv -f mkf.new Makefile

fmt:
	for a in *.go ; do \
	  $(FMT) $$a >$$a.new && mv $$a $$a~ && mv $$a.new $$a ; \
	done

clean:
	-rm *~

distclean: clean
	-rm $(ALL) $(PIGGY)

# -- depends --
graf.$O: fancy.$O pdfread.$O strm.$O util.$O
lzw.$O: crush.$O
pdfread.$O: fancy.$O hex.$O lzw.$O
pdtest.$O: fancy.$O pdfread.$O svg.$O svgtext.$O util.$O
pdtosvg.$O: fancy.$O pdfread.$O svg.$O svgtext.$O util.$O
svg.$O: graf.$O strm.$O util.$O
svgtext.$O: graf.$O pdfread.$O strm.$O util.$O
