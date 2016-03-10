package main

import (
	"fmt"
)

type HandleFunc func(*Mux)

type Mux struct {
	handler []HandleFunc
	index   int
	written bool
}

func (m *Mux) next() {
	m.index += 1
	m.run()
}

func (m *Mux) run() {
	for m.index < len(m.handler) {
		m.handler[m.index](m)
		m.index += 1

		if m.written {
			return
		}
	}
}

func main() {
	mux := new(Mux)

	handerls := make([]HandleFunc, 4)
	handerls[0] = func(m *Mux) {
		fmt.Println("A0")
		m.next()
		fmt.Println("A1")
	}
	handerls[1] = func(m *Mux) {
		fmt.Println("B0")
		m.next()
		fmt.Println("B1")
	}
	handerls[2] = func(m *Mux) {
		fmt.Println("C0")
		fmt.Println("C1")
		//m.written = true //类似下面的goto
	}
	handerls[3] = func(m *Mux) {
		fmt.Println("D0")
		fmt.Println("D1")
	}
	mux.handler = handerls

	mux.run()

	//---等价于

	func() {
		fmt.Println("A0")
		{
			fmt.Println("B0")
			{
				fmt.Println("C0")
				fmt.Println("C1")
				//goto LABEL1
				fmt.Println("D0")
				fmt.Println("D1")
			}
			//LABEL1:
			fmt.Println("B1")
		}
		fmt.Println("A1")
	}()
}
