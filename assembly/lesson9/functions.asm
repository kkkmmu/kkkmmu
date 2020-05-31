;strlen function
strlen:
	push ebx
	mov ebx, eax

next_char:
	cmp byte [eax], 0
	jz finished
	inc eax
	jmp next_char

finished:
	sub eax, ebx,              ;sub eax from ebx and save the reslult in eax
	pop ebx
	ret

;sprint function
sprint:
	push edx
	push ecx
	push ebx
	push eax
	call strlen

	mov edx, eax               ;save strlen in edx
	pop eax                    ;restore eax to string pointer

	mov ecx, eax               ;paramter 2 (buffer) to write
	mov ebx, 1
	mov eax, 4
	int 0x80

	pop ebx
	pop ecx
	pop edx
	ret

;sprintln function
sprintln:
	call sprint                ;message have be printed at this point
	push eax                   ;the flowing source code is just use to print '\n'
	mov eax, 0xA               ;the following three step convet char to string pointer.
	push eax
	mov eax, esp
	call sprint
	pop eax
	pop eax
	ret

;exit function
quit:
	mov ebx, 0
	mov eax, 1
	int 0x80
	ret
