%include "functions.asm"

section .data
	filename db "readme.txt", 0x0
	content  db "Hi, How are you!", 0x0

section .text
	global _start

_start:
	mov ecx, 0x777
	mov ebx, filename
	mov eax, 8          ;create the file 
	int 0x80

	mov edx, 16
	mov ecx, content
	mov ebx, eax
	mov eax, 4
	int 0x80
	call quit
