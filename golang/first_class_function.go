What are first class functions?
A language which supports first class functions allows functions to be assigned to variables, passed as arguments to other functions and returned from other functions. Go has support for first class functions.

In this tutorial we will discuss about the syntax and various use cases of first class functions.

Anonymous functions
Let's start with a simple example which assigns a function to a variable.

package main

import (  
	"fmt"
)

func main() {  
	a := func() {
		fmt.Println("hello world first class function")
	}
	a()
	fmt.Printf("%T", a)
}
Run in playground

In the program above, we have assigned a function to the variable a in line no. 8. This is the syntax for assigning a function to a variable. If you notice carefully, the function assigned to a does not have a name. These kind of functions are called anonymous functions since they do not have a name.

The only way to call this function is using the variable a. We have done this in the next line. a() calls the function and this prints hello world first class function. In line no. 12 we print variable a's type . This will print func().

Running this program will output

hello world first class function  
func()  
It is also possible to call a anonymous function without assigning it to a variable. Let's see how this is done in the following example.

package main

import (  
	"fmt"
)

func main() {  
	func() {
		fmt.Println("hello world first class function")
	}()
}
Run in playground

In the program above, an anonymous function is defined in line no. 8 and immediately after the function definition, we call the function using () in line no. 10. This program will output,

hello world first class function  
It is also possible to pass arguments to anonymous functions just like any other function.

package main

import (  
	"fmt"
)

func main() {  
	func(n string) {
		fmt.Println("Welcome", n)
	}("Gophers")
}
Run in playground

In the program above, a string argument is passed to the anonymous function in line no. 10. Running this program will print,

Welcome Gophers  
User defined function types
Just like we define our own struct types, it is possible to define our own function types.

type add func(a int, b int) int  
The code snippet above creates a new function type add which accepts two integer arguments and returns a integer. Now we can define variables of type add.

Let's write a program which defines a variable of type add.

package main

import (  
	"fmt"
)

type add func(a int, b int) int

func main() {  
	var a add = func(a int, b int) int {
		return a + b
	}
	s := a(5, 6)
	fmt.Println("Sum", s)
}
Run in playground


Get the free Golang tools cheat sheet

In the program above, in line no.10, we define a variable a of type add and assign to it a function whose signature matches the type add. We call the function in line no. 13 and assign the result to s. This program will print,

Sum 11  
Higher-order functions
The definition of Higher-order function from wiki is a function which does at least one of the following

takes one or more functions as arguments
returns a function as its result
Let's look at some simple examples for the above two scenarios.

Passing functions as arguments to other functions
package main

import (  
	"fmt"
)

func simple(a func(a, b int) int) {  
	fmt.Println(a(60, 7))
}

func main() {  
	f := func(a, b int) int {
		return a + b
	}
	simple(f)
}
Run in playground

In the above example, in line no. 7 we define a function simple which takes a function which accepts two int arguments and returns a int as a parameter. Inside the main function in line no. 12 we create a anonymous function f whose signature matches the parameter of the function simple. We call simple and pass f as an argument to it in the next line. This program prints 67 as output.

Returning functions from other functions
Now let's rewrite the program above and return a function from the simple function.

package main

import (  
	"fmt"
)

func simple() func(a, b int) int {  
	f := func(a, b int) int {
		return a + b
	}
	return f
}

func main() {  
	s := simple()
	fmt.Println(s(60, 7))
}
Run in playground

In the program above, the simple function in line no 7 returns a function that takes two int arguments and returns a int argument.

This simple function is called from line no. 15. The return value from simple is assigned to s. Now s contains the function returned by simple function. We call s and pass it two int arguments in line no. 16. This program outputs 67.

Closures
Closures are a special case of anonymous functions. Closures are anonymous functions which access the variables defined outside the body of the function.

An example will make things more clear.

package main

import (  
	"fmt"
)

func main() {  
	a := 5
	func() {
		fmt.Println("a =", a)
	}()
}
Run in playground

In the program above, the anonymous function accesses the variable a which is present outside its body in line no. 10. Hence this anonymous function is a closure.

Every closure is bound to its own surrounding variable. Let's understand what this means by using a simple example.

package main

import (  
	"fmt"
)

func appendStr() func(string) string {  
	t := "Hello"
	c := func(b string) string {
		t = t + " " + b
		return t
	}
	return c
}

