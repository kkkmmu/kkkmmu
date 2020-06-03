%include "functions.asm"

section .data
	filename db "readme.txt", 0x0

section .text
	global _start

_start:
	mov ebx, filename
	mov eax, 10
	int 0x80
	call quit

	
