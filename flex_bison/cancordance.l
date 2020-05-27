%option noyywrap nodefault yylineno case-insensitive
/* the symbol table */
%{
	struct symbol { /* a word */
		char *name;
		struct ref *reflist;
	};
	struct ref {
		struct ref *next;
		char *filename;
		int flags;
		int lineno;
	};
	/* simple symtab of fixed size */
#define NHASH 9997
	struct symbol symtab[NHASH];
	struct symbol *lookup(char*);
	void addref(int, char*, char*,int);
	char *curfilename; /* name of current input file */
%}
%%
a |
an |
and |
are |
as |
at |
be |
but |
for |
in |
is |
it |
of |
on |
or |
that |
the |
this |
to /* ignore */         {}
[a-z]+(\'(s|t))? 		{ addref(yylineno, curfilename, yytext, 0); }
.|\n 					{/* ignore everything else */}
%%

main(argc, argv)
	int argc;
	char **argv;
{
	int i;
	if(argc < 2) { /* just read stdin */
		curfilename = "(stdin)";
		yylineno = 1;
		yylex();
	} else
		for(i = 1; i < argc; i++) {
			FILE *f = fopen(argv[i], "r");
			if(!f) {
				perror(argv[1]);
				return (1);
			}
			curfilename = argv[i]; /* for addref */
			yyrestart(f);
			yylineno = 1;
			yylex();
			fclose(f);
		}
	printrefs();
}

/* hash a symbol */
static unsigned
symhash(char *sym)
{
	unsigned int hash = 0;
	unsigned c;
	while(c = *sym++) hash = hash*9 ^ c;
	return hash;
}

struct symbol *
lookup(char* sym)
{
	struct symbol *sp = &symtab[symhash(sym)%NHASH];
	int scount = NHASH; /* how many have we looked at */
	while(--scount >= 0) {
		if(sp->name && !strcmp(sp->name, sym)) return sp;
		if(!sp->name) { /* new entry */
			sp->name = strdup(sym);
			sp->reflist = 0;
			return sp;
		}
		if(++sp >= symtab+NHASH) sp = symtab; /* try the next entry */
	}
	fputs("symbol table overflow\n", stderr);
	abort(); /* tried them all, table is full */
}

void
addref(int lineno, char *filename, char *word, int flags)
{
	struct ref *r;
	struct symbol *sp = lookup(word);
	/* don't do dups of same line and file */
	if(sp->reflist &&
			sp->reflist->lineno == lineno &&
			sp->reflist->filename == filename) return;
	r = malloc(sizeof(struct ref));
	if(!r) {fputs("out of space\n", stderr); abort(); }
	r->next = sp->reflist;
	r->filename = filename;
	r->lineno = lineno;
	r->flags = flags;
	sp->reflist = r;
}

/* print the references
   * sort the table alphabetically
   * then flip each entry's reflist to get it into forward order
   * and print it out
   */
/* aux function for sorting */
static int
symcompare(const void *xa, const void *xb)
{
	const struct symbol *a = xa;
	const struct symbol *b = xb;
	if(!a->name) {
		if(!b->name) return 0; /* both empty */
		return 1; /* put empties at the end */
	}
	if(!b->name) return -1;
	return strcmp(a->name, b->name);
}

void
printrefs()
{
	struct symbol *sp;
	qsort(symtab, NHASH, sizeof(struct symbol), symcompare); /* sort the symbol table */
	for(sp = symtab; sp->name && sp < symtab+NHASH; sp++) {
		char *prevfn = NULL; /* last printed filename, to skip dups */
		/* reverse the list of references */
		struct ref *rp = sp->reflist;
		struct ref *rpp = 0; /* previous ref */
		struct ref *rpn; /* next ref */
		do {
			rpn = rp->next;
			rp->next = rpp;
			rpp = rp;
			rp = rpn;
		} while(rp);
		/* now print the word and its references */
		printf("%10s", sp->name);
		for(rp = rpp; rp; rp = rp->next) {
			if(rp->filename == prevfn) {
				printf(" %d", rp->lineno);
			} else {
				printf(" %s:%d", rp->filename, rp->lineno);
				prevfn = rp->filename;
			}
		}
		printf("\n");
	}
}

/*
   The %option line has two options we haven’t seen before, both of which are quite useful.
   The %yylineno option tells flex to define an integer variable called yylineno and to
   maintain the current line number in it. What that means is that every time the scanner
   reads a newline character, it increments yylineno, and if the scanner backs up over a
   newline (using some features we’ll get to later), it decrements it. It’s still up to you to
   initialize yylineno to 1 at the beginning of each file and to save and restore it if you’re
   handling include files. Even with those limitations, it’s still easier than doing line numbers by hand. (In this example, there’s only a single pattern that matches \n, which
   wouldn’t be hard to get right, but it’s quite common to have several patterns that match,
   causing hard-to-track bugs when some but not all of them update the line number.)
   The other new option is case-insensitive, which tells flex to build a scanner that treats
   upper- and lowercase the same. What this means is that a pattern like abc will match
   abc, Abc, ABc, AbC, and so forth. It does not have any effect on your input; in particular,
   the matched string in yytext is not case folded or otherwise modified.
   The symbol table is just an array of symbol structures, each of which contains a pointer
   to the name (i.e., the word in the concordance) and a list of references. The references
   are a linked list of line numbers and pointers to the filename. We also define
   curfilename, a static pointer to the name of the current file, for use when adding references.

   Concordances usually don’t index common short words, so the first set of patterns
   matches and ignores them. An action consisting solely of a vertical bar tells flex that
   the action for this rule is the same as the action for the next rule. The action on the last
   ignored word to does nothing, which is all we need to do to ignore a word.
   The next rule is the meat of the scanner and matches a reasonable approximation of
   an English word. It matches a string of letters, [a-z]+, optionally followed by an apostrophe and either s or t, to match words such as owner’s and can’t. Each matched word is passed to addref(), described in a moment, along with the current filename and line
   number.
   The final pattern is a catchall to match whatever the previous patterns didn’t.
   Note that this scanner is extremely ambiguous, but flex’s rules for resolving ambiguity
   make it do what we want. It prefers longer matches to shorter ones, so the word toad
   will be matched by the main word pattern, not to. If two patterns make an exact match,
   it prefers the earlier one in the program, which is why we put the ignore rules first and
   the catchall last.
   */
