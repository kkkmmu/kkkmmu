Sockets - Bind
Building on the previous lesson we will now associate the created socket with a local IP address and port which will allow us to connect to it. We do this by calling the second subroutine of sys_socketcall which is called 'bind'.

We begin by storing the file descriptor we recieved in lesson 29 into EDI. EDI was originally called the Destination Index and is traditionally used in copy routines to store the location of a target file.

sys_socketcall's subroutine 'bind' expects 2 arguments - a pointer to an array of arguments in ECX and the integer value 2 in EBX. The sys_socketcall opcode is then loaded into EAX and the kernel is called to bind the socket.
