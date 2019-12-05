package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func readInput() int {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter an integer parameter: ")
	scanner.Scan()
	input, err := strconv.Atoi(scanner.Text())
	check(err)
	return input
}

func executeIntcode(x []int) []int {
	for i := 0; i < len(x); {
		op := x[i]

		if op == 99 {
			break
		}

		switch op {
		case 1:
			val1 := x[x[i+1]]
			val2 := x[x[i+2]]
			iOut := x[i+3]
			x[iOut] = val1 + val2
			i += 4
		case 2:
			val1 := x[x[i+1]]
			val2 := x[x[i+2]]
			iOut := x[i+3]
			x[iOut] = val1 * val2
			i += 4
		case 3:
			addr := x[i+1]
			x[addr] = readInput()
			i += 2
		case 4:
			addr := x[i+1]
			fmt.Println(x[addr])
			i += 2
		case 99:
			break
		default:
			fmt.Printf("Unknown opcode %d encountered\n", op)
			break
		}
	}
	return x
}

func main() {
	input, err := readCsv("input.txt")
	check(err)

	memory := input
	memory = []int{3, 0, 4, 0, 99}
	output := executeIntcode(memory)
	fmt.Println(output)
}
