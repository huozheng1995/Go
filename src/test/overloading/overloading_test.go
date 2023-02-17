package overloading

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

func isInt(i interface{}) bool {
	switch i.(type) {
	case int:
		return true
	default:
		return false
	}
}

func sumAndPrint(vals ...interface{}) error {
	if len(vals) > 3 {
		return errors.New("Too many parameters")
	}

	nonNumber := false
	sum := 0
	concat := ""

	for _, val := range vals {
		if isInt(val) == false {
			nonNumber = true
			concat += val.(string) + " "
		} else {
			concat += strconv.Itoa(val.(int)) + " "
			sum += val.(int)
		}
	}

	if nonNumber {
		fmt.Println(concat)
	} else {
		fmt.Println(sum)
	}
	return nil
}

func Test(t *testing.T) {
	sumAndPrint()                 //Prints 0
	sumAndPrint(1)                //Prints 1
	sumAndPrint(1, 2)             //Prints 3
	sumAndPrint(1, 2, 3)          //Prints 6
	sumAndPrint("Hi")             //Prints Hi
	sumAndPrint("Hello", "World") //Prints Hello World
	sumAndPrint("Hi", 1, 2)       //Prints Hi 1 2

	if err := sumAndPrint(1, 2, 3, 4); err != nil {
		fmt.Println(err) //Prints Too many parameters
	}
}
