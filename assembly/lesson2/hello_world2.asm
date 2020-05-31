section .data:
	message db "Hello world!", 0xA

section .text:
global _start

_start:
	mov edx, 13
	mov ecx, message
	mov ebx, 1
	mov eax, 4
	int 0x80

	mov ebx, 0
	mov eax, 1
	int 0x80


