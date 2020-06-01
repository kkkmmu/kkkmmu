%include "functions.asm"

section .text
global _start

_start:
	mov eax, 90
	mov ebx, 9
	add eax, ebx   ; add eax and ebx, save the value in eax
	call iprintln

	call quit
