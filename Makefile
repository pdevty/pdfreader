include $(GOROOT)/src/Make.$(GOARCH)

FMT = gofmt -spaces -tabwidth=2
ALL = graf.go pdtest
PIGGY = *.$O DEADJOE

all: $(ALL)

%: %.$O
	$(LD) -o $* $*.$O

%.$O: %.go
	$(GC) -I. $*.go

%.go: %.in
	perl $*.in | $(FMT) >$*.go

fmt:
	for a in *.go ; do \
	  $(FMT) $$a >$$a.new && mv $$a $$a~ && mv $$a.new $$a ; \
	done

clean:
	-rm *~

distclean: clean
	-rm $(ALL) $(PIGGY)

pdfread.$O: fancy.$O hex.$O lzw.$O
pdtest: pdfread.$O svg.$O fancy.$O
lzw.$O: crush.$O
graf.$O: util.$O fancy.$O pdfread.$O
svg.$O: util.$O graf.$O
