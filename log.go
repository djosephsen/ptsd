package main

import (
	"fmt"
	"os"
)

func debug(msg string){
	if os.Getenv("DEBUG") != `` {
		fmt.Printf("DEBUG :: %s\n",msg)
	}
}
