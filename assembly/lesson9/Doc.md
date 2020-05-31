User input
Introduction to the .bss section
So far we've used the .text and .data section so now it's time to introduce the .bss section. BSS stands for BLOCK Started by Symbol. It is an area in our program that is used to reserve space in memory for uninitialised variables. We will use it to reserve some space in memory to hold our user input since we don't know how many bytes we'll need to store.

The syntax to declare variables is as follows:

.bss section example
1 SECTION .bss
2 variableName1:      RESB    1       ; reserve space for 1 byte
3 variableName2:      RESW    1       ; reserve space for 1 word
4 variableName3:      RESD    1       ; reserve space for 1 double word
5 variableName4:      RESQ    1       ; reserve space for 1 double precision float (quad word)
6 variableName5:      REST    1       ; reserve space for 1 extended precision float

Writing our program
We will be using the system call sys_read to receive and process input from the user. This function is assigned OPCODE 3 in the Linux System Call Table. Just like sys_write this function also takes 3 arguments which will be loaded into EDX, ECX and EBX before requesting a software interrupt that will call the function.

The arguments passed are as follows:

EDX will be loaded with the maximum length (in bytes) of the space in memory.
ECX will be loaded with the address of our variable created in the .bss section.
EBX will be loaded with the file we want to write to – in this case STDIN.
As always the datatype and meaning of the arguments passed can be found in the function's definition.
