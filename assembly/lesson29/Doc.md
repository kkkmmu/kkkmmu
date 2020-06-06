Sockets - Create
Firstly, some background
Socket Programming in Linux is achieved through the use of the sys_socketcall kernel function. The sys_socketcall function is somewhat unique in that it encapsulates a number of different subroutines, all related to socket operations, within the one function. By passing different integer values in EBX we can change the behaviour of this function to create, listen, send, receive, close and more. Click here to view the full commented source code of the completed program.

Writing our program
We begin the tutorial by first initalizing some of our registers which we will use later to store important values. We will then create a socket using sys_socketcall's first subroutine which is called 'socket'. We will then build upon our program in each of the following socket programming lessons, adding code as we go. Eventually we will have a full program that can create, bind, listen, accept, read, write and close sockets.

sys_socketcall's subroutine 'socket' expects 2 arguments - a pointer to an array of arguments in ECX and the integer value 1 in EBX. The sys_socketcall opcode is then loaded into EAX and the kernel is called to create the socket. Because everything in linux is a file, we recieve back the file descriptor of the created socket in EAX. This file descriptor can then be used for performing other socket programming functions.

Note: XORing a register by itself is an efficent way of ensuring the register is initalised with the integer value zero and doesn't contain an unexpected value that could corrupt your program.
