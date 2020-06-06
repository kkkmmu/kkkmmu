Sockets - Accept
In the previous lessons we created a socket and used the 'bind' subroutine to associate it with a local IP address and port. We then used the 'listen' subroutine of sys_socketcall to tell our socket to listen for incoming TCP requests. Now we will use the 'accept' subroutine of sys_socketcall to tell our socket to accept those incoming requests. Our socket will then be ready to read and write to remote connections.

sys_socketcall's subroutine 'accept' expects 2 arguments - a pointer to an array of arguments in ECX and the integer value 4 in EBX. The sys_socketcall opcode is then loaded into EAX and the kernel is called. The 'accept' subroutine will create another file descriptor, this time identifying the incoming socket connection. We will use this file descriptor to read and write to the incoming connection in later lessons.

Note: Run the program and use the command sudo netstat -plnt in another terminal to view the socket listening on port 9001.
