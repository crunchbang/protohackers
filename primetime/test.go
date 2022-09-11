package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	x, _ := strconv.ParseFloat("57027812634041659986718814375065495113012574076530602544491630", 64)
	fmt.Println(math.Floor(x))
	fmt.Println(math.Ceil(x))
	fmt.Println(math.Ceil(x) == math.Floor(x))
	// fmt.Println(strconv.ParseUint("57027812634041659986718814375065495113012574076530602544491630", 10, 64))
	// fmt.Println(strconv.ParseInt("57027812634041659986718814375065495113012574076530602544491630", 10, 64))
	// "47000137448339464127393145078746062236088188569847018608221417"
}