%include "functions.asm"

section .text:
global _start

_start:
	pop ecx                   ;first value on the stack is the number of arguments

nextArg:
	cmp ecx, 0x0              ;check to see if we have any arguments left
	jz noMoreArgs             ;if zero flag is set jump to noMoreArgs label
	pop eax                   ;pop the next arguments of the stack
	call sprintln
	dec ecx                   ;decrease ecx (number of arguments left) by 1
	jmp nextArg

noMoreArgs:
	call quit
