## Goals of this lesson
- Look at the two main collection types slices and maps
- Add some new endpoints to our api using slices and maps
- Start working on a cli to interact with our new server

## Note this lesson builds on top of the code from lesson 2

We started a basic api web server in lesson 2 the code can be found here: As we will be building on this it is tagged as chapter-2
https://github.com/maleck13/api

- set up for this lesson 
```bash

cd $GOPATH/src/github.com/<YOUR_USER>
git clone git@github.com:maleck13/api.git
cd api 
go run main.go 

```
Feel free to fork the repo. I am also happy to review any changes you make.

## Collections types 

Golang has two main collection types.

1) maps.

2) slices.

### Maps 

Maps are somewhat similar to what other languages call “dictionaries” or “hashes”.
The keys of  maps can be any type that is comparable. You can read more about comparing on the Golang spec:
[comparing](https://golang.org/ref/spec#Comparison_operators). I will leave it as homework to go through this.
Maps are not naturally thread safe. We are not using concurrency yet, so I will go over making data safe for concurrent use when we get to that.

Maps can be created in two ways: 

1) Using map literals:

```go 

package main

import "fmt" 

func main (){
    m := map[string]string{ //here we are declaring the key to be of type string and the value to be of type string
        "Picard":"The next generation",
        "Janeway":"Voyager",
        "Kirk":"Star Trek",
    }
    fmt.Println(m)
}

```
[playground](https://play.golang.org/p/BVXeyQTMNV)

Again this is very similar to object creation in javascript. Remove the types and it is exactly the same.

Maps can contain any valid type. When using a custom type in a map, you can omit the type name in a map literal definition:

```go 
package main

import "fmt"

type Location struct {
	Lat, Lng float64
}

func main() {
    m := map[string]Location{
	"Bell Labs": {40.68433, -74.39967},
	"Google": {37.42202, -122.08408},
	}

    fmt.Println(m)
}
```
[playground](https://play.golang.org/p/JAVAZPlm5C)

2) With the make keyword:
The make built-in function allocates and initializes an object of type slice, map, or chan (only). 

```go 

package main 

import "fmt"

func main (){
    m := make(map[string]string,3) // The second argument here is the initial size of the map (optional).
    m["Picard"] = "The next generation"
    m["Janeway"] = "Voyager"
    m["Kirk"] = "Star Trek"
    fmt.Println(m)
}
```

[playground](https://play.golang.org/p/QdmVNGKxfH)

The only reason to use the make key word is if you want to decide the size of the map before hand.

Lets go through some common actions with maps.

- Adding values

See the above examples that show us adding values to maps.

- Checking if a key exists in a map 1

```go 
package main

import "fmt" 

func main (){
    m := map[string]string{
        "Picard":"The next generation",
        "Janeway":"Voyager",
        "Kirk":"Star Trek",
    }
    if v,ok := m["Picard"]; ok{ // in this case ok is a bool that will be true if the key exists.
        fmt.Println(v)
    }
    if v := m["Picard"]; v != ""{ // in this case you have to assert something about the returned value. I prefer the example above.
        fmt.Println(v)
    }
    if _,ok := m["not here"]; !ok{ //notice the _ here this is used to ignore a return value
           fmt.Println("no key with that name")
    }
}
``` 
[playground](https://play.golang.org/p/C1hfGP7ikv)

- Iterating through a map 

```go 
package main

import "fmt" 

func main (){
    m := map[string]string{ 
        "Picard":"The next generation",
        "Janeway":"Voyager",
        "Kirk":"Star Trek",
    }
    for key,val := range m{ //notice the range key word we will talk more about this
        fmt.Printf("map has key %s and value %s \n",key,val)
    }
}

```
[playground](https://play.golang.org/p/PPGnySKyCv)

*Range* 
So we saw the range keyword above. A range clause provides a way to iterate over an array, slice, string, map, or channel (more on channels in the future).
It is only used in conjunction with a for loop. One thing to note is that when you use the range key word, for each
loop, it copies the current key and value into the vars. This means that if you change the value in the loop it will not be reflected in the map or slice.

- Get the size of a map.

To find out the size of a map (i.e the number of keys it has) we use the builtin ```len``` function.  

```go 

package main

import "fmt" 

func main (){
    m := map[string]string{ 
        "Picard":"The next generation",
        "Janeway":"Voyager",
        "Kirk":"Star Trek",
    }
    fmt.Println(len(m)) //outputs 3
}
```
[playground](https://play.golang.org/p/2rYP28LdnZ)

- Deleting items from a map. 

To Delete an item from a map use the buit in ```delete``` function:  
```go 
package main

import "fmt" 

func main (){
    m := map[string]string{
        "Picard":"The next generation",
        "Janeway":"Voyager",
        "Kirk":"Star Trek",
    }
    delete(m,"Picard") //delete takes the map as the first arg and the key as the second
    fmt.Println(m)
}
```

[playground](https://play.golang.org/p/SfM4Vipf6S)

### Slices 

A slice in Golang is an abstraction that sits on top of an array. An array normally has a set size 
and must be initialised with a size. A slice is a little more like an array in javascript or an ArrayList in Java.
A slice is a descriptor of an array segment. It consists of a pointer to the array, the length of the segment, and its capacity (the maximum length of the segment). The internals of a slice
are quite clever and mean that working with a slice is as efficient as working with an array. 
[slice internals](https://blog.golang.org/go-slices-usage-and-internals).
It is unusual in Golang to use an array and for our needs it is unlikely you will need one, so we will focus on slices.

As with maps a slice is not naturally thread safe. As with maps we will go over these things when we get to concurrency.

As with a map there are two ways of making a slice. 

1) Slice literals 

```go 

package main 

import "fmt"

func main (){
    captains := []string{"Picard","Janeway","Kirk"} //the type the slice will hold must be specified
    fmt.Println(captains)
}

``` 

[playground](https://play.golang.org/p/PdPD5t0Kb0)

2) As with maps you can use the make function:

```go 
package main 

import "fmt"

func main (){
    captains := make([]string,3) //notice the second arg here, this is the initial length and is required in the case of a slice
    captains[0] = "Picard"
    captains[1] = "Janeway"
    captains[2] = "Kirk" 
    fmt.Println(captains)
}

```
[playground](https://play.golang.org/p/hl4ulayBai)

- Appending to a slice

To add new items to a slice we use the built in append function. ```append``` takes the original slice as it's first argument
and any number of follwing args to append to that slice. It returns the new slice to you.

```go 
package main 

import "fmt"

func main (){
    captains := make([]string,0) //notice the second arg here, this is the initial length
    captains =  append(captains,"Picard","Janeway","Kirk") //append takes any number of additions
    fmt.Println(captains)
}

```

- Getting the length of a slice 

This is the same as a map, you make use of the builtin ```len``` function. 

``` go 
package main 


import "fmt"


func main (){
    captains := []string{"Picard","Janeway","Kirk"} //the type the slice will hold must be specified
    length := len(captains)
    fmt.Println(length)
}

```

[playground](https://play.golang.org/p/fD7XVzrcGJ)

- Iterate through a slice 

Once again we will use the range key word. After all underneath every good map is an array.

```go 
package main 

import "fmt"

func main (){
    captains := []string{"Picard","Janeway","Kirk"}
    for index,value := range captains{
        fmt.Printf("index %d value %s",index,value)
    } 
}

``` 
[playground](https://play.golang.org/p/Q_-AE2lizL)


- Access elements in a slice:

Acessing specific elments in a slice is the same as accessing them in an array in javascript.
Golang does give some additional features though as seen below: 

```go 
package main 

import "fmt"

func main (){
    captains := []string{"Picard","Black Beard", "Janeway","Kirk"}
    fmt.Println(captains[0]) //access the 0 index
    fmt.Println(captains[1:3]) //return a slice containing elements from 1-3
    fmt.Println(captains[:2]) //return a slice containing elements from 2 to the first element
    fmt.Println(captains[2:]) //return a slice containing elements from 2 till the end
}

```


- Delete an item from a slice 
To delete an item from a slice we need to get a new slice with the element removed. 
To do this we use the append function to append all the items but the one we don't want.
```go
package main 

import "fmt"

func main (){
    captains := []string{"Picard","Black Beard", "Janeway","Kirk"}
    captains = append(captains[:1],captains[2:]...) //lets take a minute and go over this
    fmt.Println(captains)
}
```
So what is happening here. It looks a little strange. 

We are assigning the returned slice from append to the captains reciever.
Inside append we are appending everything after the [2] index to everything before the [1] first index. So cutting out the value at [1] 
```[:1]``` before ```[2:]``` after the ```...``` is how Golang declares a variadic argument. So can be read as every value including the one at the [2] index.  

[playground](https://play.golang.org/p/eArgeZieHw)



## Add some new endpoints to our api using slices and maps

Ok lets add a new endpoint that returns a map and a new endpoint that returns a slice. I will leave it for homework to implement adding an api that stores new values in a map and another that
stores new values in a slice. 

So lets add a new endpoint by adding the following to our existing main.go 

```go 

func router() http.Handler {
	//http.HandleFunc expects a func that takes a http.ResponseWriter and http.Request
    http.HandleFunc("/api/echo", Echo)
    http.HandleFunc("/api/map", MapMe)
    http.HandleFunc("/api/slice", SliceMe)
    return http.DefaultServeMux //this is a stdlib http.Handler
}

```

Next add the following handlers to main.go 

```go 

func MapMe (rw http.ResponseWriter, req *http.Request){
    var(
        jsonEncoder = json.NewEncoder(rw) 
    )

    myMap := map[string]Message{
        "message":Message{
            Message : "hello map",
            Stamp : time.Now().Unix(),
        }
    }

    if err := jsonEncoder.Encode(myMap); err != nil{
        rw.WriteHeader(http.StatusInternalServerError)
        return
    }
}

func SliceMe (rw http.ResponseWriter, req *http.Request){
    var(
        jsonEncoder = json.NewEncoder(rw) 
    )

    mySlice := []Message{
        Message{
            Message : "hello map",
            Stamp : time.Now().Unix(),
        },
    }

    if err := jsonEncoder.Encode(mySlice); err != nil{
        rw.WriteHeader(http.StatusInternalServerError)
        return
    }
}

```

I will leave it as an exercise to add a new test for each of these endpoints.

## The Cli 

Our cli will live in a package within our api. 

I have created a PR to the repo which we can go through to take a look at the changes:

[first simple pass](https://github.com/maleck13/api/pull/1/commits/4424e579952799c31121b1eef0f0bf467e9ae8c4)

[second more advanced with tests](https://github.com/maleck13/api/pull/1/commits/fbd9ec115100fd758cb8399ff5a596ceaf15ae27)

[full pr](https://github.com/maleck13/api/pull/1)



### Some Extra notes on maps and slices 

## What if I want to add multiple types to a map or slice?

I would generally recommend avoiding this, but sometimes it could be neccessary. To add multiple types to a single map 
you must declare the map to take the interface type. The interface type is comparable to Object in Java. Everyhing is an interface.

```go 
package main 

import "fmt"

func main(){
    m := map[string]interface{}{
        "test":1,
        "test2":"test",
    }
    fmt.Println(m)
    s := []interface{}{
        "test",
        2,
    }
    fmt.Println(s)
}
```

## Where is the sugar? map, reduce, filter etc?

Golang prides itself on being a simple , pragmatic language and tried to avoid sugar that its creators
feel are unnecessary. So lets add some map reduce functions ourselves to see what is involved:


```go 
package main

import "fmt"

type MyList []string // a custom type can be a builtin type if you want it to be 

func (ml MyList) Filter(f func(string) bool) []string {
	na := []string{}
	for _, v := range ml {
		if add := f(v); add {
			na = append(na, v)
		}
	}
	return na
}

func (ml MyList) Reduce(f func(prev, current string, index int) string) []string {
	na := []string{}
	for i, v := range ml {
		if i == 0 {
			na = append(na, f("", v, i))
		} else {
			na = append(na, f(ml[i-1], v, i))
		}
	}
	return na
}

func (ml MyList) Map(f func(val string) *string) []string {
	na := []string{}
	for _, v := range ml {
		mVal := f(v)
		if nil != mVal {
			na = append(na, *mVal)
		}
	}
	return na
}

//Lets try them out 
func main() {
	list := MyList{"test", "test2"}
	fList := list.Filter(func(v string) bool {
		if v == "test" {
			return true
		}
		return false
	})
	fmt.Println(fList)

	rList := list.Reduce(func(prev, current string, index int) string {
		return fmt.Sprintf("%s%s:%d", prev, current, index)
	})
	fmt.Println(rList)

	mList := list.Map(func(val string) *string {
            val = val + " I WAS MAPPED"
			return &val
		return nil
	})
	fmt.Println(mList)
}

```

[playground](https://play.golang.org/p/WfV_OFBa3o)
