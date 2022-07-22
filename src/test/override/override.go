package main

import "fmt"

type Human struct {
	name  string
	age   int
	phone string
}
type Student struct {
	Human
	school string
}
type Employer struct {
	Human
	company string
}

//implement Human method
func (h *Human) SetName(name string) {
	fmt.Print("human")
	h.name = name
}
func (h *Human) SetAge(age int) {
	h.age = age
}
func (h *Human) SetPhone(phone string) {
	h.phone = phone
}
func (h *Human) GetInfo() Human {
	return *h
}

func (s *Student) SetName(name string) {
	fmt.Print("student")
	/*about here we can use two wanys to change the value ,so ,how different there ?????*/
	s.name = name
	//s.Human.name = name
}

func main() {
	s := Student{}
	s.SetPhone("18755201184")
	s.SetName("tsong")
	s.SetAge(26)
	fmt.Print(s.GetInfo())
}
