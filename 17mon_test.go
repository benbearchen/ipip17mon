package ipip17mon

import "testing"

import (
	"fmt"
	"strings"
)

func Test(t *testing.T) {
	x := new(IP17monDatax)
	err := x.Init("./17monipdb.datx")
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}

	ips := []string{
		"116.228.111.18",
		"8.8.8.8",
		"119.131.76.177",
		"100.64.36.232",
	}

	for _, ip := range ips {
		loc, err := x.Find(ip)
		if err != nil {
			fmt.Println("err:", err)
			break
		}

		fmt.Println(loc)
		fmt.Printf("%q\n", strings.Split(loc, "\t"))
	}
}
