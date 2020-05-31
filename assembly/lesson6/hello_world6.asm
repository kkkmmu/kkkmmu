%include "functions.asm"

section .data:
	message1 db "Hello world from China!", 0xA, 0x0
	message2 db "Hello world from ShannXi!", 0xA, 0x0

section .text:

global _start

_start:
	mov eax, message1 
	call sprint
	mov eax, message2
	call sprint
	call quit
