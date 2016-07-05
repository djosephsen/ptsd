package main

import (
	"fmt"
	"os"
)

type TextOutputter struct{}

var TXT = &TextOutputter{}

func init() {
	OUTPUTTERS = append(OUTPUTTERS, TXT)
}

func (this *TextOutputter) Name() string {
	return `Text Outputter`
}

func (this *TextOutputter) Increment(key string, value int) {
	fmt.Printf("%s :: +%d\n", key, value)
}

func (this *TextOutputter) Enabled() bool {
	if os.Getenv(`PTSD_TXT`) != `` {
		return true
	} else {
		return false
	}
}
