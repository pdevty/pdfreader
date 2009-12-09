include $(GOROOT)/src/Make.$(GOARCH)

FMT = gofmt -spaces -tabwidth=2
ALLGO = graf.go cmapi.go
ALL = $(ALLGO) pdtosvg pdtest pdstream
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
cmap.$O: fancy.$O ps.$O strm.$O util.$O
cmapi.$O: cmapt.$O fancy.$O ps.$O util.$O
graf.$O: fancy.$O ps.$O strm.$O util.$O
lzw.$O: crush.$O
pdfread.$O: fancy.$O hex.$O lzw.$O ps.$O
pdstream.$O: cmapi.$O fancy.$O pdfread.$O util.$O
pdtest.$O: pdfread.$O
pdtosvg.$O: fancy.$O pdfread.$O strm.$O svg.$O svgtext.$O util.$O
ps.$O: fancy.$O hex.$O
svg.$O: graf.$O strm.$O util.$O
svgtext.$O: cmap.$O cmapt.$O fancy.$O graf.$O pdfread.$O ps.$O strm.$O util.$O
