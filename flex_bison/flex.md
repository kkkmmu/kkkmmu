A flex program basically consists of a list of regexps with instructions about what to do when the input matches any of them, known as actions. A flex-generated scanner reads through its input, matching the input against all of the regexps and doing the appropriate action on each match. Flex translates all of the regexps into an efficient internal form that lets it match the input against all the patterns simultaneously, so it’s just as fast for 100 patterns as for one.

 The observant reader may ask, if a dot matches anything, won’t it also match the letters the first pattern is supposed to match? It does, but flex breaks a tie by preferring longer matches, and if two patterns match the same thing, it prefers the pattern that appears first in the flex program. This is an utter hack, but a very useful
 one we’ll see frequently.

Most programs with flex scanners use the scanner to return a stream of tokens that are handled by a parser. Each time the program needs a token, it calls yylex(), which reads a little input and returns the token. When it needs another token, it calls yylex() again.
The scanner acts as a coroutine; that is, each time it returns, it remembers where it was, and on the next call it picks up where it left off. Within the scanner, when the action code has a token ready, it just returns it as the value from yylex(). The next time the program calls yylex(), it resumes scanning with the next input characters. Conversely, if a pattern doesn’t produce a token for the calling program and doesn’t return, the scanner will just keep going within the same call to yylex(), scanning the next input characters. 

When a flex scanner returns a stream of tokens, each token actually has two parts, the token and the token’s value. The token is a small integer. The token numbers are arbitrary, except that token zero always means end-of-file. When bison creates a parser, bison assigns the token numbers automatically starting at 258 (this avoids collisions with literal character tokens, discussed later) and creates a .h with definitions of the tokens numbers. 

A token’s value identifies which of a group of similar tokens this one is. In our scanner, all numbers are NUMBER tokens, with the value saying what number it is. When parsing more complex input with names, floating-point numbers, string literals, and the like, the value says which name, number, literal, or whatever, this token is. 

Most flex programs now use %option noyywrap and provide their own main routine, so they don’t need the flex library.
