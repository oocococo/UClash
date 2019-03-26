package main

import "fmt"

func Checkerr(err error) {
	if err != nil {
		fmt.Print(err)
	}
}
