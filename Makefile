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

pdfread.$O: fancy.$O hex.$O lzw.$O
pdtest: pdfread.$O svg.$O fancy.$O
lzw.$O: crush.$O
graf.$O: fancy.$O pdfread.$O
svg.$O: util.$O graf.$O grafgs.$O

