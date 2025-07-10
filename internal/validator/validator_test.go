package validator

import (
	"log"
	"testing"
)

type test struct {
	Id     string `valid:"required,min=6"`
	Length int    `valid:"min=10,max=12"`
	Email  string `valid:"email"`
}

func TestSuccess1(t *testing.T) {
	t.Parallel()
	test1 := test{
		Id:     "tomlovevim",
		Length: 11,
		Email:  "tomlovevim@gmail.com",
	}
	va := NewValidator()
	oerr := va.Check(test1)
	if oerr != nil {
		log.Printf("-- %v", oerr.Error())
	} else {
		log.Printf("-- test1 passed")
	}
}

func TestSuccess2(t *testing.T) {
	t.Parallel()
	test2 := test{
		Id:     "tomlov",
		Length: 10,
		Email:  "tomlovevim@163.com",
	}
	va := NewValidator()
	oerr := va.Check(test2)
	if oerr != nil {
		log.Printf("-- %v", oerr.Error())
	} else {
		log.Printf("-- test2 passed")
	}
}

func TestFail1(t *testing.T) {
	t.Parallel()
	test1 := test{
		Id:     "tomlo",
		Length: 9,
		Email:  "tomlovevim@163.com",
	}
	va := NewValidator()
	oerr := va.Check(test1)
	if oerr != nil {
		log.Printf("-- %v", oerr.Error())
	} else {
		log.Printf("-- test1 passed")
	}
}

func TestFail2(t *testing.T) {
	t.Parallel()
	test2 := test{
		Id:     "tomlo",
		Length: 9,
		Email:  "tomlovevim@163.com",
	}
	va := NewValidator()
	oerr := va.Check(test2)
	if oerr != nil {
		log.Printf("-- %v", oerr.Error())
	} else {
		log.Printf("-- test2 passed")
	}
}

func TestFail3(t *testing.T) {
	t.Parallel()
	test3 := test{
		Id:     "tomlov",
		Length: 10,
		Email:  "@163.com",
	}
	va := NewValidator()
	oerr := va.Check(test3)
	if oerr != nil {
		log.Printf("-- %v", oerr.Error())
	} else {
		log.Printf("-- test3 passed")
	}
}
