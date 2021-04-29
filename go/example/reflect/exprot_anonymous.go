package main

import (
	"fmt"
	"reflect"
)

type Z struct {
	L int
	k int
}

type z1 struct {
	L int
	k int
}

type T struct {
	A int
	b int
	z1
	Other struct {
		X int
		y int
	}
	Z
	O *T
}

func main() {
	var t T

	v := reflect.ValueOf(t)
	ty := v.Type()

	for i, n := 0, ty.NumField(); i < n; i++ {
		f := ty.Field(i)
		fmt.Printf("%+v\n", f)
		fmt.Println(f.Name, f.PkgPath, f.Anonymous)
		fmt.Println("---")

	}
}

// v := reflect.ValueOf(map[string]string{}) ; v.CanSet()==false
