 Subroutines
 Introduction to subroutines
 Subroutines are functions. They are reusable pieces of code that can be called by your program to perform various repeatable tasks. Subroutines are declared using labels just like we've used before (eg. _start:) however we don't use the JMP instruction to get to them - instead we use a new instruction CALL. We also don't use the JMP instruction to return to our program after we have run the function. To return to our program from a subroutine we use the instruction RET instead.

 Why don't we JMP to subroutines?
 The great thing about writing a subroutine is that we can reuse it. If we want to be able to use the subroutine from anywhere in the code we would have to write some logic to determine where in the code we had jumped from and where we should jump back to. This would litter our code with unwanted labels. If we use CALL and RET however, assembly handles this problem for us using something called the stack.

 Introduction to the stack
 The stack is a special type of memory. It's the same type of memory that we've used before however it's special in how it is used by our program. The stack is what is call Last In First Out memory (LIFO). You can think of the stack like a stack of plates in your kitchen. The last plate you put on the stack is also the first plate you will take off the stack next time you use a plate.

 The stack in assembly is not storing plates though, its storing values. You can store a lot of things on the stack such as variables, addresses or other programs. We need to use the stack when we call subroutines to temporarily store values that will be restored later.

 Any register that your function needs to use should have it's current value put on the stack for safe keeping using the PUSH instruction. Then after the function has finished it's logic, these registers can have their original values restored using the POP instruction. This means that any values in the registers will be the same before and after you've called your function. If we take care of this in our subroutine we can call functions without worrying about what changes they're making to our registers.

 The CALL and RET instructions also use the stack. When you CALL a subroutine, the address you called it from in your program is pushed onto the stack. This address is then popped off the stack by RET and the program jumps back to that place in your code. This is why you should always JMP to labels but you should CALL functions.
