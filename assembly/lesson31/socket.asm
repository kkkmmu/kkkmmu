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

_bind:
	mov edi, eax
	push dword 0x00000000 ;ip 0.0.0.0
	push word 0x2923      ;port 9001
	push word 2           ;AF_INET
	mov ecx, esp
	push byte 16          ;argument length
	push ecx
	push edi
	mov ecx, esp
	mov ebx, 2            ;bind
	mov eax, 102
	int 0x80

_listen:
	push byte 1           ;queue length
	push edi
	mov ecx, esp
	mov ebx, 4            ;listen
	mov eax, 102
	int 0x80

	call iprintln

	call quit

