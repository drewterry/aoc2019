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

func getUserInput() int {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter an integer parameter: ")
	scanner.Scan()
	input, err := strconv.Atoi(scanner.Text())
	check(err)
	return input
}

func writeUserOutput(output int) {
	fmt.Println(output)
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

func executeIntcode(name string, x []int, readInput func(string) int, writeOutput func(string, int), outputs []chan int) []int {
execution:
	for i := 0; i < len(x); {
		op, paramModes := parseOpCode(x[i])
		// fmt.Println(name, "op:", op)

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
			x[x[i+1]] = readInput(name)
			i += 2
		case 4:
			params := getParams(x, i, paramModes, 1)
			writeOutput(name, params[0])
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
			fmt.Println(name, "complete", outputs)
			for _, output := range outputs {
				close(output)
			}
			break execution
		default:
			fmt.Printf("Unknown opcode %d encountered\n", op)
			break execution
		}
	}
	return x
}

// Perm calls f with each permutation of a.
func Perm(a []int, f func([]int)) {
	perm(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func perm(a []int, f func([]int), i int) {
	if i > len(a) {
		f(a)
		return
	}
	perm(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		perm(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func getProgrammaticInput(inputs []int) func() int {
	i := 0
	return func() int {
		i++
		return inputs[i-1]
	}
}

func getProgrammaticInputChannel(inputChan chan int) func(string) int {
	return func(name string) int {
		input := <-inputChan
		// fmt.Println(name, "input:", input)
		return input
	}
}

func createChannelInitial(inputs []int) chan int {
	i := make(chan int, 2)
	for _, input := range inputs {
		i <- input
	}
	return i
}

func writeProgrammaticOutput(outputPtr *int) func(int) {
	return func(output int) {
		*outputPtr = output
	}
}

func writeProgrammaticOutputChannel(outputChans []chan int) func(string, int) {
	return func(name string, output int) {
		// fmt.Println(name, "output:", output)
		for _, outputChan := range outputChans {
			outputChan <- output
		}
	}
}

func executeChain(memory []int, settings []int, initalVal int) int {

	programs := make([][]int, len(settings))
	chans := make([]chan int, len(settings))

	for i := range programs {
		programs[i] = make([]int, len(memory))
		copy(programs[i], memory)
		chans[i] = createChannelInitial([]int{settings[i]})
	}

	chans[0] <- 0
	outputs := createChannelInitial([]int{})

	for i := 0; i < len(settings)-1; i++ {
		fmt.Println(i)
		go executeIntcode(fmt.Sprint("amp", i), programs[i], getProgrammaticInputChannel(chans[i]), writeProgrammaticOutputChannel([]chan int{chans[i+1]}), []chan int{chans[i+1]})
	}

	lastChan := len(settings) - 1
	go executeIntcode(fmt.Sprint("amp", lastChan), programs[lastChan], getProgrammaticInputChannel(chans[lastChan]), writeProgrammaticOutputChannel([]chan int{chans[0], outputs}), []chan int{chans[0], outputs})

	final := 0
	for output := range outputs {
		final = output
		fmt.Println("output: ", output)
	}
	return final
}

func main() {
	input, err := readCsv("input.txt")
	check(err)

	memory := input
	// memory = []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5}
	// memory = []int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10}
	// memory = []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0}
	outputs := []int{}
	Perm([]int{5, 6, 7, 8, 9}, func(a []int) {
		outputs = append(outputs, executeChain(memory, a, 0))
	})

	// fmt.Println(executeChain(memory, []int{5, 6, 7, 8, 9}, 0))
	fmt.Println(executeChain(memory, []int{9, 8, 7, 6, 5}, 0))

	max := 0
	for i, e := range outputs {
		if i == 0 || e > max {
			max = e
		}
	}
	fmt.Println(max)
}
