Nssing arguments
Passing arguments to your program from the command line is as easy as popping them off the stack in NASM. When we run our program, any passed arguments are loaded onto the stack in reverse order. The name of the program is then loaded onto the stack and lastly the total number of arguments is loaded onto the stack. The last two stack items for a NASM compiled program are always the name of the program and the number of passed arguments.

So all we have to do to use them is POP the number of arguments off the stack first, then iterate once for each argument and perform our logic. In our program that means calling our print function.

Note: We are using the ECX register as our counter for the loop. Although it's a general-purpose register it's original intention was to be used as a counter.ote: I've highlighted the new code in functions.asm below.

