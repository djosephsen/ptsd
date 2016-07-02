package main

import (
	"os"
	"strconv"
	statsd "github.com/etsy/statsd/examples/go"
)

type SDOutputter struct{
	Server	string
	Port		int
}
var SD = &SDOutputter{}

func init(){
	OUTPUTTERS = append(OUTPUTTERS,SD)
}

func (this *SDOutputter) Name() string{
	return `Statsd Outputter`
}

func (this *SDOutputter) Increment(key string, value int){
	metrics := statsd.New(this.Server, this.Port)
	metrics.IncrementByValue(key, value)
}

func (this *SDOutputter) Enabled() bool{
	if server := os.Getenv("SD_SERVER"); server == ``{
		return false
	}else{
		this.Server = server
	}
	if port := os.Getenv("SD_PORT"); port == ``{
		return false
	}else{
		this.Port,_ = strconv.Atoi(port)
	}
	return true
}
