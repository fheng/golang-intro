* Note a lot of this is taken from the excellent tour.golang.org
* Some information was taken from http://www.golangbootcamp.com/

##Goals of this lesson
- Introduce some basic building blocks to allow the exploration of creating actual programs.
- Introduce how Go deals with data structures.
- Introduce how Go does json encoding
- Show how to write a basic web server and test it.


The examples can be run by adding the code to a main.go and running:

```bash

go run main.go 

```

## Exported Names

In Golang a packages' publicly accessible properties all begin with a capital letter. From within a package namespace you can
refer to private functions and variables (starting with lower case) but from outside the package, you can only access the things exported from that package.
Think public and private key word in Java.

```go 
//exported 

func MyPublicFunc(){

}

//private 
func myPrivateFunc(){

}

```

## Functions

Functions
A function can take zero or more arguments.

In this example, add takes two parameters of type int.

```go 
func add(x int, y int) int {
	return x + y
}
```

Notice that the type comes after the variable name.

When two or more consecutive named function parameters share a type, you can omit the type from all but the last.

```go 
func add(x, y int) int {
	return x + y
}
```

A function can return any number of results to do this you add braces around the return types:

```go 
func xandy(x,y int)(int,int){ //returns two ints
    return x,y
}

```

Lets try this out: 

Create a new dir in your gopath
```
mkdir -p $GOPATH/src/github.com/YOUR_USER/functions
touch $GOPATH/src/github.com/YOUR_USER/functions/main.go
```

Add the following to main.go

```go 

package main

import "fmt"

func swap(x, y string) (string, string) {
	return y, x
}

func main() {
    //notice the assignment here. This is called a short assignment and it inferrs the type. More on this later.
	//also notice we have 2 variables on the left hand side
	a, b := swap("hello", "world") 
	fmt.Println(a, b)
}

```

Just as in javascript, functions have closure and can be passed as values:

```go 
package main

import (
	"fmt"
	"math"
)

func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func main() {
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}
	fmt.Println(hypot(5, 12)) //call hypot directly
	fmt.Println(compute(hypot)) //pass hypot as value to compute
	fmt.Println(compute(math.Pow)) //pass math.Pow as a value
	fmt.Println(func(x,y float64)float64{
		return hypot(x,y)
	}) //define a function inline
}

```
## Variables and infferred types

There are two main types of assignment in Go. 
1) the ``` var ``` key word and the ```:=``` short assignment.

#### VAR:

The var keyword can be used both outside function scopes and within.

```go 
package main 

import(
	"fmt"
)

var MyGlobal string = "hello world"


func main (){
	var msg = "hello world again"
	fmt.Println(MyGlobal)
	fmt.Println(msg)
}

``` 

In the above example we set the type, but if an initializer is present, the type can be omitted, the variable will take the type of the initializer (inferred typing).
The below outlines the different ways you can use the var keyword:

```go 
package main

import "fmt"

func main(){
	//as a list of vars
  var(
	  one = 1 //infferred type 
	  two = 2 // infferred type
	  three int // default value of type int 0
  )
  //individual declarations
  var four int // default value of type int 0
  var five = 5 // inferred type 

  // multiple inline initialization is also fine
  var name, location, age = "golang", "is great", 42 

  fmt.Printf("one %d two %d three %d four %d five %d", one,two,three,four,five)

}

```
The above will ouput 

```
one 1 two 2 three 0 four 0 five 5
```

Notice that three and four have a value of 0 that is because the default value for int is 0

#### Short Assignment:

Inside a function, the := short assignment statement can be used.
Unlike regular variable declarations, a short variable declaration may redeclare variables provided they were originally declared earlier in the same block
Here are some example of short assignment:

```go 

func main (){
	msg := "hello world again"
	msg,age,fingers := "john",27,10 // notice we redeclare msg
}
```

A variable can contain any type, including functions:

```go 

package main 

import "fmt"

func main(){
	fn := func (){
		fmt.Println("hello world")
	}
	fn()
}
``` 
These things allow go to feel quite familar to javascript developers.


### The defer keyword

The defer keyword is a useful piece of utility provided by the golang language. It allows you to defer a function call until the current function has completed
this can be extremely useful for doing things like closing files, closing sockets etc. Instead of closing files etc on err and on success instead you can defer 
its close function and be sure your function will clean up after itself. (we will be using this later) 

```go 
package main 

import "fmt"

func hello1(){
	defer hello2() //this will always run when hello1 has completed
	fmt.Println("hello world 1")
}

func hello2(){
	fmt.Println("hello world 2")
}

func main(){
	hello1()
}
``` 

### Structs

Structs are the basic building blocks of your data structures in Go. They look a lot like javascript objects:

```go 
package main 

import "fmt"

type MyType struct {
	Name string
}

func main (){
	myValue := MyType{Name:"Janeway"} //this will give myValue a value of type MyType
	fmt.Println(myValue.Name)

	myPointer := &MyType{Name:"Piccard"}
	fmt.Println(myPointer.Name) // this will give myPointer the address of where MyType is stored.

} 
```

We will go more into pointers and values in Golang in a future lesson. As with most languages, pointers are much cheaper and means anything recieving that
pointer can change attributes of whatever is stored in the pointer's address. With a value your function will get a copy of the value so it can change
the attributes of that value within the scope of the function but not outside of that scope.  


