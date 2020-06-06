Lesson 36
Download a Webpage
In the previous lessons we have been learning how to use the many subroutines of the sys_socketcall kernel function to create, manage and transfer data through Linux sockets. We will continue that theme in this lesson by using the 'connect' subroutine of sys_socketcall to connect to a remote webserver and download a webpage.

These are the steps we need to follow to connect a socket to a remote server:

Call sys_socketcall's subroutine 'socket' to create an active socket that we will use to send outbound requests.
Call sys_socketcall's subroutine 'connect' to connect our socket with a socket on the remote webserver.
Use sys_write to send a HTTP formatted request through our socket to the remote webserver.
Use sys_read to recieve the HTTP formatted response from the webserver.
We will then use our string printing function to print the response to our terminal.
What is a HTTP Request
The HTTP specification has evolved through a number of standard versions including 1.0 in RFC1945, 1.1 in RFC2068 and 2.0 in RFC7540. Version 1.1 is still the most common today.

A HTTP/1.1 request is comprised of 3 sections:

   A line containing the request method, request url, and http version
   An optional section of request headers
   An empty line that tells the remote server you have finished sending the request and you will begin waiting for the response.
   A typical HTTP request for the root document on this server would look like this:

GET / HTTP/1.1                  ; A line containing the request method, url and version
Host: asmtutor.com              ; A section of request headers
                                ; A required empty line

Writing our program
This tutorial starts out like the previous ones by calling sys_socketcall's subroutine 'socket' to initially create our socket. However, instead of calling 'bind' on this socket we will call 'connect' with an IP Address and Port Number to connect our socket to a remote webserver. We will then use the sys_write and sys_read kernel methods to transfer data between the two sockets by sending a HTTP request and reading the HTTP response.

sys_socketcall's subroutine 'connect' expects 2 arguments - a pointer to an array of arguments in ECX and the integer value 3 in EBX. The sys_socketcall opcode is then loaded into EAX and the kernel is called to connect to the socket.

Note: In Linux we can use the following command ./crawler > index.html to save the output of our program to a file instead.
