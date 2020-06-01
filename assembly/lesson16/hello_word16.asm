%include "functions.asm"

section .text
	global _start

_start:
	pop ecx         ;save argc into ecx
	pop edx         ;save program name into edx
	sub ecx, 1      ;decrease argc by 1 for program name
	mov edx, 0

__inext:
	cmp ecx, 0x0    ;if there is no user argument
	jz __no_arg
	pop eax         ;get next argv
	call atoi
	add edx, eax,
	dec ecx
	jmp __inext

__no_arg:
	mov  eax, edx,
	call iprintln
	call quit
	
