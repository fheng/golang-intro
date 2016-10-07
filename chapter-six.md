## Goals of this lesson

- Recap on go routines
- Expand on channels and their usages
- Discuss locks / mutexs and their uses
- Discuss other useful concurrency tools such as ``` sync.WaitGroup ```

## Goroutines
So just to recap, A goroutine is a light weight thread or "actor" that multiplexes ontop of a single "real" os thread. 
This model allows Go to create many thousands of goroutines that span on a small number of actual os threads.

The ```go ``` keyword is used to create a goroutine.

```go 

func main(){
    go func(){
        fmt.Println("I will execute async")    
    }()
}

```
So our func will execute within another goroutine. 


## Channels

If you want to communicate something from one executing goroutine to another, you should use a channel. 

### Creating a channel
To make a new channel you use the builtin make command just like a map
```go 
 myStrChan := make(chan string) //this type can be any valid type
 myIntChan := make(chan int)
 //you can also create a buffered channel we will expand on these shortly.
 myBChan := make(chan int, 5) //5 here is how many messages to buffer before begining to block.
```

### Closing a channel  
To close a channel or turn it off so that nothing new can be sent through it you use the ```close``` builtin function.

```go
myStrChan := make(chan string) //this type can be any valid type
close(myStrChan)
```

### Some behaviour to understand

- sending on a channel when nothing is reading from it will block until something reads. Remember channels are a way for routines to synchronise and share data.

