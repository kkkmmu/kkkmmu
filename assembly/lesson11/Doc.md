Count to 10 (itoa)
So why did our program in Lesson 10 print out a colon character instead of the number 10?. Well lets have a look at our ASCII table. We can see that the colon character has a ASCII value of 58. We were adding 48 to our integers to convert them to their ASCII string representations so instead of passing sys_write the value '58' to print ten we actually need to pass the ASCII value for the number 1 followed by the ASCII value for the number 0. Passing sys_write '4948' is the correct string representation for the number '10'. So we can't just simply ADD 48 to our numbers to convert them, we first have to divide them by 10 because each place value needs to be converted individually.

We will write 2 new subroutines in this lesson 'iprint' and 'iprintLF'. These functions will be used when we want to print ASCII string representations of numbers. We achieve this by passing the number in EAX. We then initialise a counter in ECX. We will repeatedly divide the number by 10 and each time convert the remainder to a string by adding 48. We will then PUSH this onto the stack for later use. Once we can no longer divide the number by 10 we will enter our second loop. In this print loop we will print the now converted string representations from the stack and POP them off. Popping them off the stack moves ESP forward to the next item on the stack. Each time we print a value we will decrease our counter ECX. Once all numbers have been converted and printed we will return to our program.

How does the divide instruction work?
The DIV and IDIV instructions work by dividing whatever is in EAX by the value passed to the instruction. The quotient part of the value is left in EAX and the remainder part is put into EDX (Originally called the data register).

	For example.

	IDIV instruction example
	mov     eax, 10         ; move 10 into eax
	mov     esi, 10         ; move 10 into esi
	idiv    esi             ; divide eax by esi (eax will equal 1 and edx will equal 0)
	idiv    esi             ; divide eax by esi again (eax will equal 0 and edx will equal 1)
If we are only storing the remainder won't we have problems?
No, because these are integers, when you divide a number by an even bigger number the quotient in EAX is 0 and the remainder is the number itself. This is because the number divides zero times leaving the original value as the remainder in EDX. How good is that?

Note: Only the new functions iprint and iprintLF have comments.
