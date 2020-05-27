/* calculator with AST */
%{
# include <stdio.h>
# include <stdlib.h>
# include "calculator_1.h"
%}
%union {
struct ast *a;
double d;
}
/*
 The first section of the
parser uses the %union construct to declare types to be used in the values of symbols in
the parser. In a bison parser, every symbol, both tokens and nonterminals, can have a
value associated with it. By default, the values are all integers, but useful programs
generally need more sophisticated values. The %union construct, as its name suggests,
is used to create a C language union declaration for symbol values. In this case, the
union has two members; a, which is a pointer to an AST, and d, which is a double
precision number
*/
/* declare tokens */
%token <d> NUMBER
%token EOL
%type <a> exp factor term

/*
Once the union is defined, we need to tell bison what symbols have what types of values
by putting the appropriate name from the union in angle brackets (< >). The token
NUMBER, which represents numbers in the input, has the value <d> to hold the value of
the number. The new declaration %type assigns the value <a> to exp, factor, and term,
which we’ll use as we build up our AST.

You don’t have to declare a type for a token or declare a nonterminal at all if you don’t
use the symbol’s value. If there is a %union in the declarations, bison will give you an
error if you attempt to use the value of a symbol that doesn’t have an assigned type.
Keep in mind that any rule without explicit action code gets the default action $$ =
$1;, and bison will complain if the LHS symbol has a type and the RHS symbol doesn’t
have the same type.
*/

%%
calclist: /* nothing */
		| calclist exp EOL 			{
										printf("= %4.4g\n", eval($2));
										treefree($2);
										printf("> ");
									}
		| calclist EOL 				{ printf("> "); } /* blank line or a comment */
;
exp: factor
   		| exp '+' factor 			{ $$ = newast('+', $1,$3); }
		| exp '-' factor 			{ $$ = newast('-', $1,$3);}
;
factor: term
	  	| factor '*' term 			{ $$ = newast('*', $1,$3); }
		| factor '/' term 			{ $$ = newast('/', $1,$3); }
;
term: NUMBER 						{ $$ = newnum($1); }
		| '|' term 					{ $$ = newast('|', $2, NULL); }
		| '(' exp ')' { $$ = $2; }
		| '-' term { $$ = newast('M', $2, NULL); }
;
%%
