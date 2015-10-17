# Philosophy #

It's as of today (end of 2009) sometimes a little bit strange to write larger programs with Go. The language and the libraries are still in development and they may change at every point of a day. The thing is more a moving target than something you can build cathedrals on. Does it sound bad enough for you? It might become even worse if you insist on the chance to get something changed in Go that fits your predeclared mind about what a programming language is or how it has to be. The gods of Go will follow their own ways.

For my mind they crippled recursion - especially recursion of local/unnamed functions. I guess they recognized the problem but they did not manage to jump to the attitude to add a new keyword that would be needed to support recursion for unnamed functions (How else to call a function that has no name?).

Other people do have other issues. So some people are not willing to type ";" (semicolon) at end of statements. Well, the Go-gods decided to remove the need for semicolons more or less. One could say this is minimalism, but even a lot earlier they decided to not have something like
```
  a = b < 0 ? "less zero" : "greater/equal zero"
```
This is C-like and the same statement now needs several lines of code to be done. Others than me would cry about the infix- and postfix-operators that are now statements, but these are really not that useful in my mind.

So well, the things are interesting.

I do use
  * code generators written in Perl
  * an autoimporter for needed packages
  * a dependency-finder for the Makefiles
to make work with the sources more suitable. I hope that at some day I do not need the code generators in some special cases (as I use they now, they will be useful forever ;) ) and that I do not need the autoimporter anymore. The autoimporter would be better located inside the compiler. How these two things are organized will change what is needed for the dependency-finder. Until a really good solution is out there, I'd say the things should stay the same as they are now (Dec. 22 2009 update).

**Please comment if you have some ideas.** I'll update this page as far as I have time.