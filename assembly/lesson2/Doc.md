Some more background
After successfully learning how to execute a system call in Lesson 1 we now need to learn about one of the most important system calls in the kernel, sys_exit.

Notice how after our 'Hello, world!' program ran we got a Segmentation fault? Well, computer programs can be thought of as a long strip of instructions that are loaded into memory and divided up into sections (or segments). This general pool of memory is shared between all programs and can be used to store variables, instructions, other programs or anything really. Each segment is given an address so that information stored in that section can be found later.

To execute a program that is loaded in memory, we use the global label _start: to tell the operating system where in memory our program can be found and executed. Memory is then accessed sequentially following the program logic which determines the next address to be accessed. The kernel jumps to that address in memory and executes it.

It's important to tell the operating system exactly where it should begin execution and where it should stop. In Lesson 1 we didn't tell the kernel where to stop execution. So, after we called sys_write the program continued sequentially executing the next address in memory, which could have been anything. We don't know what the kernel tried to execute but it caused it to choke and terminate the process for us instead - leaving us the error message of 'Segmentation fault'. Calling sys_exit at the end of all our programs will mean the kernel knows exactly when to terminate the process and return memory back to the general pool thus avoiding an error.

Writing our program
sys_exit has a simple function definition. In the Linux System Call Table it is allocated OPCODE 1 and is passed a single argument through EBX.

In order to execute this function all we need to do is:

Load EBX with 0 to pass zero to the function meaning 'zero errors'.
Load EAX with 1 to call sys_exit.
Then request an interrupt on libc using INT 80h.
We then compile, link and run it again.

Note: Only new code added in each lesson will be commented.
