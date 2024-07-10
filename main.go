package main

import (
	"fmt"
	"github.com/kritibb/ktz/cmd"
	"os"
    "strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a city name!")
		return
	}
    city := os.Args[1:]
    cityString:=strings.Join(city, " ")
    tz, err:=cmd.GetCityTZ(cityString)
    if err!=nil{
        fmt.Println("Error:", err)
        return
    }
    cmd.DisplayTime(tz)


}
