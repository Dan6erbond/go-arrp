package main

import "openwt.com/go-arrp/internal"

func main() {
	r := internal.MakeServer()
	r.Run()
}
