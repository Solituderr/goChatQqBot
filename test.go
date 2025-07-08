package main

import (
	"fmt"
	"regexp"
)

func Test() {
	re := regexp.MustCompile("\\.r(\\d*)d(\\d*)")
	match := re.FindStringSubmatch(".r1d100 1241321r")
	match1 := re.FindStringSubmatch(".rd100")
	match2 := re.FindStringSubmatch(".r1d")
	fmt.Println(match, match1, match2)

}
