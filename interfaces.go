package main

type Collector interface{
	Enabled() bool
	Run(int)
	Name() string
}

type Outputter interface{
	Enabled() bool
	Increment(string, int)
}
