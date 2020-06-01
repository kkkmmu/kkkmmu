Calculator - division
In this program we will be dividing the value in EBX by the value present in EAX. We've used division before in our integer print subroutine. Our program requires a few extra strings in order to print out the correct answer but otherwise there's nothing complicated going on.

Firstly we load EAX and EBX with integers in the same way as Lesson 12. Division logic is performed using the DIV instruction. The DIV instruction always divides EAX by the value passed after it. It will leave the quotient part of the answer in EAX and put the remainder part in EDX (the original data register). We then MOV and call our strings and integers to print out the correct answer.
