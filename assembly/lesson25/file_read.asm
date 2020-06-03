%include "functions.asm"

section .data
	filename db "readme.txt", 0x0
	content  db "Hi, How are you!", 0x0

section .bss
	file_buffer resb 255,

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
	mov eax, 4          ;sys_write
	int 0x80

	mov ecx, 0
	mov ebx, filename
	mov eax, 5          ;sys_open
	int 0x80
	
	mov edx, 16
	mov ecx, file_buffer
	mov ebx, eax,
	mov eax, 3
	int 0x80

	mov eax, file_buffer
	call sprintln 
	call quit
