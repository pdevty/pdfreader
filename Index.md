Introduction to the pdfreader library for Go/[Issue9](https://code.google.com/p/pdfreader/issues/detail?id=9).

# Introduction #

The pdfreader library for [Go/[Issue9](https://code.google.com/p/pdfreader/issues/detail?id=9)] is a library to read contents of PDF files.

Basically it turns out that this will result in a PDF to SVG converter at the first stage.

## Details ##

PDF files are basically something that is usually at the end of some workflow and that is intended to conserve informations in a way that allows the informations to be accessed as they where intended (e.g. in terms of typographical layout). These informations need to be fetched. The project here tries to make this possible with a library for [Go/[Issue9](https://code.google.com/p/pdfreader/issues/detail?id=9)].

If you are not a [Go/[Issue9](https://code.google.com/p/pdfreader/issues/detail?id=9)]-programmer, just move away or play with an example application.

Currently everything is at it's premature state and there is no production-ready library to be expected. Well, the things work usually fine for many tasks.

If you are willing to make experiments, just checkout at http://code.google.com/p/pdfreader/source/checkout

## Basic design principles ##

- Using this library with a malformed PDF might crash the program. This is intentional.

- Keep things simple - no reason to produce billions of lines of code.

- Make the crash to be late. As late as possible. If there is something really wrong it will crash earlier or later. Why using a "safe" programming language if not using it and adding useless tests for validity of input?

- Avoid endless recursions. There are many places where this could occur in PDF-files. A fixing of [issue 226](https://code.google.com/p/pdfreader/issues/detail?id=226) with golang would help, but the gurus of Google did decide to do different. So be prepared to have no real fun with the implementation language. See [Philosophy](Philosophy.md)-page.

## Example ##

This shows an SVG displayed in Inkscape that was converted from a PDF:

![http://pdfreader.googlecode.com/files/example-convert.png](http://pdfreader.googlecode.com/files/example-convert.png)