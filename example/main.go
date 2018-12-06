package main

import (
	"errors"
	"fmt"

	"github.com/MagalixTechnologies/common/alltogether"
)

func main() {
	arr := []string{"uno", "due", "tre", "quattro", "cinque", "sei", "sette", "otto", "nove", "dieci"}

	pr, err := alltogether.NewConcurrentProcessor(arr, func(task string) error {
		fmt.Println(task)
		if task == "due" {
			return errors.New("some error for due")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	errs := pr.Do()
	fmt.Println(errs)
	fmt.Println("allNil", errs.AllNil())

	arrOrStructs := []*SomeType{
		&SomeType{Value: "uno"},
		&SomeType{Value: "due"},
		&SomeType{Value: "tre"},
		&SomeType{Value: "quattro"},
		&SomeType{Value: "cinque"},
		&SomeType{Value: "sei"},
		&SomeType{Value: "sette"},
		&SomeType{Value: "otto"},
		&SomeType{Value: "nove"},
		&SomeType{Value: "dieci"},
	}

	pr, err = alltogether.NewConcurrentProcessor(arrOrStructs, func(task *SomeType) error {
		fmt.Println(task)
		if task.Value == "due" {
			return errors.New("some error for due")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	errs = pr.Do()
	fmt.Println(errs)
	fmt.Println("allNil", errs.AllNil())
}

type SomeType struct {
	Value string
}
