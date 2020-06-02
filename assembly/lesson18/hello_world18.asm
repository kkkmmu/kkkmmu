%include "functions.asm"

section .data:
	fizz db "Fizz", 0x0
	buzz db "Buzz", 0x0

section .text:
	global _start:

_start:
	mov esi, 0
	mov edi, 0
	mov ecx, 0

.next:
	inc ecx

.check_fizz:
	mov edx, 0
	mov eax, ecx
	mov ebx, 3
	div ebx                 ;divide eax by ebx, remainder is saved in edx
	mov edi, edx
	cmp edi, 0              
	jne .check_buzz         ;if this not a fizz
	mov eax, fizz           ;if this is a fizz, print it out
	call sprint

.check_buzz:
	mov edx, 0              ;reset edx
	mov eax, ecx
	mov ebx, 5
	div ebx                 ;divide eax by ebx, remainder is saved in edx
	mov esi, edx
	cmp esi, 0
	jne .check_int          ;if this is not a buzz
	mov eax, buzz           ;if this is a buzz, print it out
	call sprint

.check_int:
	cmp edi, 0  
	je .continue
	cmp esi, 0
	je .continue
	mov eax, ecx            ;if this not a fizz and not a buzz, just print it out
	call iprint

.continue:
	mov eax, 0xa
	push eax
	mov eax, esp            ;print a string end '0'
	call sprint
	pop eax
	cmp ecx, 100            ;only print smaller than 100
	jne .next

	call quit
