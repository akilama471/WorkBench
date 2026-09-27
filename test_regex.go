package main

import (
	"fmt"
	"regexp"
)

func main() {
	fmt.Println(regexp.MustCompile(`\b(\d+\.\d+(?:\.\d+)?)\b`).FindString("php-8.4.26-Win32-vs17-x64.zip"))
}
