# roundrobin
 A simple roundobin algorithm implementation written in Go.
 
 This can be used to balance addresses on layer 4.
 
 You can create a new service from addresses (ip:port). Addresses will be validated.

 ## Example

 ```go
package main

import "github.com/mdnix/roundrobin"

func main() {
    service, err := roundrobin.NewService([]string{
    	"192.168.1.1:4444",
    	"192.168.1.2:4444",
    	"192.168.1.3:4444",
    })
    if err != nil {
		// handle error
	}

    first := service.Next()
    second := service.Next()
    ...
}
 ```
