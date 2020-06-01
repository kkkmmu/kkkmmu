%include "functions.asm"

section .text
	global _start

_start:
	mov eax, 90
	mov ebx, 9
	sub eax, ebx     ;sub eax with ebx and save the result in eax
	call iprintln

	call quit