func main() {  
	a := appendStr()
	b := appendStr()
	fmt.Println(a("World"))
	fmt.Println(b("Everyone"))

	fmt.Println(a("Gopher"))
	fmt.Println(b("!"))
}
Run in background

In the program above, the function appendStr returns a closure. This closure is bound to the variable t. Let's understand what this means.

The variables a and b declared in line nos. 17, 18 are closures and they are bound to their own value of t.

We first call a with the parameter World. Now the value of a's version of t becomes Hello World.

In line no. 20 we call b with the parameter Everyone. Since b is bound to its own variable t, b's version of t has a initial value of Hello again. Hence after this function call, the value of b's version of t becomes Hello Everyone. The rest of the program is self explanatory.

This program will print,

Hello World  
Hello Everyone  
Hello World Gopher  
Hello Everyone !  

Get the free Golang tools cheat sheet

Practical use of first class functions
Till now we have defined what first class functions are and we have seen a few contrived examples to learn how they work. Now lets write a concrete program which shows the practical use of first class functions.

We will create a program which filters a slice of students based on some criteria. Let's approach this step by step.

First lets define the student type.

type student struct {  
	firstName string
	lastName string
	grade string
	country string
}
The next step is to write the filter function. This function takes a slice of students and a function which determines whether a student matches the filtration criteria as parameters. We will understand better once we write this function. Let's go ahead and do it.

func filter(s []student, f func(student) bool) []student {  
	var r []student
	for _, v := range s {
		if f(v) == true {
			r = append(r, v)
		}
	}
	return r
}
In the above function, the second parameter to filter is a function which takes a student as parameter and returns a bool. This function determines whether a particular student matches a criteria or not. We iterate through the student slice in line no. 3 and and we pass each student as parameter to the function f. If this returns true, it means that that the student has passed the filter criteria and he is added to the result slice r. You might be a little confused about the real use of this function, but it will be clear once we complete the program. I have added the main function and have provided the full program below.

package main

import (  
	"fmt"
)

type student struct {  
	firstName string
	lastName  string
	grade     string
	country   string
}

func filter(s []student, f func(student) bool) []student {  
	var r []student
	for _, v := range s {
		if f(v) == true {
			r = append(r, v)
		}
	}
	return r
}

func main() {  
	s1 := student{
		firstName: "Naveen",
		lastName:  "Ramanathan",
		grade:     "A",
		country:   "India",
	}
	s2 := student{
		firstName: "Samuel",
		lastName:  "Johnson",
		grade:     "B",
		country:   "USA",
	}
	s := []student{s1, s2}
	f := filter(s, func(s student) bool {
		if s.grade == "B" {
			return true
		}
		return false
	})
	fmt.Println(f)
}
Run in playground

In the main function, we first create two students s1 and s2 and add them to slice s. Now let's say we want to find out all students who have grade B. We have established this in the above program by passing a function which checks whether the student has grade B and if so then returning true, as parameter to the filter function in line no. 38. The above program will print,

[{Samuel Johnson B USA}]
Let's say we want to find all students from India. This can be done easily by changing the function parameter to the filter function. 
I have provided code that does this below,

c := filter(s, func(s student) bool {  
	if s.country == "India" {
		return true
	}
	return false
})
fmt.Println(c)  
Please add this to the main function and check the output.

Let's conclude this section by writing one more program. This program will perform the same operations on each element of a slice and return the result. For example if we want to multiply all integers in a slice by 5 and return the output, it can be easily done using first class functions. These kind of functions which operate on every element of a collection are called map functions. I have provided the program below. It is self explanatory.

package main

import (  
	"fmt"
)

func iMap(s []int, f func(int) int) []int {  
	var r []int
	for _, v := range s {
		r = append(r, f(v))
	}
	return r
}
func main() {  
	a := []int{5, 6, 7, 8, 9}
	r := iMap(a, func(n int) int {
		return n * 5
	})
	fmt.Println(r)
}
Run in playground

The above program will print,

[25 30 35 40 45]
Here's a quick recap of what we learnt in this tutorial,

What are first class functions?
Anonymous functions
User defined function types
Higher-order functions
Passing functions as arguments to other functions
Returning functions from other functions
Closures
Practical use of first class functions
That's about it for first class functions. Have a good day.
