%include "functions.asm"

section .data:
	message1 db "Jumping to finished label", 0x0 
	message2 db "Inside subroutine number: ", 0x0
	message3 db 'Inside sburoutine "finished".', 0x0

section .text:
	global _start

_start:

subroutine1:	
	mov eax, message1
	call sprintln
	jmp .finished

.finished:
	mov eax, message2
	call sprint
	mov eax, 1
	call iprintln

subroutine2:
	mov eax, message1
	call sprintln
	jmp .finished

.finished:
	mov eax, message2
	call sprint
	mov eax, 2
	call iprintln

	mov eax, message1
	call sprintln
	jmp finished

finished:
	mov eax, message3
	call sprintln
	call quit
