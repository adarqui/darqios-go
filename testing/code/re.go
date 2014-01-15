package main

import (
	"regexp"
	"fmt"
)

func main() {
	str := "active"
	m1 := "Active"
	re, err := regexp.CompilePOSIX(m1)
	match := re.MatchString(str)
	fmt.Println(match,err)
}
