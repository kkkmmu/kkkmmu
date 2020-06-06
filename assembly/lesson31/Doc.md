Sockets - Listen
In the previous lessons we created a socket and used the 'bind' subroutine to associate it with a local IP address and port. In this lesson we will use the 'listen' subroutine of sys_socketcall to tell our socket to listen for incoming TCP requests. This will allow us to read and write to anyone who connects to our socket.

sys_socketcall's subroutine 'listen' expects 2 arguments - a pointer to an array of arguments in ECX and the integer value 4 in EBX. The sys_socketcall opcode is then loaded into EAX and the kernel is called. If succesful the socket will begin listening for incoming requests.
