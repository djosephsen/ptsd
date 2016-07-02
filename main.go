package main

import (
	"os"
	"fmt"
	"strconv"
)

var COLLECTORS = []Collector{}
var OUTPUTTERS = []Outputter{}

func main(){
	var offset int
	if interval := os.Getenv(`PTSD_INTERVAL`); interval == ``{
		offset = 60	
	}else{
		offset,_ = strconv.Atoi(interval)
	}
	debug(fmt.Sprintf("Configured offset: %d",offset))

	for _,collector := range COLLECTORS{
		if collector.Enabled(){
			debug(fmt.Sprintf("collector enabled: %s",collector.Name()))
			collector.Run(offset)
		}
	}
}

func increment(key string, value int){
	for _,outputter := range OUTPUTTERS {
		if outputter.Enabled(){
			outputter.Increment(key,value)
		}
	}
}
