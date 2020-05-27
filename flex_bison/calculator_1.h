/*
 * * Declarations for a calculator fb3-1
 * */
/* interface to the lexer */
extern int yylineno; /* from lexer */
void yyerror(char *s, ...);
/* nodes in the abstract syntax tree */
struct ast {
	int nodetype;
	struct ast *l;
	struct ast *r;
};
struct numval {
	int nodetype; /* type K for constant */
	double number;
};
/* build an AST */
struct ast *newast(int nodetype, struct ast *l, struct ast *r);
struct ast *newnum(double d);
/* evaluate an AST */
double eval(struct ast *);
/* delete and free an AST */
void treefree(struct ast *);

/*
 * The variable yylineno and routine yyerror are familiar from the flex example. Our
 * yyerror is slightly enhanced to take multiple arguments in the style of printf.
 * The AST consists of nodes, each of which has a node type. Different nodes have different
 * fields, but for now we have just two kinds, one that has pointers to up to two subnodes
 * and one that contains a number. Two routines, newast and newnum, create AST nodes;
 * eval walks an AST and returns the value of the expression it represents; and treefree
 * walks an AST and deletes all of its nodes.
 */
