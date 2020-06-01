%include "functions.asm"

section .data
	message1 db " remainder "

section .text
	global _start

_start:
	mov eax, 91
	mov ebx, 9
	div ebx                  ;divide eax by eax, keep the quitient part in eax, and keep the remainder in edx
	call iprint
	mov eax, message1
	call sprint
	mov eax, edx
	call iprintln
	call quit
