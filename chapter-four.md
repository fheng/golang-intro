## Goals of this lesson

- Structs form a large part of the Go type system. In this chapter we will go into more details about how you can use structs in Golang
including adding behaviour, embedding structs to "inherit" / compose objects and the difference between value recievers and pointer recievers.

- Introduce interfaces and how they work and show how they can help with encapsulation and testing.

- Talk about how Golang treats errors and how you should deal with them.


## Structs

As already shown structs form the main basis for creating our own custom types. They are similar, although quite different, from classes.
It is worth noting though that they are not the only way to create custom types below shows how you can use builtin types to create new types:

```go 
type MyString string

type IP string
type HostName string

type router map[HostName]IP

type HostList []Hostname

```

Through out this chapter we will focus on structs, but alot of what we cover here applies to any custom type.

```go 
type MyStruct struct{}
```
This is the simplist form of struct. It's not very interesting and doesn't really do anything.

- Adding behaviour or methods to a struct.

To add some behaviour to a struct you can have it be a method reciever. You can use a value reciever or a pointer reciever.

```go
package main 

import "fmt"

type MyType struct{
    Name string
}
//value reciever
func (mt MyType)SayName(){
    fmt.Println(mt.Name)
}
//pointer reciever
func (mt *MyType)SayAgain(){
    fmt.Println(mt.Name)
}

//example of custom builtin type 
type MySlice []string 
func (MySlice)Length()int{
    return len(MySlice)
}

func main (){
    mt := MyType{Name:"Janeway"}
    mt.SayName()

    ms := MySlice{"test"}
    fmt.Println(ms.Length())

}
```

