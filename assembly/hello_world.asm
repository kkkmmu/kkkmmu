global _start;

section .text

_start:
	mov eax, 0x4               ;system call number should be passed with eax
	mov ebx, 1                 ;first parameter for write. (fd)
	mov ecx, message           ;second parameter for write. (buf)
	mov edx, message_length    ;thired parameter for write. (buf_len)
	int 0x80                   ;do system call

	mov eax, 0x1               ;another system call (exit)
	mov ebx, 0                 ;parameter
	int 0x80                   ;do system call

section .data:
	message db "Hello World!", 0xA
	message_length equ $-message


;nasm -f elf64 -o hello_world.o hello_world.asm
;ld -m elf_x86_64 -o hello_world hello_world.o
