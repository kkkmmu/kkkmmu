%include "functions.asm"

section .text:
	global _start

_start:
	xor eax, eax,
	xor ebx, ebx,
	xor edi, edi,
	xor esi, esi

_socket:
	push byte 6         ;IPPROTO_TCP
	push byte 1         ;SOCK_STREAM
	push byte 2         ;PF_INET
	mov ecx, esp
	mov ebx, 1          ;socket
	mov eax, 102        ;sys_socketcall
	int 0x80

	call iprintln

	call quit

