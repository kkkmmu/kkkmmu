%include 'functions.asm'

section .data:
	message1 db "Hello world from China!", 0xA
	message2 db "Hello world from ShannXi!", 0xA

section .text:
global _start

_start:
	mov eax, message1
	call sprint

	mov eax, message2
	call sprint

	call quit
