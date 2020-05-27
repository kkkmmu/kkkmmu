In order to write a parser, we need some way to describe the rules the parser uses to turn a sequence of tokens into a parse tree. The most common kind of language that
computer parsers handle is a context-free grammar (CFG).
The standard form to write down a CFG is Backus-Naur Form (BNF), created around 1960 to describe Algol 60 and named after two members of the Algol 60 committee.

CFGs are also known as phrase-structure grammars or type-3 languages. Computer theorists and natural language linguists independently developed them at about the same time in the late 1950s. If you’re a computer scientist, you usually call them CFGs, and if you’re a linguist, you usually call them PSGs or type-3, but they’re the same thing.

Fortunately, BNF is quite simple. Here’s BNF for simple arithmetic expressions enough to handle 1 * 2 + 3 * 4 + 5:
		<exp> ::= <factor>
			  | <exp> + <factor>
		<factor> ::= NUMBER
			  | <factor> * NUMBER

	Each line is a rule that says how to create a branch of the parse tree. 
	In BNF, ::= can be read “is a” or “becomes,” and | is “or,” another way to create a branch of the same kind. 
	The name on the left side of a rule is a symbol or term. By convention, all tokens are considered to be symbols, but there are also symbols that are not tokens.

Bison rules are basically BNF, with the punctuation simplified a little to make them easier to type. 

Bison takes a grammar that you specify and writes a parser that recognizes valid “sentences” in that grammar. We use the term sentence here in a fairly general way—for a C language grammar, the sentences are syntactically valid C programs. Programs can be syntactically valid but semantically invalid, for example, a C program that assigns a string to an int variable. Bison handles only the syntax; other validation is up to you.

A grammar is a series of rules that the parser uses to recognize syntactically valid input. 
	statement: NAME '=' expression
	expression: NUMBER '+' NUMBER
				| NUMBER '−' NUMBER
	The vertical bar, |, means there are two possibilities for the same symbol; that is, an
	expression can be either an addition or a subtraction. The symbol to the left of the : is
	known as the left-hand side of the rule, often abbreviated LHS, and the symbols to the
	right are the right-hand side, usually abbreviated RHS. Several rules may have the same
	left-hand side; the vertical bar is just shorthand for this. Symbols that actually appear
	in the input and are returned by the lexer are terminal symbols or tokens, while those
	that appear on the left-hand side of each rule are nonterminal symbols or nonterminals.
	Terminal and nonterminal symbols must be different; it is an error to write a rule with
	a token on the left side.

Every grammar includes a start symbol, the one that has to be at the root of the parse tree. In this grammar, statement is the start symbol.

Rules can refer directly or indirectly to themselves; this important ability makes it possible to parse arbitrarily long input sequences. Let’s extend our grammar to handle longer arithmetic expressions:
	expression: NUMBER
		| expression + NUMBER
		| expression − NUMBER
	Now we can parse a sequence like fred = 14 + 23 − 11 + 7 by applying the expression rules repeatedly. Bison can parse recursive rules very efficiently, so we will see recursive rules in nearly every grammar we use.

A bison parser works by looking for rules that might match the tokens seen so far. When bison processes a parser, it creates a set of states, each of which reflects a possible position in one or more partially parsed rules. As the parser reads tokens, each time it reads a token that doesn’t complete a rule, it pushes the token on an internal stack and switches to a new state reflecting the token it just read. This action is called a shift. When it has found all the symbols that constitute the right-hand side of a rule, it pops the right-hand side symbols off the stack, pushes the left-hand side symbol onto the stack, and switches to a new state reflecting the new symbol on the stack. This action is called a reduction, since it usually reduces the number of items on the stack.* Whenever bison reduces a rule, it executes user code associated with the rule. This is how you actually do something with the material that the parser parses. 

A Bison Parser
	A bison specification has the same three-part structure as a flex specification. (Flex copied its structure from the earlier lex, which copied its structure from yacc, the predecessor of bison.) 
	The first section, the definition section, handles control information for the parser and generally sets up the execution environment in which the parser will operate. 
	The second section contains the rules for the parser. 
	The third section is C code copied verbatim into the generated C program.
	Bison creates the C program by plugging pieces into a standard skeleton file. The rules are compiled into arrays that represent the state machine that matches the input tokens. The actions have the $N and @N values translated into C and then are put into a switch statement within yyparse() that runs the appropriate action each time there’s a reduction. Some bits of the skeleton have multiple versions from which bison chooses depending on what options are in use; for example, if the parser uses the locations feature, it includes code to handle location data.

An AST is basically a parse tree that omits the nodes for the uninteresting rules. Once a parser creates an AST, it’s straightforward to write recursive routines that “walk” the tree.

There are two ways to specify precedence and associativity in a grammar, implicitly and explicitly. So far, we’ve specified them implicitly, by using separate nonterminal symbols for each precedence level. This is a perfectly reasonable way to write a grammar, and if bison didn’t have explicit precedence rules, it would be the only way.
But bison also lets you specify precedence explicitly. 
