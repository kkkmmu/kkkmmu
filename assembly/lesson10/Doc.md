Count to 10
Firstly, some background
Counting by numbers is not as straight forward as you would think in assembly. Firstly we need to pass sys_write an address in memory so we can't just load our register with a number and call our print function. Secondly, numbers and strings are very different things in assembly. Strings are represented by what are called ASCII values. ASCII stands for American Standard Code for Information Interchange. A good reference for ASCII can be found here. ASCII was created as a way to standardise the representation of strings across all computers.

Remember, we can't print a number - we have to print a string. In order to count to 10 we will need to convert our numbers from standard integers to their ASCII string representations. Have a look at the ASCII values table and notice that the string representation for the number '1' is actually '49' in ASCII. In fact, adding 48 to our numbers is all we have to do to convert them from integers to their ASCII string representations.

Writing our program
What we will do with our program is count from 1 to 10 using the ECX register. We will then ADD 48 to our counter to convert it from a number to it's ASCII string representation. We will then PUSH this value to the stack and call our print function passing ESP as the memory address to print from. Once we have finished counting to 10 we will exit our counting loop and call our quit function.
