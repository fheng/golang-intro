## Goals of this lesson

- Continue to explore structs and how you can compose them together, similar to how we saw for interfaces.

- Explore how Go deals with dependencies (its short comings and the tools that resolve these issues)

- Introduce some of the constructs of concurrency in Go.

Structs can "inherit" from other structs by a method known as embeding. Here is a simple example:

```go 
package main

import (
	"fmt"
)

type a struct {
	Name string
}
//embeds a value of type a
type b struct {
	a
}
//embeds a pointer to an a
type c struct {
	*a
}

func main() {
	a := a{Name: "Janeway"}
	fmt.Println(a.Name)
	b := &b{a: a}
	fmt.Println(b.Name)
	c := &c{a: &a}
	fmt.Println(c.Name)
}
```
[playground](https://play.golang.org/p/6AataGi66F)

If the struct that is embeded has methods then the struct into which it is embedded will also have those methods

```go 
package main

import (
	"fmt"
)

type a struct {
	Name string
}

func (as a) NameLength() int {
	return len(as.Name)
}

type b struct {
	a
}

type c struct {
	*a
}

func main() {
	a := a{Name: "Janeway"}
	fmt.Println(a.Name)
	fmt.Println(a.NameLength())
	b := &b{a: a}
	fmt.Println(b.Name)
	fmt.Println(b.NameLength())
	c := &c{a: &a}
	fmt.Println(c.Name)
	fmt.Println(c.NameLength())
}

``` 

So what happens if b also has a method called NameLength?

```go 
package main

import (
	"fmt"
)

type a struct{
	Name string
}

func (as a)NameLength() int{
    return len(as.Name)
}

type b struct{
 a
}

func (bs b)NameLength() int{
    return len(bs.Name) -1
}

type c struct{
*a
}

func main() {
	a := a{Name:"Janeway"}
	fmt.Println(a.Name)
	fmt.Println(a.NameLength())
	b := &b{a:a}
	fmt.Println(b.Name)
	fmt.Println(b.NameLength()) //this is the method on b itself
	fmt.Println(b.a.NameLength()) //notice we can still reference the behaviour of a
	c := &c{a:&a}
	fmt.Println(c.Name)
	fmt.Println(c.NameLength())
}
//will output
//Janeway
//7
//Janeway
//6
//7
//Janeway
//7
```

[playground](https://play.golang.org/p/-NFr1QSmp6)

So it acts as I think we would expect it to. If there is a method on the referenced struct that matches it will chose that implementation, however if you want to call the method on the embedded struct you can do so by going up the chain of properties

```
b.a.MethodName()
```

So what's the difference between this and classical inheritance?

1) The embedded struct knows nothing about the struct in which it has been embedded. There is no way to go down the chain you cannot do ``` a.b.MethodName ```
2) There is no way to have the concept of ```abstract``` methods that the embedding struct is forced to implement.

Again the rules around pointers apply. If you modify an embeded pointer you will also modify the value of that pointer for anything else using the same pointer address.

## Dependencies 

So far we haven't needed to use a single dependency. But it is naive to think you will never need one.
Although much of the tooling around Go is excellent. The Go team took what is considered a fairly crude approach to dependency management. We have to remember that initally Go was designed to suit Google's use case. It was not really a concern of theirs to think of lots of other use cases.
Google is well known to have an immense mono repo so ``` go get ``` works well for them. This is not the case with everyone else. By default go get pulls the master branch of the repo you point it at. When you do a go get it pulls in the required dependencies, this means there are issues with
reproducibility. 
As of go 1.5 they looked to address some of the issues by introducing the vendor directory. If a directory called vendor exists in the current package it will first attempt to resolve dependencies their. Thing node_modules, except you are encouraged to check in your dependencies.
As go is compiled to machine code, there is no issues with checking in dependencies, unlike node modules that might have native dependencies.
Once vendoring was in place several project sprang up to manage these dependencies and updating these dependencies. There are several but they all work in similar ways:

Create a file that registers the dependencies and their current commit hash. When you update a dependency you use the tool and it will update the commit. This makes dependency management much nicer. There is still a way to go and 
due to the popularity of Go, there is currently a community panel set up to better address dependency.

I like a tool called glide:
[glide tool](https://github.com/Masterminds/glide)

Here is an example of the files glide creates in your project.

[lock file](https://github.com/feedhenry/rhm/blob/master/glide.lock)

[glide.yaml think package.json](https://github.com/feedhenry/rhm/blob/master/glide.yaml)

[how glide works](https://github.com/Masterminds/glide#how-it-works)

Personally I think vendoring dependencies is a good thing overall as long as the versions are managed with a tool like glide. I think it forces you to see the changes happening in your dependency when you update it as it will show up
in your PR 

With glide you would install a dependency as follows. 

- Having done a ```glide init ``` you would install the dependency by doing a 
``` glide get <some package> ``` when you do this glide will look for release tags on the repo and prompt you for the version you want to use. Once this is done it will update the lock and glide file and add the dependency to your vendor directory (creating one if it doesn't exist).

## Touching on concurrency before next week.