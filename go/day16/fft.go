package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// input := "12345678"
	// input := "80871224585914546619083218645595"
	input := "59758034323742284979562302567188059299994912382665665642838883745982029056376663436508823581366924333715600017551568562558429576180672045533950505975691099771937719816036746551442321193912312169741318691856211013074397344457854784758130321667776862471401531789634126843370279186945621597012426944937230330233464053506510141241904155782847336539673866875764558260690223994721394144728780319578298145328345914839568238002359693873874318334948461885586664697152894541318898569630928429305464745641599948619110150923544454316910363268172732923554361048379061622935009089396894630658539536284162963303290768551107950942989042863293547237058600513191659935"
	// input := "03036732577212944063491565474664"
	pattern := []int{0, 1, 0, -1}

	// input = strings.Repeat(input, 10000)
	// fmt.Println(input)

	var b strings.Builder
	fmt.Fprintf(&b, "%s", input)

	// offset, err := strconv.Atoi(input[:7])
	// check(err)

	for i := 0; i < 100; i++ {
		fmt.Println(i)
		b = runPhase(b, pattern)
	}
	fmt.Println(b.String())
	// fmt.Println(input[offset:8])

}

func runPhase(input strings.Builder, pattern []int) strings.Builder {
	var result strings.Builder

	for i := 0; i < input.Len(); i++ {
		sum := 0
		for j, digit := range input.String() {
			intDigit, err := strconv.Atoi(string(digit))
			check(err)

			sum += intDigit * pattern[(((j+1)/(i+1))%4)]
			// fmt.Println(intDigit, (((j + 1) / (i + 1)) % 4), pattern[(((j+1)/(i+1))%4)], sum)
		}
		// sumString := strconv.Itoa(sum)
		fmt.Fprintf(&result, "%d", int(math.Abs(float64(sum%10))))
	}
	// fmt.Println(result.String())
	return result
}
