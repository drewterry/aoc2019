package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func userInput(name string) int {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter an integer parameter: ")
	scanner.Scan()
	input, err := strconv.Atoi(scanner.Text())
	check(err)
	return input
}

func userOutput(name string, output int) {
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

func executeIntcode(name string, program []int, readInput func() int, writeOutput func(int, bool)) []int {
	x := make([]int, len(program)+3000)
	copy(x, program)
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
			x[params[0].index] = readInput()
			// fmt.Println(relativeBase, paramModes, params, x[params[0].index], i, x[:i], x)
			i += 2
		case 4:
			params := getParams(x, i, paramModes, 1, relativeBase)
			writeOutput(params[0].value, false)
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
			fmt.Println(name, "complete")
			writeOutput(0, true)
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

func staticInput(inputs []int) func() int {
	i := 0
	return func() int {
		i++
		return inputs[i-1]
	}
}

func channelInput(inputChan chan int) func() int {
	return func() int {
		input := <-inputChan
		return input
	}
}

func channelInputDefault(inputChan chan int, defaultInput int) func() int {
	return func() int {
		input := defaultInput
		select {
		case input = <-inputChan:
		default:
		}
		return input
	}
}

func initalizeChannel(channel chan int, inputs []int) {
	for _, input := range inputs {
		channel <- input
	}
}

func pointerOutput(outputPtr *int) func(int) {
	return func(output int) {
		*outputPtr = output
	}
}

func channelOutput(outputChans []chan int) func(int, bool) {
	return func(output int, done bool) {
		for _, outputChan := range outputChans {
			if done {
				close(outputChan)
			} else {
				outputChan <- output
			}
		}
	}
}

func channelOutputCollate(collated chan []int, num int) func(int, bool) {
	output := make(chan int, num)
	tempOutput := make([]int, 0)

	go func() {
		for out := range output {
			tempOutput = append(tempOutput, out)

			if len(tempOutput) == num {
				collated <- tempOutput
				tempOutput = make([]int, 0)
			}
		}
		fmt.Println("output closed")
		close(collated)
	}()

	return channelOutput([]chan int{output})
}

func executeChain(memory []int, settings []int, initalVal int) int {

	programs := make([][]int, len(settings))
	chans := make([]chan int, len(settings))

	for i := range programs {
		programs[i] = make([]int, len(memory))
		copy(programs[i], memory)
		chans[i] = make(chan int, 2)
		initalizeChannel(chans[i], []int{settings[i]})
	}

	chans[0] <- 0
	outputs := make(chan int, 2)
	initalizeChannel(outputs, []int{})

	for i := 0; i < len(settings)-1; i++ {
		fmt.Println(i)
		go executeIntcode(fmt.Sprint("amp", i), programs[i], channelInput(chans[i]), channelOutput([]chan int{chans[i+1]}))
	}

	lastChan := len(settings) - 1
	go executeIntcode(fmt.Sprint("amp", lastChan), programs[lastChan], channelInput(chans[lastChan]), channelOutput([]chan int{chans[0], outputs}))

	final := 0
	for output := range outputs {
		final = output
		fmt.Println("output: ", output)
	}
	return final
}

type coord struct {
	x int
	y int
}

type robot struct {
	loc coord
	dir string
}

func main() {
	program, err := readCsv("input.txt")
	check(err)

	// memory = []int{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5}

	input := make(chan int, 2)
	input <- 0

	output := make(chan int, 2)
	go executeIntcode("game", program, channelInput(input), channelOutput([]chan int{output}))

	scaffold := make([][]string, 1)
	scaffold[0] = make([]string, 0)
	line := 0

	for code := range output {
		switch code {
		case 10:
			scaffold = append(scaffold, make([]string, 0))
			line++
		default:
			scaffold[line] = append(scaffold[line], string(code))
		}
	}

	scaffold = scaffold[:len(scaffold)-2]

	for i, line := range scaffold {
		fmt.Println(i, line)
	}

	type coord struct {
		x int
		y int
	}

	intersections := make([]coord, 0)
	for i := 1; i < len(scaffold)-1; i++ {
		for j := 1; j < len(scaffold[i])-1; j++ {
			if scaffold[i][j] == "#" &&
				scaffold[i+1][j] == "#" &&
				scaffold[i-1][j] == "#" &&
				scaffold[i][j+1] == "#" &&
				scaffold[i][j-1] == "#" {
				intersections = append(intersections, coord{i, j})
			}
		}
	}

	fmt.Println(intersections)

	total := 0
	for _, intersection := range intersections {
		total += intersection.x * intersection.y
	}
	fmt.Println(total)

	// PART 2 +++++++

	x0, y0 := 0, 0

start:
	for i := 0; i < len(scaffold); i++ {
		for j := 0; j < len(scaffold[i]); j++ {
			if scaffold[i][j] == "^" {
				x0, y0 = i, j
				break start
			}
		}
	}
	fmt.Println(x0, y0)

	x, y := x0, y0
	xD, yD := nextDir(scaffold, x, y, 0, 0)

	for xD != 0 || yD != 0 {

		length := 0
		for {
			if x+xD < 0 || x+xD >= len(scaffold) || y+yD < 0 || y+yD >= len(scaffold[0]) || scaffold[x+xD][y+yD] != "#" {
				break
			} else {
				x, y = x+xD, y+yD
				length++
			}
		}

		fmt.Println(length)

		xD, yD = nextDir(scaffold, x, y, xD, yD)
	}

	// a := "L,12,L,10,R,8,L,12,R,8,R,10,R,12"
	aCode := []int{76, 44, 49, 50, 44, 76, 44, 49, 48, 44, 82, 44, 56, 44, 76, 44, 49, 50, 10}
	// b := "L,10,R,12,R,8"
	bCode := []int{76, 44, 49, 48, 44, 82, 44, 49, 50, 44, 82, 44, 56, 10}
	// c := "R,8,R,10,R,12"
	cCode := []int{82, 44, 56, 44, 82, 44, 49, 48, 44, 82, 44, 49, 50, 10}
	// routine := "A,A,B,B,C,A,B"  "A,C,A,C,B,B,C,A,C,B"
	routineCode := []int{65, 44, 67, 44, 65, 44, 67, 44, 66, 44, 66, 44, 67, 44, 65, 44, 67, 44, 66, 10}

	program[0] = 2
	inputCodes := append(routineCode, aCode...)
	inputCodes = append(inputCodes, bCode...)
	inputCodes = append(inputCodes, cCode...)
	input = make(chan int, len(inputCodes)+2)
	for _, code := range inputCodes {
		input <- code
	}
	input <- 89
	input <- 10

	output = make(chan int, 2)
	go executeIntcode("game", program, channelInput(input), channelOutput([]chan int{output}))

	last := 0
	for out := range output {
		if last == 10 && out == 10 {
			clearScreen()
		}

		if out > 100 {
			fmt.Println(out)
		}

		fmt.Print(string(out))

		last = out
	}
}

func nextDir(scaffold [][]string, x, y, xD, yD int) (int, int) {

	if xD != 0 || (xD == 0 && yD == 0) {
		if y-1 >= 0 && scaffold[x][y-1] == "#" {
			if xD > 0 {
				fmt.Println("R")
			} else {
				fmt.Println("L")
			}
			return 0, -1
		}
		if scaffold[x][y+1] == "#" {
			if xD > 0 {
				fmt.Println("L")
			} else {
				fmt.Println("R")
			}
			return 0, 1
		}
	}
	if yD != 0 || (xD == 0 && yD == 0) {
		if x-1 >= 0 && scaffold[x-1][y] == "#" {
			if yD > 0 {
				fmt.Println("L")
			} else {
				fmt.Println("R")
			}
			return -1, 0
		}
		if scaffold[x+1][y] == "#" {
			if yD > 0 {
				fmt.Println("R")
			} else {
				fmt.Println("L")
			}
			return 1, 0
		}
	}
	return 0, 0
}
func clearScreen() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
