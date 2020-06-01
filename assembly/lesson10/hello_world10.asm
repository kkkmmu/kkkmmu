;Counting from 0 to 10

%include "functions.asm"

section .text
	global _start


_start:
	mov ecx, 0
next:
	inc ecx

	mov eax, ecx,
	add eax, 48            ;convert interger value to its ascii presentation
	push eax               ;in order to call our print function, a pointer is necessary
	mov eax, esp           
	call sprintln

	pop eax
	cmp ecx, 10
	jne next

	call quit
