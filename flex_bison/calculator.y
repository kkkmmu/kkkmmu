/* simplest version of calculator */
%{
#include <stdio.h>
%}
/* declare tokens */
%token NUMBER
%token ADD SUB MUL DIV ABS
%token EOL
%token OP CP
%token AND OR
%%
calclist: /* nothing  matches at beginning of input */
		| calclist exp EOL { printf("= %d\n", $2); } /*EOL is end of an expression; */
;
exp: cexp
   | exp AND cexp { $$ = $1 && $3; }
   | exp OR cexp { $$ = $1 || $3; }

cexp: factor /*default $$ = $1*/
   | cexp ADD factor { $$ = $1 + $3; printf(" %d %d %d\n", $1, $3, $$);}
	| cexp SUB factor { $$ = $1 - $3; printf(" %d %d %d\n", $1, $3, $$);}
;
factor: term /*default $$ = $1 */
	  | factor MUL term { $$ = $1 * $3; printf(" %d %d %d\n", $1, $3, $$);}
	  | factor DIV term { $$ = $1 / $3; printf(" %d %d %d\n", $1, $3, $$);}
;
term: NUMBER /*default $$ = $1 */
	| ABS term { $$ = $2 >= 0? $2 : - $2; printf(" %d %d\n", $2, $$);}
	| OP exp CP {$$ = $2;}
;
%%
main(int argc, char **argv)
{
yyparse();
}
yyerror(char *s)
{
fprintf(stderr, "error: %s\n", s);
}

/*
Precedence and Ambiguity. The separate symbols for term, factor, and exp tell bison to handle ABS, then MUL and DIV, and then ADD and SUB. In general, whenever a grammar has multiple levels of precedence where one kind of operator binds “tighter” than another, the parser will need a level of rule for each level.
*/

/*
	Bison programs have (not by coincidence) the same three-part structure as flex programs, with declarations, rules, and C code. 
		The declarations here include C code to be copied to the beginning of the generated C parser, again enclosed in %{ and %}.
 		Following that are %token token declarations, telling bison the names of the symbols in
			the parser that are tokens. By convention, tokens have uppercase names, although bison
			doesn’t require it. 
		Any symbols not declared as tokens have to appear on the left side of at least one rule in the program. 
			(If a symbol neither is a token nor appears on the left side of a rule, it’s like an unreferenced variable in a C program. 
			It doesn’t hurt anything, but it probably means the programmer made a mistake.)

		The second section contains the rules in simplified BNF. Bison uses a single colon rather than ::=, 
			and since line boundaries are not significant, a semicolon marks the end of a
	 		rule. Again, like flex, the C action code goes in braces at the end of each rule.
		Bison automatically does the parsing for you, remembering what rules have been
			matched, so the action code maintains the values associated with each symbol. Bison
			parsers also perform side effects such as creating data structures for later use or, as in
			this case, printing out results. 
		The symbol on the left side of the first rule is the start
			symbol, the one that the entire input has to match. There can be, and usually are, other
			rules with the same start symbol on the left.
		Each symbol in a bison rule has a value; the value of the target symbol (the one to the
			left of the colon) is called $$ in the action code, and the values on the right are numbered
			$1, $2, and so forth, up to the number of symbols in the rule. The values of tokens are
			whatever was in yylval when the scanner returned the token; the values of other symbols are set in rules in the parser. 
			In this parser, the values of the factor, term, and exp symbols are the value of the expression they represent.
		In this parser, the first two rules, which define the symbol calcset, implement a loop
			that reads an expression terminated by a newline and prints its value. The definition
			of calclist uses a common two-rule recursive idiom to implement a sequence or list:
			the first rule is empty and matches nothing; the second adds an item to the list. The
			action in the second rule prints the value of the exp in $2.
			The rest of the rules implement the calculator. The rules with operators such as exp
			ADD factor and ABS term do the appropriate arithmetic on the symbol values. The rules
			with a single symbol on the right side are syntactic glue to put the grammar together;
			for example, an exp is a factor. In the absence of an explicit action on a rule, the parser
			assigns $1 to $$. This is a hack, albeit a very useful one, since most of the time it does
			the right thing.
*/
