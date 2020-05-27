/*
 flex is not a bad tool to use for doing modest
 text transformations and for programs that collect statistics on input. More often than
 not, though, you’ll want to use flex to generate a scanner that divides the input into
 tokens that are then used by other parts of your program.
 */

%%
"colour" { printf("color"); }
"flavour" { printf("flavor"); }
"clever" { printf("smart"); }
"smart" { printf("elegant"); }
"conservative" { printf("liberal"); }
. { printf("%s", yytext); }
%%

/* 
   flex test.l
   gcc lex.yy.c -lfl
   ./a.out
 */
