package main

import (
	"encoding/json"
	"fmt"
)

type A struct {
	Uno int
}

type B struct {
	Uno int
	Due int
}

func main() {
	b := B{
		Uno: 1,
		Due: 2,
	}
	if jb, err := json.Marshal(b); err == nil {
		var a A
		if err = json.Unmarshal(jb, &a); err == nil {
			fmt.Print(a)
		}
	}
}
