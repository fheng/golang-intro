## Goals of this lesson

- Structs form a large part of the Golang type system. In this chapter we will go into more details about how you can use structs in Golang
including adding behaviour, embedding structs to "inherit" / compose objects and the difference between value recievers and pointer recievers.

- Introduce interfaces and how they work and show how they can help with encapsulation and testing.

- Talk about how Golang treats errors and how you should deal with them.


## Structs

As already shown structs form the main basis for creating our own custom types. They are similar, although quite different, from classes.

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

func main (){
    mt := MyType{Name:"Janeway"}
    mt.SayName()
}

```

[playground](https://play.golang.org/p/4rQH9jLY2j)

So what's happening here? Well we are creating and assigning a value of type MyType to the variable mt. 
Then we are calling the method SayName on that Value. This is actually a bit of sugar from Golang. 
Under the covers this is essentially what is happening:

```go 
package main

import "fmt"

type MyType struct{
    Name string
}

func SayName(mt MyType){
    fmt.Println(mt.Name)
}

func SayAgain(mt *MyType){
    fmt.Println(mt.Name)
}

func main(){
    mt := MyType{Name:"Janeway"}
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
Struct values when used as method recievers are exactly that values and so are read only.
Pointer recievers are mutable as you are recieving a pointer to the same address in memory.

```go 
package main

import "fmt"

type test_struct struct {
	Message string
}

//value reciever
func (t test_struct)Say (){
	fmt.Println(t.Message)
}
//value reciever
func (t test_struct)Update(m string){
	t.Message = m;
}
//pointer reciever
func (t * test_struct) SayP(){
	fmt.Println(t.Message)
}
//pointer reciever
func (t* test_struct) UpdateP(m string)  {
	t.Message = m;
}

func main(){
    ts := test_struct{}
    //call update on the value reciever
    ts.Update("test");
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
    tsp.SayP();
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

ts.Update(test_struct{},"test")

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
## Interfaces 
  TODO

## Errors 

 TODO 
