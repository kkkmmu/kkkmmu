File Handling - Read
Building upon the previous lesson we will now use sys_read to read the content of a newly created and opened file. We will store this string in a variable.

sys_read expects 3 arguments - the number of bytes to read in EDX, the memory address of our variable in ECX and the file descriptor in EBX. We will use the previous lessons sys_open code to obtain the file descriptor which we will then load into EBX. The sys_read opcode is then loaded into EAX and the kernel is called to read the file contents into our variable and is then printed to the screen.

Note: We will reserve 255 bytes in the .bss section to store the contents of the file. See Lesson 9 for more information on the .bss section.
