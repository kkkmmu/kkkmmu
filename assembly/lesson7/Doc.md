Linefeeds
Linefeeds are essential to console programs like our 'hello world' program. They become even more important once we start building programs that require user input. But linefeeds can be a pain to maintain. Sometimes you will want to include them in your strings and sometimes you will want to remove them. If we continue to hard code them in our variables by adding 0Ah after our declared message text, it will become a problem. If there's a place in the code that we don't want to print out the linefeed for that variable we will need to write some extra logic remove it from the string at runtime.

It would be better for the maintainability of our program if we write a subroutine that will print out our message and then print a linefeed afterwards. That way we can just call this subroutine when we need the linefeed and call our current sprint subroutine when we don't.

A call to sys_write requires we pass a pointer to an address in memory of the string we want to print so we can't just pass a linefeed character (0Ah) to our print function. We also don't want to create another variable just to hold a linefeed character so we will instead use the stack.

The way it works is by moving a linefeed character into EAX. We then PUSH EAX onto the stack and get the address pointed to by the Extended Stack Pointer. ESP is another register. When you PUSH items onto the stack, ESP is decremented to point to the address in memory of the last item and so it can be used to access that item directly from the stack. Since ESP points to an address in memory of a character, sys_write will be able to use it to print.

Note: I've highlighted the new code in functions.asm below.
