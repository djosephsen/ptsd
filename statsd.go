package main

import (
	"fmt"
	statsd "github.com/etsy/statsd/examples/go"
)

var metrics = statsd.New(`localhost`, 8125)

func increment(key string, value int){
	fmt.Printf("%s:: +%d\n",key,value)
	metrics.IncrementByValue(key, value)
}
