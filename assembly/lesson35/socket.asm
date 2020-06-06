%include "functions.asm"

section .data
	response db 'HTTP/1.1 200 OK', 0Dh, 0Ah, 'Content-Type: text/html', 0Dh, 0Ah, 'Content-Length: 14', 0Dh, 0Ah, 0Dh, 0Ah, 'Hello World!', 0Dh, 0Ah, 0h

section .bss
	buffer resb 1024

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

_accept:
	push byte 0           ;address length
	push byte 0           ;address argument
	push edi              ;listen file descriptor
	mov ecx, esp
	mov ebx, 5            ;accept
	mov eax, 102
	int 0x80

_fork:
	mov esi, eax
	mov eax, 2
	int 0x80
	jz _read
	jmp _accept

_read:
	mov edx, 1024
	mov ecx, buffer
	mov ebx, esi
	mov eax, 3
	int 0x80
	mov eax, buffer

_write:
	mov edx, 78
	mov ecx, response
	mov ebx, esi
	mov eax, 4
	int 0x80
_close:
	mov ebx, esi
	mov eax, 6
	int 0x80

	call quit

