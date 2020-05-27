/* just like Unix wc 
 To run this program: 
 	flex wc.l
	gcc lex.yy.c -lfl
	./a.out
	The boy stood on the burning deck
	shelling peanuts by the peck
	^D
 */
%{
	int chars = 0;
	int words = 0;
	int lines = 0;
%}
%%
[^ \t\n\r\f\v]+ { words++; chars += strlen(yytext); }
[a-zA-Z]+ 	{ words++; chars += strlen(yytext); }
\n 			{ chars++; lines++; }
. 			{ chars++; }
%%
main(int argc, char **argv)
{
	yylex();
	printf("%8d%8d%8d\n", lines, words, chars);
}

/* 
 A flex program consists of three sections, separated by %% lines. 
 	The first section contains declarations and option settings. 
	The second section is a list of patterns and actions,
	The third section is C code that is copied to the generated scanner, usually small routines related to the code in the actions.
	In the declaration section, code inside of %{ and %} is copied through verbatim near the
		beginning of the generated C source file. In this case it just sets up variables for lines, words, and characters.
	In the second section, each pattern is at the beginning of a line, followed by the C code
		to execute when the pattern matches. The C code can be one statement or possibly a
		multiline block in braces, { }. (Each pattern must start at the beginning of the line,
		since flex considers any line that starts with whitespace to be code to be copied into
		the generated C program.)

In this program, there are only three patterns. 
	The ^ at the beginning of the character class means to match any character other than
		the ones in the class, and the + once again means to match one or more of the preceding
		patterns. This demonstrates one of flex’s strengths—it’s easy to make small changes to
		patterns and let flex worry about how they might affect the generated code.
	The first one, [a-zA-Z]+, matches a word.
		The characters in brackets, known as a character class, match any single upper- or
		lowercase letter, and the + sign means to match one or more of the preceding thing,
		which here means a string of letters or a word. The action code updates the number of
		words and characters seen. 
	In any flex action, the variable yytext is set to point to the
		input text that the pattern just matched. In this case, all we care about is how many
		characters it was so we can update the character count appropriately.
	The second pattern, \n, just matches a new line. The action updates the number of lines
		and characters.
	The final pattern is a dot, which is regex-ese for any character. (It’s similar to a ? in shell scripts.) 
		The action updates the number of characters. And that’s all the patterns we need.

The C code at the end is a main program that calls yylex(), the name that flex gives to
	the scanner routine, and then prints the results. In the absence of any other arrangements, the scanner reads from the standard input. 
*/
