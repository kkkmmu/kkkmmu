section .data:
	message db "hello brave new world!", 0xA

section .text:
global _start
	
_start:
	mov ebx, message
	mov eax, ebx

next_char:
	cmp byte [eax], 0   ;compare the byte pointed to by eax against zero.
	jz finished         ;jump (if the zero flags has been set) to finished.
	inc eax             ;increase eax
	jmp next_char

finished:
	sub eax, ebx        ;substract the address in ebx from the address in eax. (Get the string length)
	mov edx, eax,
	mov ecx, message
	mov ebx, 1
	mov eax, 4
	int 0x80

	mov ebx, 0
	mov eax, 1
	int 0x80
