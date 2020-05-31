NULL terminating bytes
Ok so why did our second message print twice when we only called our sprint function on msg2 once? Well actually it did only print once. You can see what I mean if you comment out our second call to sprint. The output will be both of our message strings.

But how is this possible?

What is happening is we weren't properly terminating our strings. In assembly, variables are stored one after another in memory so the last byte of our msg1 variable is right next to the first byte of our msg2 variable. We know our string length calculation is looking for a zero byte so unless our msg2 variable starts with a zero byte it keeps counting as if it's the same string (and as far as assembly is concerned it is the same string). So we need to put a zero byte or 0h after our strings to let assembly know where to stop counting.

Note: In programming 0h denotes a null byte and a null byte after a string tells assembly where it ends in memory.
