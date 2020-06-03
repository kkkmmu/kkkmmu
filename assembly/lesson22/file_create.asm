%include "functions.asm"

section .data
	filename db "readme.txt"

section .text
	global _start

_start:
	mov ecx, 0x777
	mov ebx, filename
	mov eax, 8

	int 0x80

	call iprintln

	call quit
