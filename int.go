package main

import "fmt"

func main() {
	var o = 0666
	fmt.Printf("%d %[1]o %#[2]o\n",o,0777)
}