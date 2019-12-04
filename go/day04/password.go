package main

import (
	"fmt"
	"log"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func checkDecreasing(digits []int) bool {
	for i := 0; i < 5; i++ {
		if digits[i] > digits[i+1] {
			return false
		}
	}
	return true
}

func checkDouble(digits []int) bool {
	for i := 0; i < 5; i++ {
		if digits[i] == digits[i+1] {
			return true
		}
	}
	return false
}

func checkDoubleStrict(digits []int) bool {
	matchCount := make(map[int]int)

	for i := 0; i < 5; i++ {
		if digits[i] == digits[i+1] {
			matchCount[digits[i]]++
		}
	}

	for _, v := range matchCount {
		if v == 1 {
			return true
		}
	}
	return false
}

func checkRequirements(current int) bool {
	num := current
	var digits []int
	for num >= 10 {
		digits = append([]int{num % 10}, digits...)
		num = num / 10
	}
	digits = append([]int{num}, digits...)

	return checkDecreasing(digits) && checkDouble(digits)
}

func checkRequirementsStrict(current int) bool {
	num := current
	var digits []int
	for num >= 10 {
		digits = append([]int{num % 10}, digits...)
		num = num / 10
	}
	digits = append([]int{num}, digits...)

	return checkDecreasing(digits) && checkDoubleStrict(digits)
}

func nextCountUp(current int) int {
	new := current + 1
	for !checkRequirements(new) {
		new++
	}
	return new
}

func nextCountUpStrict(current int) int {
	new := current + 1
	for !checkRequirementsStrict(new) {
		new++
	}
	return new
}

func main() {
	// fmt.Println(checkRequirements(111111))
	// fmt.Println(checkRequirements(223450))
	// fmt.Println(checkRequirements(123789))

	start := 193651
	end := 649729

	current := start
	var matches []int
	for current <= end {
		current = nextCountUp(current)
		if current <= end {
			matches = append(matches, current)
		}
	}

	fmt.Println(len(matches))

	// fmt.Println(checkDoubleStrict([]int{1, 1, 2, 2, 3, 3}))
	// fmt.Println(checkDoubleStrict([]int{1, 2, 3, 4, 4, 4}))
	// fmt.Println(checkDoubleStrict([]int{1, 1, 1, 1, 2, 2}))

	currentStrict := start
	var matchesStrict []int
	for currentStrict <= end {
		currentStrict = nextCountUpStrict(currentStrict)
		if currentStrict <= end {
			matchesStrict = append(matchesStrict, currentStrict)
		}
	}

	fmt.Println(len(matchesStrict))

}
