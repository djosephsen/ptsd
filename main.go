package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var COLLECTORS = []Collector{}
var OUTPUTTERS = []Outputter{}

func main() {
	var offset int
	if interval := os.Getenv(`PTSD_INTERVAL`); interval == `` {
		offset = 60
	} else {
		offset, _ = strconv.Atoi(interval)
	}
	debug(fmt.Sprintf("Configured offset: %d", offset))

	if enabledCollectors() == 0{
		fmt.Println("No collectors enabled!")
		os.Exit(2)
	}
	if enabledOutputters() == 0{
		fmt.Println("No outputters enabled!")
		os.Exit(2)
	}

	ticker := time.NewTicker(time.Duration(offset) * time.Minute)
	for {
		for _, collector := range COLLECTORS {
			if collector.Enabled() {
				debug(fmt.Sprintf("collector enabled: %s", collector.Name()))
				collector.Run(offset)
			}
		}
		<-ticker.C
	}
}

func increment(key string, value int) {
	for _, outputter := range OUTPUTTERS {
		if outputter.Enabled() {
			outputter.Increment(key, value)
		}
	}
}

func enabledCollectors() int{
	c := 0
	for _, collector := range COLLECTORS {
		if collector.Enabled(){
			c += 1
		}
	}
	return c
}

func enabledOutputters() int{
	c := 0
	for _, outputter := range OUTPUTTERS {
		if outputter.Enabled(){
			c += 1
		}
	}
	return c
}
