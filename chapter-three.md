## Goals of this lesson
- Expolre structs and their uses in Golang.
- Look at the two main collection types slices and maps
- Start working on a cli to interact with our new server

## Note this lesson builds on top of the code from lesson 2

We started a basic api web server in lesson 2 the code can be found here:
https://github.com/maleck13/api

- set up for this lesson 
```bash

cd $GOPATH/src/github.com/<YOUR_USER>
git clone git@github.com:maleck13/api.git
cd api 
go run main.go 

```
Feel free to fork the repo. I am also happy to review any changes you make.

We are going to learn about collection types and use them to add two new endpoints to our server.

## Collections types 

Golang has two main collection types.
1) maps.
2) slices.

### Maps 

Maps are somewhat similar to what other languages call “dictionaries” or “hashes”.

maps can created in two ways: 

1) using map literals:

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

2) With the make keyword:
The make built-in function allocates and initializes an object of type slice, map, or chan (only).

```go 

package main 

import "fmt"

func main (){
    m := make(map[string]string) //again we are declaring the type of the key and the type of the value
    m["Picard"] = "The next generation"
    m["Janeway"] = "Voyager"
    m["Kirk"] = "Star Trek"
    fmt.Println(m)
}
```

[playground](https://play.golang.org/p/Gz9OfA5rFR)

Maps can hold any valid type or pointer.

- checking if a key exists in a map 


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

    if v,ok := m["Picard"]; ok{
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
    m := map[string]string{ //here we are declaring the key to be of type string and the value to be of type string
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
A range clause provides a way to iterate over an array, slice, string, map, or channel.
It is used in conjunction with a for loop. One thing to note is that when you use the range key work for each
loop it copies the current key and value into the vars so if you change the value in the loop it will not be reflected in the map or slice.

- Get the size of a map using the built in len function 

```go 

package main

import "fmt" 

func main (){
    m := map[string]string{ //here we are declaring the key to be of type string and the value to be of type string
        "Picard":"The next generation",
        "Janeway":"Voyager",
        "Kirk":"Star Trek",
    }
    fmt.Println(len(m)) //outputs 3
}
```
[playground](https://play.golang.org/p/2rYP28LdnZ)

- To Delete an item from a map use the buit in delete function  
```go 

package main

import "fmt" 

func main (){
    m := map[string]string{ //here we are declaring the key to be of type string and the value to be of type string
        "Picard":"The next generation",
        "Janeway":"Voyager",
        "Kirk":"Star Trek",
    }
    delete(m,"Picard")
    fmt.Println(m)
}
```

[playground](https://play.golang.org/p/SfM4Vipf6S)

### Slices 

A slice in Golang is an abstraction that sits on top of an array. An array normally has a set size 
and must be initialised with a size. A slice is a little more like an array in javascript or an ArrayList in Java.
A slice is a descriptor of an array segment. It consists of a pointer to the array, the length of the segment, and its capacity (the maximum length of the segment). The internals of a slice
are quite clever and mean that working with a slice is as efficient as working with an array. 
[slice internals] (https://blog.golang.org/go-slices-usage-and-internals).
It is unusual in Golang to use an array and for our needs it is unlikely you will need one, so we will focus on slices.

Again there are two ways of making a slice. 

1) slice literals 

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
    captains := make([]string,3) //notice the second arg here, this is the initial length
    captains[0] = "Picard"
    captains[1] = "Janeway"
    captains[2] = "Kirk" 
    fmt.Println(captains)
}

```
[playground](https://play.golang.org/p/hl4ulayBai)

- Appending to a slice
To add new items to a slice we use the built in append function. 

```go 
package main 

import "fmt"

func main (){
    captains := make([]string,0) //notice the second arg here, this is the initial length
    captains =  append(captains,"Picard","Janeway","Kirk") //append takes any number of additions
    fmt.Println(captains)
}

```

- Iterate through a slice 

Once again we will use the range key word. After all underneath a map is an array.

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


- Access elements in a slice 

```go 
package main 

import "fmt"

func main (){
    captains := []string{"Picard","Black Beard", "Janeway","Kirk"}
    fmt.Println(captains[0])
    fmt.Println(captains[0:2])
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