package main

import "fmt"

type A struct {
	AField int
}

type B struct {
	AField int
}

func (b B) Hello() {
	fmt.Printf("Hello! I'm %d\n", b.AField)
}

type SayHello interface {
	Hello()
}

type SayCiao SayHello

func niceToMeetYou(b B) {
	fmt.Print("Hi B!")
	b.Hello()
}

func main() {
	a := A{
		AField: 3,
	}

	var b interface{}
	b = B(a)
	if c, ok := b.(A); ok {
		d := A(c)
		fmt.Print(d.AField)
	} else if c, ok := b.(B); ok {
		fmt.Print(c)
	}

	//niceToMeetYou(a)
}
