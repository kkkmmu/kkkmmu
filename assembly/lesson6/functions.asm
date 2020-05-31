;function to calculate string len
strlen:
	push ebx
	mov ebx, eax

next_char:
	cmp byte [eax], 0
	jz finished
	inc eax
	jmp next_char

finished:
	sub eax, ebx
	pop ebx
	ret

;function to print string
sprint:
	push edx
	push ecx
	push ebx
	push eax
	call strlen

	mov edx, eax,
	pop eax,
	mov ecx, eax,
	mov ebx, 1
	mov eax, 4
	int 0x80

	pop ebx,
	pop ecx,
	pop edx,
	ret

;quit function
quit:
	mov ebx, 0
	mov eax, 1
	int 0x80
	ret
