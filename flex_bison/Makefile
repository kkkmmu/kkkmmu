calculator: calculator.y calculator.l
	bison -d calculator.y 
	flex calculator.l 
	gcc -o $@ calculator.tab.c lex.yy.c -lfl

calculator_1: calculator_1.y calculator_1.l calculator_1.h calculator_1.c
	bison -d calculator_1.y 
	flex calculator_1.l 
	gcc -o $@ calculator_1.tab.c lex.yy.c calculator_1.c -lfl

compiler: compiler.y compiler.l compiler.h compiler.c
	bison -d compiler.y 
	flex compiler.l 
	gcc -o $@ compiler.tab.c lex.yy.c compiler.c -lfl -lm

targets1 := c99
$(targets1):
	bison -d $@.y 
	flex $@.l 
	gcc -o $@ $@.tab.c lex.yy.c -lfl

clean:
	rm -rf *.tab.h
	rm -rf *.tab.c
	rm -rf lex.yy.c
	find . -xtype f -executable | xargs rm -rf

targets := file mfile include cancordance ccr 
$(targets):
	flex $@.l
	gcc -o $@ lex.yy.c -lfl
