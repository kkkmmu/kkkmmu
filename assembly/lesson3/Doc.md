Firstly, some background
Why do we need to calculate the length of a string?

Well sys_write requires that we pass it a pointer to the string we want to output in memory and the length in bytes we want to print out. If we were to modify our message string we would have to update the length in bytes that we pass to sys_write as well, otherwise it will not print correctly.

You can see what I mean using the program in Lesson 2. Modify the message string to say 'Hello, brave new world!' then compile, link and run the new program. The output will be 'Hello, brave ' (the first 13 characters) because we are still only passing 13 bytes to sys_write as its length. It will be particularly necessary when we want to print out user input. As we won't know the length of the data when we compile our program, we will need a way to calculate the length at runtime in order to successfully print it out.

Writing our program
To calculate the length of the string we will use a technique called pointer arithmetic. Two registers are initialised pointing to the same address in memory. One register (in this case EAX) will be incremented forward one byte for each character in the output string until we reach the end of the string. The original pointer will then be subtracted from EAX. This is effectively like subtraction between two arrays and the result yields the number of elements between the two addresses. This result is then passed to sys_write replacing our hard coded count.

The CMP instruction compares the left hand side against the right hand side and sets a number of flags that are used for program flow. The flag we're checking is the ZF or Zero Flag. When the byte that EAX points to is equal to zero the ZF flag is set. We then use the JZ instruction to jump, if the ZF flag is set, to the point in our program labeled 'finished'. This is to break out of the nextchar loop and continue executing the rest of the program.
