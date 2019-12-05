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

func parseOpCode(input int) (int, []int) {
	opCode := input % 100

	paramCode := input / 100
	var paramModes []int

	for paramCode >= 10 {
		paramModes = append(paramModes, paramCode%10)
		paramCode = paramCode / 10
	}
	paramModes = append(paramModes, paramCode)

	return opCode, paramModes
}

func getParams(x []int, addr int, paramModes []int, numParams int) []int {
	params := make([]int, numParams)
	copy(params, paramModes)

	for j, param := range params {
		if param == 0 {
			params[j] = x[x[addr+j+1]]
		} else {
			params[j] = x[addr+j+1]
		}
	}

	return params
}

func executeIntcode(x []int) []int {
execution:
	for i := 0; i < len(x); {
		op, paramModes := parseOpCode(x[i])

		switch op {
		case 1:
			params := getParams(x, i, paramModes, 2)
			x[x[i+3]] = params[0] + params[1]
			i += 4
		case 2:
			params := getParams(x, i, paramModes, 2)
			x[x[i+3]] = params[0] * params[1]
			i += 4
		case 3:
			x[x[i+1]] = readInput()
			i += 2
		case 4:
			params := getParams(x, i, paramModes, 1)
			fmt.Println(params[0])
			i += 2
		case 5:
			params := getParams(x, i, paramModes, 2)
			if params[0] != 0 {
				i = params[1]
			} else {
				i += 3
			}
		case 6:
			params := getParams(x, i, paramModes, 2)
			if params[0] == 0 {
				i = params[1]
			} else {
				i += 3
			}
		case 7:
			params := getParams(x, i, paramModes, 2)
			if params[0] < params[1] {
				x[x[i+3]] = 1
			} else {
				x[x[i+3]] = 0
			}
			i += 4
		case 8:
			params := getParams(x, i, paramModes, 2)
			if params[0] == params[1] {
				x[x[i+3]] = 1
			} else {
				x[x[i+3]] = 0
			}
			i += 4
		case 99:
			break execution
		default:
			fmt.Printf("Unknown opcode %d encountered\n", op)
			break execution
		}
	}
	return x
}

func main() {
	input, err := readCsv("input.txt")
	check(err)

	memory := input
	// memory = []int{3, 0, 4, 0, 99}
	// memory = []int{1002, 5, 3, 5, 99, 33}
	// memory = []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}
	// memory = []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}
	// memory = []int{3, 3, 1108, -1, 8, 3, 4, 3, 99}
	// memory = []int{3, 3, 1107, -1, 8, 3, 4, 3, 99}
	// memory = []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}
	// memory = []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}
	memory = []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31, 1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104, 999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}
	executeIntcode(memory)
	// fmt.Println(output)
}
