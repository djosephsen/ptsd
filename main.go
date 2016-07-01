package main

import (
	"os"
)

func main(){
	if token := os.Getenv(`PDTOKEN`); token != ``{
		runPagerDuty(token)
	}
}
