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
	// inString = `109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99`
	// inString = `1102,34915192,34915192,7,4,7,99,0`
	// inString = `104,1125899906842624,99`
	// inString = `109,19,99`

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

func getUserInput(name string) int {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter an integer parameter: ")
	scanner.Scan()
	input, err := strconv.Atoi(scanner.Text())
	check(err)
	return input
}

func writeUserOutput(name string, output int) {
	fmt.Println("output:", output)
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

type param struct {
	index int
	value int
}

func getParams(x []int, addr int, paramModes []int, numParams int, base int) []param {
	params := make([]param, numParams)
	// fmt.Println(x[addr], params, paramModes, numParams, len(paramModes))
	paramModes = append(paramModes, make([]int, numParams-len(paramModes))...)

	for j := range params {
		if paramModes[j] == 2 {
			// fmt.Println(base, addr, j, x)
			params[j] = param{base + x[addr+j+1], x[base+x[addr+j+1]]}
		} else if paramModes[j] == 1 {
			params[j] = param{addr + j + 1, x[addr+j+1]}
		} else {
			params[j] = param{x[addr+j+1], x[x[addr+j+1]]}
		}
	}

	return params
}

func executeIntcode(name string, x []int, readInput func(string) int, writeOutput func(string, int), outputs []chan int) []int {
	relativeBase := 0

execution:
	for i := 0; i < len(x); {
		op, paramModes := parseOpCode(x[i])
		// fmt.Println(name, "op:", op)

		switch op {
		case 1:
			params := getParams(x, i, paramModes, 3, relativeBase)
			x[params[2].index] = params[0].value + params[1].value
			i += 4
		case 2:
			params := getParams(x, i, paramModes, 3, relativeBase)
			x[params[2].index] = params[0].value * params[1].value
			i += 4
		case 3:
			params := getParams(x, i, paramModes, 1, relativeBase)
			x[params[0].index] = readInput(name)
			// fmt.Println(relativeBase, paramModes, params, x[params[0].index], i, x[:i], x)
			i += 2
		case 4:
			params := getParams(x, i, paramModes, 1, relativeBase)
			writeOutput(name, params[0].value)
			i += 2
		case 5:
			params := getParams(x, i, paramModes, 2, relativeBase)
			if params[0].value != 0 {
				i = params[1].value
			} else {
				i += 3
			}
		case 6:
			params := getParams(x, i, paramModes, 2, relativeBase)
			if params[0].value == 0 {
				i = params[1].value
			} else {
				i += 3
			}
		case 7:
			params := getParams(x, i, paramModes, 3, relativeBase)
			if params[0].value < params[1].value {
				x[params[2].index] = 1
			} else {
				x[params[2].index] = 0
			}
			i += 4
		case 8:
			params := getParams(x, i, paramModes, 3, relativeBase)
			if params[0].value == params[1].value {
				x[params[2].index] = 1
			} else {
				x[params[2].index] = 0
			}
			i += 4
		case 9:
			params := getParams(x, i, paramModes, 1, relativeBase)
			relativeBase += params[0].value
			i += 2
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

	memory := make([]int, len(input)+1000)
	copy(memory, input)
	// memory = []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5}
	// memory = []int{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10}
	// memory = []int{3, 15, 3, 16, 1002, 16, 10, 16, 1, 16, 15, 15, 4, 15, 99, 0, 0}
	// memory = []int{3, 0, 4, 0, 99}
	// memory = []int{1002, 5, 3, 5, 99, 33}
	// memory = []int{3, 9, 8, 9, 10, 9, 4, 9, 99, -1, 8}
	// memory = []int{3, 9, 7, 9, 10, 9, 4, 9, 99, -1, 8}
	// memory = []int{3, 3, 1108, -1, 8, 3, 4, 3, 99}
	// memory = []int{3, 3, 1107, -1, 8, 3, 4, 3, 99}
	// memory = []int{3, 12, 6, 12, 15, 1, 13, 14, 13, 4, 13, 99, -1, 0, 1, 9}
	// memory = []int{3, 3, 1105, -1, 9, 1101, 0, 0, 12, 4, 12, 99, 1}
	// memory = []int{3, 21, 1008, 21, 8, 20, 1005, 20, 22, 107, 8, 21, 20, 1006, 20, 31, 1106, 0, 36, 98, 0, 0, 1002, 21, 125, 20, 4, 20, 1105, 1, 46, 104, 999, 1105, 1, 46, 1101, 1000, 1, 20, 4, 20, 1105, 1, 46, 98, 99}

	executeIntcode("boost", memory, getUserInput, writeUserOutput, make([]chan int, 0))

}
