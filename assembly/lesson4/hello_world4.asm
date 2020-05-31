section .data:
	message db "Hello world from China", 0xA

section .text:
global _start

_start:
	mov eax, message
	call strlen           ;call function strlen
	mov edx, eax
	mov ecx, message
	mov ebx, 1
	mov eax, 4
	int 0x80
	
	mov ebx, 0
	mov eax, 1
	int 0x80
	;_start function ends here

	;strlen function starts here
strlen:
	push ebx              ;save ebx onto the stack
	mov ebx, eax
	
next_char:
	cmp byte [eax], 0
	jz finished
	inc eax
	jmp next_char

finished:
	sub eax, ebx
	pop ebx               ;retore ebx from stack
	ret
