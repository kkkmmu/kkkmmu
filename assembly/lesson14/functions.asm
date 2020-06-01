iprint:
	push eax
	push ecx
	push edx
	push esi
	mov ecx, 0

__divide:
	inc ecx                  ;keep tracking of the counter
	mov edx, 0
	mov esi, 10
	idiv esi                 ;divide eax by esi, the quotient part of the value is left in eax and the remainder part is put into edx
	add edx, 48              ;convert the integer to its ascii presentation
	push edx
	cmp eax, 0               ;is all value been processed ?
	jnz __divide

__print:
	dec ecx
	mov eax, esp              ;print the value
	call sprint
	pop eax
	cmp ecx, 0
	jnz __print

	pop esi
	pop edx
	pop ecx
	pop eax
	ret
	

iprintln:
	call iprint

	push eax
	mov eax, 0xa
	push eax
	mov eax, esp
	call sprint
	pop eax
	pop eax
	ret

strlen:
	push ebx
	mov ebx, eax
__next:
	cmp byte [eax], 0
	jz __finished
	inc eax
	jmp __next

__finished:
	sub eax, ebx
	pop ebx
	ret

sprint:
	push edx
	push ecx
	push ebx
	push eax
	call strlen

	mov edx, eax
	pop eax

	mov ecx, eax
	mov ebx, 1
	mov eax, 4
	int 0x80

	pop ebx
	pop ecx
	pop edx
	ret

sprintln:
	call sprint

	push eax
	mov eax, 0xa
	push eax
	mov eax, esp
	call sprint

	pop eax,
	pop eax
	ret

quit:
	mov ebx, 0
	mov eax, 1
	int 0x80
	ret

