/* recognize tokens for the calculator and print them out */
%%
"+" { printf("PLUS\n"); }
"-" { printf("MINUS\n"); }
"*" { printf("TIMES\n"); }
"/" { printf("DIVIDE\n"); }
"|" { printf("ABS\n"); }
[0-9]+ { printf("NUMBER %s\n", yytext); }
\n { printf("NEWLINE\n"); }
[ \t] { }
. { printf("Mystery character %s\n", yytext); }
%%

/*
The first five patterns are literal operators, written as quoted strings, and the actions,
	for now, just print a message saying what matched. The quotes tell flex to use the strings
	as is, rather than interpreting them as regular expressions.

The sixth pattern matches an integer. The bracketed pattern [0-9]matches any single
	digit, and the following + sign means to match one or more of the preceding item, which
	here means a string of one or more digits. The action prints out the string that’s matched, 
	using the pointer yytext that the scanner sets after each match.
The seventh pattern matches a newline character, represented by the usual C \n sequence.
The eighth pattern ignores whitespace. It matches any single space or tab (\t), and the empty action code does nothing.
The final pattern is the catchall to match anything the other patterns didn’t. Its action
	code prints a suitable complaint.
These nine patterns now provide rules to match anything that the user might enter. As we continue 
	to develop the calculator, we’ll add more rules to match more tokens, but
	these will do to get us started.

 */
