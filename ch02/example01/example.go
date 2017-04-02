package main

import "fmt"

type person struct {
	Age  int
	Name string
}

func (p person) greeting() string {
	return fmt.Sprintf("hey! I'm %s", p.Name)
}

func (p person) name() string {
	return p.Name
}

type dog struct {
	Age   int
	Name  string
	Owner string
}

func (d dog) greeting() string {
	return "wan!"
}

func (d dog) name() string {
	return d.Name
}

type nameTyp string

func (s nameTyp) name() string {
	return string(s)
}

func (s nameTyp) greeting() string {
	return fmt.Sprintf("nameType: %s", s)
}

type animal interface {
	greeting() string
	name() string
}

func main() {
	p1 := person{
		Name: "achiku",
		Age:  31,
	}
	p2 := person{
		Name: "moqada",
		Age:  31,
	}
	d1 := dog{
		Name:  "taro",
		Age:   2,
		Owner: "moqada",
	}
	n1 := nameTyp("8maki")
	for _, a := range []animal{p1, p2, d1, n1} {
		fmt.Printf("%s: %s\n", a.name(), a.greeting())
	}
}
