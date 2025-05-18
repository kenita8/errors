package main

import (
	"fmt"

	"github.com/kenita8/errors"
)

var Cause1 = errors.New("cause1")
var Cause2 = errors.New("cause2")
var Cause3 = errors.New("cause3")
var Cause4 = errors.New("cause4")

func main() {
	err1 := Cause1.WithDetails("key1", "val1")
	fmt.Printf("%s\n", err1) // cause1(key1=val1)

	err2 := Cause2.WithDetails("key2", "val2").Wrap(err1)
	fmt.Printf("%s\n", err2) // cause2(key2=val2): cause1(key1=val1)

	err3 := Cause3.WithDetails("key3", "val3").Wrap(err2)
	fmt.Printf("%s\n", err3) // cause3(key3=val3): cause2(key2=val2): cause1(key1=val1)

	err4 := Cause4.WithDetails("key4", "val4").Wrap(err3)
	fmt.Printf("%s\n", err4) // cause4(key4=val4): cause3(key3=val3): cause2(key2=val2): cause1(key1=val1)
}
