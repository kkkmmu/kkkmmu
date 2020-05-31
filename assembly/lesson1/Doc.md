https://asmtutor.com/#lesson1

Lesson 1
Hello, world!
First, some background
Assembly language is bare-bones. The only interface a programmer has above the actual hardware is the kernel itself. In order to build useful programs in assembly we need to use the linux system calls provided by the kernel. These system calls are a library built into the operating system to provide functions such as reading input from a keyboard and writing output to the screen.

When you invoke a system call the kernel will immediately suspend execution of your program. It will then contact the necessary drivers needed to perform the task you requested on the hardware and then return control back to your program.

Note: Drivers are called drivers because the kernel literally uses them to drive the hardware.

We can accomplish this all in assembly by loading EAX with the function number (operation code OPCODE) we want to execute and filling the remaining registers with the arguments we want to pass to the system call. A software interrupt is requested with the INT instruction and the kernel takes over and calls the function from the library with our arguments. Simple.

For example requesting an interrupt when EAX=1 will call sys_exit and requesting an interrupt when EAX=4 will call sys_write instead. EBX, ECX & EDX will be passed as arguments if the function requires them. Click here to view an example of a Linux System Call Table and its corresponding OPCODES.

Writing our program
Firstly we create a variable 'msg' in our .data section and assign it the string we want to output in this case 'Hello, world!'. In our .text section we tell the kernel where to begin execution by providing it with a global label _start: to denote the programs entry point.

We will be using the system call sys_write to output our message to the console window. This function is assigned OPCODE 4 in the Linux System Call Table. The function also takes 3 arguments which are sequentially loaded into EDX, ECX and EBX before requesting a software interrupt which will perform the task.

The arguments passed are as follows:

EDX will be loaded with the length (in bytes) of the string.
ECX will be loaded with the address of our variable created in the .data section.
EBX will be loaded with the file we want to write to – in this case STDOUT.
The datatype and meaning of the arguments passed can be found in the function's definition.
We compile, link and run the program using the commands below.