[playground](https://play.golang.org/p/npQZWK_AQU)

- sending on a closed channel will panic 

[playground](https://play.golang.org/p/f7ABM6_FnZ)

- reading from a closed channel will result in the default value for that type 

[playground](https://play.golang.org/p/FFDs1xDDWT)

### Using channels 

```go 
package main 

import "fmt"

func worker(done chan string){
    done <- "done" //send the string done through the chanel to any other routine that is reading from that channel
}

func main(){
    //make a new channel 
    d := make(chan string)
    go worker(d)
    done := <- d //this will block until something is sent.
    fmt.Println(done) 
    close(d) //

}
```
[playground](https://play.golang.org/p/_VMlxfl59q)

You can range over channels and recieve as many values as are sent until the channel is closed.

```go 
package main 

import "fmt"
import "time"

func worker(done chan string){
    done <- "done" //send the string done through the chanel to any other routine that is reading from that channel
}

func main(){
    //make a new channel 
    d := make(chan string)
    go worker(d)
    go worker(d)
        //close in background 
    go func (){
      time.Sleep(1 * time.Second)
      close(d)
    }()

    for done := range d { // each time something is sent through d it will be assigned to done.
     fmt.Println(done) 
    }
}
```
[playground](https://play.golang.org/p/6-qIy8mlE1)

## Deadlock!!

In all concurrent languages you have the concept of a Deadlock. This essentially means there is no way for the involved actors / goroutines to proceed. This causes the program to panic and exit.

```go

package main 

import "fmt"
import "time"

func worker(done chan string){
    done <- "done" //send the string done through the chanel to any other routine that is reading from that channel
}

func main(){
    //make a new channel 
    d := make(chan string)
    go worker(d)
    go worker(d)
        //close in background 
    go func (){
      time.Sleep(1 * time.Second)
      // WE REMOVED THE CLOSE SO NOW THE GO ROUTINE WILL BECOME DEADLOCKED
    }()

    for done := range d { // each time something is sent through d it will be assigned to done.
     fmt.Println(done) 
    }
}
```
[playground](https://play.golang.org/p/5U1yInSRWg)

To help detect this issues you need to ensure you write good tests. There are also some builtin helpers. ```go vet``` will catch some things.
Also running ```go test -race``` will turn on the race detector when running your tests.
[more on races and detection](https://golang.org/doc/articles/race_detector.html)

## channels allow you to build powerful concurrent progams

These examples are taken from the [go blog](https://blog.golang.org/pipelines)

**Pipline**

```go 
package main 
// function takes any number of ints and returns a channel through which it will send the generated numbers
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}
// sq takes a channel that sends ints and returns a channel that also sends out ints
func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

func main() {
    // Set up the pipeline.
    c := gen(2, 3)
    out := sq(c)

    // Consume the output.
    fmt.Println(<-out) // 4
    fmt.Println(<-out) // 9
}

```
[playground](https://play.golang.org/p/6dAn5s2-KA)

I am concious of not overloading at this stage. I recommend taking a look at the blog post once you get a little more familiar with the main concepts.

## Buffered Channels

A buffered channel allows you to send a number of messages before it blocks.

As mentioned before sending on a channel that has nothing reading from it will block. Sometimes it may be desirable to allow the sender to queue up messages so that it can do other things.
For example, if you had something that generated a lot of work but the workers took a long time to do their work. Without a buffered channel it would look like this:

```go 
package main

import "fmt"
import "time"

func slowWorker(incoming chan string) {
	for in := range incoming {
		fmt.Println("having a sleep first")
		time.Sleep(1 * time.Second)
		fmt.Println("done ", in)
	}
}

func main() {
	work := make(chan string)
	go slowWorker(work)
	work <- fmt.Sprintf("work %d", 1)
	fmt.Println("sending some more work")
	work <- fmt.Sprintf("work %d", 2) // will block here until the slowWorker is done
	fmt.Println("finished sending work")
	close(work)
	fmt.Println("moving on to do something else")
}
```

[playground](https://play.golang.org/p/M4VFed6ejE)

So we can see that no new work can be sent on this channel until the recieving channel is ready to recieve again, this means that the sending routine is blocked until the recieving routine is ready. 
What if you don't want this? Well we can use a buffered channel. When declaring the channel above if we set a 3rd param: 

```
make(chan string 2)

```
Adding the buffer of 2 stopped the sender from being blocked and to carry on with something else until the buffer is full.
We can see this working here:
[playground](https://play.golang.org/p/NlGYWoxtkj)


## Select

What if you want to recieve messages from more than one channel at a time? To do this we use the select statement.

From the go welcom tour. 

> The select statement lets a goroutine wait on multiple communication operations. A select blocks until one of <it></it>s cases can run, then it executes that case. It chooses one at random if multiple are ready.

```go 
package main

import "fmt"

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for { // using a for loop like this is the same as a while loop 
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
```
[playground](https://play.golang.org/p/uf94rfXkva)

So in our above example we are using a for loop to constantly check the select statement. Once any of this evaluate to true the code will execute that is associated with the corrosponding case. 

Another example using a time out.

## Locks and WaitGroups

Go also offers mutexs to allow for more traditional synchronisation. 
Channels increase the complexity of the code and require a detailed understanding of how where those channels are being used. A mutex offers a more simple abstraction for more simple use cases.
A common way to see a mutex used is to protect concurrent access to something like a map.

```go 
package main

import (
	"fmt"
	"sync"
)

type concurrentMap struct {
	sync.Mutex //gives this object the behaviour of a mutex
	Data       map[string]string
}

func (cm *concurrentMap) Get(key string) string {
	cm.Lock()
	defer cm.Unlock()
	return cm.Data[key]

}
func (cm *concurrentMap) Put(key, val string) {
	cm.Lock()
	cm.Data[key] = val
	cm.Unlock()
}

// ConcurrentAccess defines access to concurrentMap
type ConcurrentAccess interface {
	Put(key, val string)
	Get(key string) string
}

// NewConcurrentMap return a ConcurrentAccess Map
func NewConcurrentMap() ConcurrentAccess {
	return &concurrentMap{
		Data: make(map[string]string),
	}
}

func main() {
	var cMap = NewConcurrentMap()
	wg := sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(index int) {
			cMap.Put("test", fmt.Sprintf("test %d", index))
		}(i)
		go func() {
			wg.Done()
			fmt.Println("get", cMap.Get("test"))
		}()
	}

	wg.Wait()

}
``` 

So there two main things worth pointing out here. The first is the use of sync.Mutex. So we can see I have created a new type that has a map and a mutex.
The mutex is embeded in the struct giving it all the methods of a mutex. We have also defined an interface specifying how a the conncurrentAccess should be used.
Now within the Put and Get methods we lock and unlock the mutex allowing us to safe gaurd access to the data in the map from concurrent go routines.


## WaitGroups

A wait group is a simple mechanism provided by the sync package to allow a go routine to wait until another set of actions have completed. This is particularly useful if 
you don't need any data back from that set of go routines / or they put there data somewhere as part of what they are doing. We can see a WaitGroup being used in the sample code
above.

I would encourage you to play with the concurrency tools provided by Go, but as with all conccurrency, it increases complexity. When faced with a problem we 
should not first reach for conccurrency, but rather exhaust other options and then reach for conccurrency.

## Wrap Up 

So this is the last lesson in the golang intro. I wanted to just briefly recap.   