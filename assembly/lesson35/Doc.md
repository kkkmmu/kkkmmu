Sockets - Close
In this lesson we will use sys_close to properly close the active socket connection in the child process after our response has been sent. This will free up some resources that can be used to accept new incoming connections.

sys_close expects 1 argument - the file descriptor in EBX. The sys_close opcode is then loaded into EAX and the kernel is called to close the socket and remove the active file descriptor.

Note: Run the program and use the command curl http://localhost:9001 in another terminal or connect to the same address using any standard web browser.
