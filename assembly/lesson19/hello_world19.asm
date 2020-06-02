%include "functions.asm"

;call sys_exec

section .data:
	command1       db "/bin/echo", 0x0
	arg1           db "Hello World!", 0x0
	arguments1     dd command1
	               dd arg1
	               dd 0x0
	environment1   dd 0x0

	command2       db "/bin/ls", 0x0
	arg2	       db "-l", 0x0  
	arguments2     dd command2
	               dd arg2
	               dd 0x0
	environment2   dd 0x0

section .text:
	global _start
	
_start:	
;	mov edx, environment1
;	mov ecx, arguments1
;	mov ebx, command1
;	mov eax, 11
;	int 0x80

	mov edx, environment2
	mov ecx, arguments2
	mov ebx, command2
	mov eax, 11
	int 0x80

	call quit
