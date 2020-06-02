%include "functions.asm"

section .data:
	child      db   "This is the child process", 0x0
	parent     db   "This is the parent process", 0x0

section .text:
	global _start

_start:
	mov eax, 2      ;call sys_fork
	int 0x80

	cmp eax, 0
	jz .child

.parent:
	mov eax, parent
	call sprintln
	call quit

.child:
	mov eax, child
	call sprintln
	call quit
