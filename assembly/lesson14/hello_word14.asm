%include "functions.asm"

section .text
	global _start

_start:
	mov eax, 90
	mov ebx, 9
	mul ebx             ;multiply eax by ebx and keep the result in eax
	
	call iprintln
	call quit