### Encoding (json)

Golang has several encoding packages built into the stdlib. Today we will focus on the json encoding package. The package details can be seen here:
[json](https://golang.org/pkg/encoding/json/)
There are not alot of methods in this package. The main ones we are interested in are Decode and Encode. Some code here may be a little confusing but try not to be put off as it will become clearer.

```go 
package main

import (
        "encoding/json" //import the json sub package of the encoding package.
        "fmt"
        "log"     // import the logging package
        "os"      //import the operating system package
        "strings" // import the strings package
)

type MyData struct { //define a data structure
        Message string
}

var jsonMyData = `{"Message":"hello"}` //define a string literal

//take a value of MyData encode it print it and return an error if there is a problem
func encodeAndPrintMyData(data MyData) error {
        enc := json.NewEncoder(os.Stdout) //NewEncoder takes a writer (we will get into this in the future) os.Stdout is a writer for stdout
        return enc.Encode(data) //enc.Encode returns an error if there was a problem.
}

//take a json string and decode it into a MyData and return the pointer return an error if there is a problem
func decodeMyData(jsonData string) (*MyData, error) {
        dec := json.NewDecoder(strings.NewReader(jsonData))
        myData := &MyData{} //create something for the data to be decoded into
        err := dec.Decode(myData)
        if err != nil {
                return nil, err //remember we can return more than one value
        }
        return myData, nil
}

func main() {
        myValue := MyData{Message: "hello value"}
        err := encodeAndPrintMyData(myValue)
        if err != nil {
                log.Fatal(err)
        }
        myPointer, err := decodeMyData(jsonMyData)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Println(myPointer)
}
```

Golangs json encoder and decoder support streaming also. See 
[https://golang.org/pkg/encoding/json/#example_Decoder_Decode_stream](https://golang.org/pkg/encoding/json/#example_Decoder_Decode_stream)


### Putting it all together the first web server

Create a new project:

```
mkdir $GOPATH/src/github/YOUR_USER/api
touch $GOPATH/src/github/YOUR_USER/api/main.go
```

Add the following to main.go

```go 
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

//Message wraps a message and stamps it
type Message struct {
        Message string `json:"message"` //tells the decoder what to decode into
        Stamp   int64  `json:"stamp,omitempty"`
}

//BuisnessLogic does awesome BuisnessLogic
func BuisnessLogic(text string) *Message {
        mess := &Message{}
        mess.Message = text
        mess.Stamp = time.Now().Unix()
        return mess
}

//Echo echoes what you send
func Echo(res http.ResponseWriter, req *http.Request) {
        var (
                jsonDecoder = json.NewDecoder(req.Body) //decoder reading from the post body
                jsonEncoder = json.NewEncoder(res)      //encoder writing to the response stream
                message     = &Message{}         // something to hold our data
        )
        res.Header().Add("Content-type", "application/json")
        if err := jsonDecoder.Decode(message); err != nil { //decode our data into our struct
                res.WriteHeader(http.StatusInternalServerError)
                return
        }
        pointless := BuisnessLogic(message.Message)
        if err := jsonEncoder.Encode(pointless); err != nil { //encode our data and write it back to the response stream
                res.WriteHeader(http.StatusInternalServerError)
                return
        }
}

//Setup our simple router
func router() http.Handler {
        http.HandleFunc("/api/echo", Echo)
        return http.DefaultServeMux //this is a stdlib http.Handler
}

func main() {
        router := router()
        //start our server on port 3001
        if err := http.ListenAndServe(":3001", router); err != nil {
                log.Fatal(err)
        }
}
```

So now we have written our webserver now lets run it.

```
go run main.go 

curl -XPOST http://localhost:3001/api/echo -H 'Content-type:application/json' -d '{"message":"hello world"}'

```

Of course we need to also test it. Here we will use the stdlib [httptest](https://golang.org/pkg/net/http/httptest/) package. Normally we would try to write the tests first, but I didn't want to overload.
In the same directory create main_test.go 

```
touch $GOPATH/src/github/YOUR_USER/api/main_test.go
```

Add the following to the test file:

```go 

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEcho(t *testing.T) {
	server := httptest.NewServer(router()) //starts a real server on a free port
	defer server.Close() //notice we use defer here to ensure our server is closed
	res, err := http.NewRequest("POST", server.URL+"/api/echo", strings.NewReader(`{"message":"test"}`))
	if err != nil {
		log.Fatal(err)
	}
	echo, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	message := &Message{}
	if err := json.Unmarshal(echo, message); err != nil {
		log.Fatal(err)
	}
	if "test" != message.Message {
		t.Fail()
		log.Println("expected the message to equal test")
	}
}

``` 

Finally run the test command. 

```bash 

go test -cover -v 

```

Notice we added the -cover flag, this will print the coverage stats for the package.

Also it is worth noticing that we still haven't used any dependencies yet. We will eventually need them, but it good to show what you can achieve just with the stdlib.

### Finally

What would people like to build as part of this? Would working on a command line client that maybe does a login against RHMAP be good or would you prefer to learn more about building rest api servers?


