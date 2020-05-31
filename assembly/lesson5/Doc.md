External include files
External include files allow us to move code from our program and put it into separate files. This technique is useful for writing clean, easy to maintain programs. Reusable bits of code can be written as subroutines and stored in separate files called libraries. When you need a piece of logic you can include the file in your program and use it as if they are part of the same file.

In this lesson we will move our string length calculating subroutine into an external file. We fill also make our string printing logic and program exit logic a subroutine and we will move them into this external file. Once it's completed our actual program will be clean and easier to read.

We can then declare another message variable and call our print function twice in order to demonstrate how we can reuse code.

Note: I won't be showing the code in functions.asm after this lesson unless it changes. It will just be included if needed.
