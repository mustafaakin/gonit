package main

import (
	"os"
)

func main() {
	println("I will arrange terminals and stuff")
	println("Getting an environment variable via init", os.Getenv("MYVAR"))
}
