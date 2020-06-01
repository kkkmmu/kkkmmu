%include "functions.asm"

section .text:
	global _start

_start:
	mov ecx, 0
__inext:
	inc ecx
	mov eax, ecx
	call iprintln
	cmp ecx, 10
	jne __inext

	call quit
