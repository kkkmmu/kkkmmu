%include "functions.asm"

section .data:
	message1 db "Hello world from China!", 0x0   ;new line char '\n' has been removed in this version
	message2 db "Hello world from ShannXi!", 0x0 ;new line char '\n' has been removed in this version.

section .text:

global _start

_start:
	mov eax, message1 
	call sprintln
	mov eax, message2
	call sprintln
	call quit
