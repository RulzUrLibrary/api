package utils

import (
	"strconv"
)

func IsIsbn10(isbn string) bool {
	if len(isbn) != 10 {
		return false
	}

	var sum int
	var multiply int = 10
	for i, v := range isbn {
		digitString := string(v)

		if i == 9 && digitString == "X" {
			digitString = "10"
		}

		digit, err := strconv.Atoi(digitString)
		if err != nil {
			panic(err)
		} else {
			sum = sum + (multiply * digit)
			multiply--
		}
	}

	return sum%11 == 0
}

func IsIsbn13(isbn string) bool {

	if len(isbn) != 13 {
		return false
	}

	var sum int
	for i, v := range isbn {
		var multiply int
		if i%2 == 0 {
			multiply = 1
		} else {
			multiply = 3
		}

		digit, err := strconv.Atoi(string(v))
		if err != nil {
			panic(err)
		} else {
			sum = sum + (multiply * digit)
		}
	}

	return sum%10 == 0
}
