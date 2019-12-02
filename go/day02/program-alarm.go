package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func readCsv(path string) ([]int, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	inString := string(in)
	// inString = `2,4,4,5,99,0`

	r := csv.NewReader(strings.NewReader(inString))
	records, err := r.Read()
	if err != nil {
		return nil, err
	}

	data := make([]int, len(records))

	for i, record := range records {
		temp, err := strconv.Atoi(record)
		if err != nil {
			return nil, err
		}

		data[i] = temp
	}

	return data, nil
}

func executeIntcode(input []int) []int {
	for i := 0; i < len(input); i += 4 {
		op := input[i]

		if op == 99 {
			break
		}

		val1 := input[input[i+1]]
		val2 := input[input[i+2]]
		iOut := input[i+3]

		switch op {
		case 1:
			input[iOut] = val1 + val2
		case 2:
			input[iOut] = val1 * val2
		case 99:
			break
		default:
			fmt.Printf("Unknown opcode %d encountered\n", op)
		}
	}
	return input
}

func main() {
	input, err := readCsv("input.txt")
	check(err)

	// part1 := executeIntcode(input)
	// fmt.Println(part1)

	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			// Initialize memory
			memory := make([]int, len(input))
			copy(memory, input)

			// Update noun and verb
			memory[1] = noun
			memory[2] = verb

			// Execute
			output := executeIntcode(memory)
			if output[0] == 19690720 {
				fmt.Println(noun, verb)
				break
			}
		}
	}
}
