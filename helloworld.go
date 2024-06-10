package main

import (
    "fmt"
    "sync"
    "time"
)

func printHello(done chan bool, wg *sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("Hello, World!")
    time.Sleep(time.Second)
    done <- true
}

func main() {
    var wg sync.WaitGroup
    done := make(chan bool, 1)

    wg.Add(1)
    go printHello(done, &wg)
    
    <-done
    wg.Wait()
    fmt.Println("Program finished executing.")
}
