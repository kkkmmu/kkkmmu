Sockets - Read
When an incoming connection is accepted by our socket, a new file descriptor identifying the incoming socket connection is returned in EAX. In this lesson we will use this file descriptor to read the incoming request headers from the connection.

We begin by storing the file descriptor we recieved in lesson 32 into ESI. ESI was originally called the Source Index and is traditionally used in copy routines to store the location of a target file.

We will use the kernel function sys_read to read from the incoming socket connection. As we have done in previous lessons, we will create a variable to store the contents being read from the file descriptor. Our socket will be using the HTTP protocol to communicate. Parsing HTTP request headers to determine the length of the incoming message and accepted response formats is beyond the scope of this tutorial. We will instead just read up to the first 255 bytes and print that to standardout.

Once the incoming connection has been accepted, it is very common for webservers to spawn a child process to manage the read/write communication. The parent process is then free to return to the listening/accept state and accept any new incoming requests in parallel. We will implement this design pattern below using sys_fork and the JMP instruction prior to reading the request headers in the child process.

To generate valid request headers we will use the commandline tool curl to connect to our listening socket. But you can also use a standard web browser to connect in the same way.

sys_read expects 3 arguments - the number of bytes to read in EDX, the memory address of our variable in ECX and the file descriptor in EBX. The sys_read opcode is then loaded into EAX and the kernel is called to read the contents into our variable which is then printed to the screen.

Note: We will reserve 255 bytes in the .bss section to store the contents being read from the file descriptor. See Lesson 9 for more information on the .bss section.

Note: Run the program and use the command curl http://localhost:9001 in another terminal to view the request headers being read by our program.
