%include "functions.asm"

section .data
	message db  "Seconds since Jan 01 1970: ", 0x0

section .text
	global _start

_start:
	mov eax, message
	call sprint

	mov eax, 13
	int 0x80

	call iprintln
	call quit
