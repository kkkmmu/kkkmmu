%option noyywrap
%x IFILE
%{
	struct bufstack {
		struct bufstack *prev; /* previous entry */
		YY_BUFFER_STATE bs; /* saved buffer */
		int lineno; /* saved line number */
		char *filename; /* name of this file */
		FILE *f; /* current file */
	} *curbs = 0;
	char *curfilename; /* name of current input file */
	int newfile(char *fn);
	int popfile(void);
%}
%%
^"#"[ \t]*include[ \t]*\[\"<] 				{ BEGIN IFILE; }
<IFILE>[^ \t\n\">]+ 						{
												{ int c;
													while((c = input()) && c != '\n') ;
												}
												yylineno++;
												if(!newfile(yytext))
													yyterminate(); /* no such file */
												BEGIN INITIAL;
											}
<IFILE>.|\n 								{ fprintf(stderr, "%4d bad include line\n", yylineno);
											  yyterminate();
											}
<<EOF>> 									{ if(!popfile()) yyterminate(); }
^. 											{ fprintf(yyout, "%4d %s", yylineno, yytext); }
^\n 										{ fprintf(yyout, "%4d %s", yylineno++, yytext); }
\n 											{ ECHO; yylineno++; }
. 											{ ECHO; }
%%
main(int argc, char **argv)
{
	if(argc < 2) {
		fprintf(stderr, "need filename\n");
		return 1;
	}
	if(newfile(argv[1]))
		yylex();
}
int
newfile(char *fn)
{
	FILE *f = fopen(fn, "r");
	struct bufstack *bs = malloc(sizeof(struct bufstack));
	/* die if no file or no room */
	if(!f) { perror(fn); return 0; }
	if(!bs) { perror("malloc"); exit(1); }
	/* remember state */
	if(curbs)curbs->lineno = yylineno;
	bs->prev = curbs;
	/* set up current entry */
	bs->bs = yy_create_buffer(f, YY_BUF_SIZE);
	bs->f = f;
	bs->filename = fn;
	yy_switch_to_buffer(bs->bs);
	curbs = bs;
	yylineno = 1;
	curfilename = fn;
	return 1;
}
int
popfile(void)
{
	struct bufstack *bs = curbs;
	struct bufstack *prevbs;
	if(!bs) return 0;
	/* get rid of current entry */
	fclose(bs->f);
	yy_delete_buffer(bs->bs);
	   /* switch back to previous */
	prevbs = bs->prev;
	free(bs);
	if(!prevbs) return 0;
	yy_switch_to_buffer(prevbs->bs);
	curbs = prevbs;
	yylineno = curbs->lineno;
	curfilename = curbs->filename;
	return 1;
}


/*
   In the patterns, the first pattern matches a #include statement up through the double
   quote that precedes the filename. The pattern permits optional whitespace in the usual
   places. It switches to IFILE state to read the next input filename. In IFILE state, the
   second pattern matches a filename, characters up to a closing quote, whitespace, or
   end-of-line. The filename is passed to newfile to stack the current input file and set up
   the next level of input, but first there’s the matter of dealing with whatever remains of
   the #include line. One possibility would be to use another start state and patterns that
   absorb the rest of the line, but that would be tricky, since the action switches to the
   included file, so the start state and pattern would have to be used after the end of the
   included file. Instead, this is one of the few places where input() makes a scanner
   simpler. A short loop reads until it finds the \n at the end of the line or EOF. Then,
   when scanning returns to this file after the end of the included one, it resumes at the
   beginning of the next line.

   Since an exclusive start state in effect defines its own mini-scanner, that scanner has to
   be prepared for any possible input. The next pattern deals with the case of an ill-formed
	#include line that doesn’t have a filename after the double quote. It simply prints an
	error message and uses the macro yyterminate(), which immediately returns from the
	scanner.§ This definition of #include is fairly casual and makes no effort to verify that
	the punctuation around the filename matches or that there isn’t extra junk after the
	filename. It’s not hard to write code to check those issues and diagnose errors, and a
	more polished version of this program should do so.‖

   Next is the special pattern <<EOF>>, which matches at the end of each input file. We
   call popfile(), defined later, to return to the previous input file. If it returns 0, meaning
   that was the last file, we terminate. Otherwise, the scanner will resume reading the
   previous file when it resumes scanning.

   The last four patterns do the actual work of printing out each line with a preceding line
   number. Flex provides a variable called yylineno that is intended to track line numbers,
   so we might as well use it. The pattern ^. matches any character at the beginning of a
   line, so the action prints the current line number and the character. Since a dot doesn’t
   match a newline, ^\n matches a newline at the beginning of a line, that is, an empty
   line, so the code prints out the line number and the new line and increments the line
   number. A newline or other character not at the beginning of the line is just printed
   out with ECHO, incrementing the line number for a new line.

   The routine newfile(fn) prepares to read from the file named fn, saving any previous
   input file. It does so by keeping a linked list of bufstack structures, each of which has
   a link to the previous bufstack along with the saved yylineno and filename. It opens
   the file; creates and switches to a flex buffer; and saves the previous open file, filename,
   and buffer. (In this program nothing uses the filename after the file is open, but we’ll
   reuse this code later in this chapter in a program that does.)
   The routine popfile undoes what newfile did. It closes the open file, deletes the current
   flex buffer, and then restores the buffer, filename, and line number from the prior stack
   entry. Note that it doesn’t call yyrestart() when it restores the prior buffer; if it did, it
   would lose any input that had already been read into the buffer.
  */
