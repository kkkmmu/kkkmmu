Calculator (atoi)
	Our program will take several command line arguments and ADD them together printing out the result in the terminal.

	Writing our program
	Our program begins by using the POP instruction to get the number of passed arguments off the stack. This value is stored in ECX (originally known as the counter register). It will then POP the next value off the stack containing the program name and remove it from the number of arguments stored in ECX. It will then loop through the rest of the arguments popping each one off the stack and performing our addition logic. As we know, arguments passed via the command line are received by our program as strings. Before we can ADD the arguments together we will need to convert them to integers otherwise our result will not be correct. We do this by calling our Ascii to Integer function (atoi). This function will convert the ascii value into an integer and place the result in EAX. We can then ADD this value to EDX (originally known as the data register) where we will store the result of our additions. If the value passed to atoi is not an ascii representation of an integer our function will return zero instead. When all arguments have been converted and added together we will print out the result and call our quit function.

	How does the atoi function work
	Converting an ascii string into an integer value is not a trivial task. We know how to convert an integer to an ascii string so the process should essentially work in reverse. Firstly we take the address of the string and move it into ESI (originally known as the source register). We will then move along the string byte by byte (think of each byte as being a single digit or decimal placeholder). For each digit we will check if it's value is between 48-57 (ascii values for the digits 0-9).

	Once we have performed this check and determined that the byte can be converted to an integer we will perform the following logic. We will subtract 48 from the value – converting the ascii value to it's decimal equivalent. We will then ADD this value to EAX (the general purpose register that will store our result). We will then multiple EAX by 10 as each byte represents a decimal placeholder and continue the loop.

	When all bytes have been converted we need to do one last thing before we return the result. The last digit of any number represents a single unit (not a multiple of 10) so we have multiplied our result one too many times. We simple divide it by 10 once to correct this and then return.

	What is the BL register
	You may have noticed that the atoi function references the BL register. So far in these tutorials we have been exclusively using 32bit registers. These 32bit general purpose registers contain segments of memory that can also be referenced. These segments are available in 16bits and 8bits. We wanted a single byte (8bits) because a byte is the size of memory that is required to store a single ascii character. If we used a larger memory size we would have copied 8bits of data into 32bits of space leaving us with 'rubbish' bits - because only the first 8bits would be meaningful for our calculation.

	The EBX register is 32bits. EBX'S 16 bit segment is referenced as BX. BX contains the 8bit segments BL and BH (Lower and Higher bits). We wanted the first 8bits (lower bits) of EBX and so we referenced that storage area using BL.

	Almost every assembly language tutorial begins with a history of the registers, their names and their sizes. These tutorials however were written to provide a foundation in NASM by first writing code and then understanding the theory. The full story about the size of registers, their history and importance are beyond the scope of this tutorial but we will return to that story in later tutorials.

	Note: Only the new function in this file 'atoi' is shown below.