[playground](https://play.golang.org/p/n_H2ZSMMdc)

So what's happening here? Well we are creating and assigning a value of type MyType to the variable mt. 
Then we are calling the method SayName on that Value. This is actually a bit of sugar from Go. 
Under the covers this is essentially what is happening:

```go 
package main

import "fmt"

type MyType struct {
	Name string
}

func SayName(mt MyType) {
	fmt.Println(mt.Name)
}

func SayAgain(mt *MyType) {
	fmt.Println(mt.Name)
}

func main() {
	mt := MyType{Name: "Janeway"}
	SayName(mt)
	SayAgain(&mt)
}

```
[playground](https://play.golang.org/p/p4qpZxSUDG)

So essentially the compiler is using the the definition before the method name as the first argument to the function.
As with Javascript you can assign functions to most types. For e.g:

```go 
package main 

import "fmt"

type MyMap map[string]string 

func(m MyMap)Get(val string)string{
    return m[val]
}

func(m MyMap)Put(key,val string){
    m[key] = val
}

func main(){
    mMap := MyMap{}
    mMap.Put("voyager","Janeway")
    fmt.Println(mMap.Get("voyager"))
    //or as MyMap is a map we could just do 
    fmt.Println(mMap["voyager"])
}

```
Struct values when used as method receivers are exactly that values and so are a shallow copy of the value allocated a new portion of memory.
The effects are not observed outside of the method as there are no references to the new value and so it is garbaged collected.
Pointer receivers allow mutation of what the pointer points to. Your function is recieving a pointer to the same address in memory even in the function stackframe.

```go 
package main

import "fmt"

type test_struct struct {
	Message string
}

//value reciever
func (t test_struct) Say() {
	fmt.Println(t.Message)
}

//value reciever
func (t test_struct) Update(m string) {
	t.Message = m
}

//pointer reciever
func (t *test_struct) SayP() {
	fmt.Println(t.Message)
}

//pointer reciever
func (t *test_struct) UpdateP(m string) {
	t.Message = m
}

func main() {
	ts := test_struct{}
	//call update on the value reciever
	ts.Update("test")
	ts.Say()
	//call update on the value reciever
	ts.Update("test2")
	ts.Say()
	//assign to the ts value created in this function
	ts.Message = "test2"
	ts.Say()
	//create a pointer
	tsp := &test_struct{}
	//set a value on the pointer
	tsp.Message = "test"
	tsp.SayP()
	//update the pointer
	tsp.UpdateP("test2")
	tsp.SayP()
}

//outputs 
//
//
//test2
//test
//test2
```
Recall that essentially what is happening here is 

```go 

Update(ts,"test")

```
So we are being passed a value and the value is being copied to the funtions stackframe. 
This is the same as if we were to do the following in Javascript :

```Javascript
var str = "mystring";

function add(mystr){
    mystr += "hasnt changed";
}

add(str);

console.log(str); //outputs mystring

```
More reading. The Go FAQ has some good information if you are interested in reading further:
[FAQ](https://golang.org/doc/faq#methods_on_values_or_pointers)

## Interfaces 

An interface type is defined by a set of methods.
``` go 
type CaptainNamer interface{
    Name()string
}
```
In Go, interfaces are implicitly implemented. There is no need for a ```implements ``` keyword.
interface types can be passed to functions in the same way as any other type.

```go
package main

import "fmt"

type CaptainNamer interface{
    Name()string
} 

type Captain struct{}
func (c Captain)Name()string{ //implementing this methods turns captain into something that satisfies the CaptainNamer interface
    return "Janeway"
}

func NameMe(cn CaptainNamer){
    fmt.Println(cn.Name())
}

func main(){
    c := Captain{}
    NameMe(c)
}

```
[playground](https://play.golang.org/p/4sO2iY3gPP)

This property of interfaces makes them extremely powerful and perhaps one the best features in Go. Lets have another example to try and make it as clear as possible.
Lets say you want to write a function that takes a file and writes some data to it. It might look like this:

```go 
func WriteData(f *os.File)error{
    written, err := f.Write([]byte("some data"))
    if err != nil{
        return err
    }
    fmt.Printf("I wrote %d ", written)
}
```
So we see that the function does what we expect, but now how do we test this? Well we have to create a file and write to it. Is there a better way?

Interfaces to the rescue:

There are many builtin interfaces within the stdlib. It turns out there is one called io.Writer 
```go
//from io package
type Writer interface {
        Write(p []byte) (n int, err error)
}
```
This interface along with the io.Reader interface are among the most commonly used in Go. So how can we use it?

Well our method takes a file. Turns out it has a method called Write that return the number of bytes written and an error if something goes wrong.
So the file is implicitly an io.Writer. So we can refactor our function to look like this:

```go
func WriteData(w io.Writer)error{
    written, err := w.Write([]byte("some data"))
    if err != nil{
        return err
    }
    fmt.Printf("I wrote %d ", written)
} 


main(){
    f ,_ := os.Create("/tmp/tmpfile") //ignoring errors as just example
    WriteData(f)
}
```

Know to test our function, we only need provide an implementation of io.Writer 

```go 

package main

import (
	"bytes"
	"testing"
)

func TestWriteData(t *testing.T) {
	//buffer is a writer
	var b bytes.Buffer
	if err := WriteData(&b); err != nil {
		t.Fail()
	}
    written := string(b.Bytes()) //casting the bytes to a string
    if "some data" != written{
        t.Fail()
    }
}
```

So a guideline that is spoken about in Go is to "accept interfaces return structs". 
Accepting interfaces allows us to easily unit test our work by providing alternative implementations

One final word on interfaces for now. We haven't covered this in structs yet, but you can embed interfaces within each other allowing for something like 
```extending``` 

```go 

//from the stdlib again 
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReadWriter interface {
	Reader
	Writer
}

```  

So here we are defing some very simple interfaces and then composing them together to form a completely new iterface.
If something wants to fulfill the ReaderWriter interface it must implment both Read and Write. We will go over embeding structs in the next lesson.

## Errors 
So we were talking about interfaces. Error is by far the most used interface in Go. lets have a look at it:

```go 
type error interface {
        Error() string
}
//The error built-in interface type is the conventional interface for representing an error condition, with the nil value representing no error.
``` 
So any type that implements the Error method is by default an error type. Some Examples:

```go

package main

import "fmt"

type MyError string
func (me MyError)Error()string{
    return string(me) //cast it back to string
}

type MyOtherError struct{
    Code int 
    Message string
}
func (me MyOtherError)Error()string{
    return fmt.Sprintf("error occurred: %s with code %d ",me.Message,me.Code)
}

func MyBadFunc()error{
    return MyError("this is an error ") //cast it to a MyError
}

func MyOtherBadFunc()error{
    return MyOtherError{Message:"MyOtherBadFunc failed ", Code:500}
}

func main(){
    fmt.Println(MyBadFunc())

    fmt.Println(MyOtherBadFunc())
}

```

[playground](https://play.golang.org/p/07vBhhV6xx)

In Go errors are values and are meant to be checked. They are not exceptions. The idiom in Go is very similar to node.js 
This is likely the most common code snippet in Go programs
```go 
if err != nil{
    return err
}

```
This is likely the most commons snippet of code in node programs

```Javascript
if (err){
  return cb(err);
}
```

Go tries to ensure you check your errors and you should always check your errors.
Error states are very normal in running programs, they are part of your program not something that should just disappear.

One note on errors. As the are values you do not naturally get a stacktrace. I think Go could have done a better job here. However there
are solutions. One of the go contributors wrote a compatible errors package github.com/pkg/errors I recommend using this package once you become a little more advanced.

```go 

go get github.com/pkg/errors 

```

```go 
package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	err := errors.New("my error with stack")
	fmt.Printf("%+v", err) //the +v indicates to print the trace
}

```

## A side project 
I have created a repo:

github.com/feedhenry/rhm 

It is a replacement for fhc. Not the full thing as it is far too big. I want to use this project as a learning and trainging ground. 
The hope is that this is something tangible that you can get your teeth into. Questions, issues PRs all welcome.
