%include "functions.asm"

section .data
	message1 db "Please enter your name: ", 0x0
	message2 db "Hello, ", 0x0
	size dw 1024

section .bss
	buf resb 1024

section .text
	global _start

_start
	mov eax, message1
	call sprint

	mov edx, size              ;buffer size to read
	mov ecx, buf               ;buffer for read
	mov ebx, 0                 ;read from standard input
	mov eax, 3                 ;read system call
	int 0x80

	mov eax, message2
	call sprint

	mov eax, buf
	call sprint

	call quit
