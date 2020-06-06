Sockets - Write
When an incoming connection is accepted by our socket, a new file descriptor identifying the incoming socket connection is returned in EAX. In this lesson we will use this file descriptor to send our response to the connection.

We will use the kernel function sys_write to write to the incoming socket connection. As our socket will be communicating using the HTTP protocol, we will need to send some compulsory headers in order to allow HTTP speaking clients to connect. We will send these following the formatting rules set out in the RFC Standard.

sys_write expects 3 arguments - the number of bytes to write in EDX, the response string to write in ECX and the file descriptor in EBX. The sys_write opcode is then loaded into EAX and the kernel is called to send our response back through our socket to the incoming connection.

Note: We will create a variable in the .data section to store the response we will write to the file descriptor. See Lesson 1 for more information on the .data section.

Note: Run the program and use the command curl http://localhost:9001 in another terminal to view the response sent via our socket. Or connect to the same address using any standard web browser.
