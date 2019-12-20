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

	"github.com/eiannone/keyboard"
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

func userInput() int {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter an integer parameter: ")
	scanner.Scan()
	input, err := strconv.Atoi(scanner.Text())
	check(err)
	return input
}

func userOutput(output int, done bool) {
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
	paramModes = append(paramModes, make([]int, numParams-len(paramModes))...)

	for j := range params {
		if paramModes[j] == 2 {
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
	x := make([]int, len(program)+1000)
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

	// limiter := time.Tick(time.Millisecond)
	// ship := make([][]int, 1)
	// score := 0

	// // for i := 0; i < (38 * 21); i++ {
	// // 	tile := <-output
	// // 	updatemap(&map, tile[0], tile[1], tile[2])
	// // }

	// paddle, ball, i := 0, 0, 0
	// // setup ui
	// for tile := range output {

	// 	if tile[0] == -1 {
	// 		score = tile[2]
	// 	} else {
	// 		updateShip(&ship, tile[0], tile[1], tile[2])
	// 	}

	// 	if tile[2] == 4 {
	// 		ball = tile[0]
	// 	}

	// 	if tile[2] == 3 {
	// 		paddle = tile[0]
	// 	}

	// 	if ball != 0 && paddle != 0 && tile[2] == 4 {
	// 		diff := ball - paddle
	// 		// fmt.Println("diff:", diff)
	// 		if diff < 0 {
	// 			input <- -1
	// 		} else if diff > 0 {
	// 			input <- 1
	// 		} else {
	// 			input <- 0
	// 		}
	// 	}

	// 	if i > (38 * 21) {
	// 		<-limiter
	// 		// print2D(ship, score)
	// 		// fmt.Println(tile)
	// 		// fmt.Println("ball:", ball, "paddle:", paddle)
	// 	} else {
	// 		i++
	// 	}
	// }
	// fmt.Println(score)

	// total := 0
	// for _, row := range ship {
	// 	for _, tile := range row {
	// 		if tile == 2 {
	// 			total++
	// 		}
	// 	}
	// }
	// fmt.Println("# of blocks:", total)

	ship := [][]string{[]string{"D"}}
	x, x0, y, y0 := 0, 0, 0, 0

	input := make(chan int, 1)
	output := make(chan int, 2)
	go executeIntcode("ship", program, channelInput(input), channelOutput([]chan int{output}))
	go getUserInput(input, &x, &y, &x0, &y0)

	for result := range output {
		// fmt.Println("result:", result, x, y, x0, y0)
		switch result {
		case 0:
			updateShip(&ship, &x, &y, &x0, &y0, ".")
			x, y = x0, y0
		case 1:
			updateShip(&ship, &x0, &y0, &x0, &y0, " ")
			updateShip(&ship, &x, &y, &x0, &y0, "D")
		case 2:
			updateShip(&ship, &x0, &y0, &x0, &y0, " ")
			updateShip(&ship, &x, &y, &x0, &y0, "S")
			fmt.Println("oxygen:", x, y)
			os.Exit(x*100 + y)
		}

		fmt.Println()
		fmt.Println()
		for _, row := range ship {
			fmt.Println(row)
		}
	}

}

func getUserInput(input chan int, x, y, x0, y0 *int) {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		rune, key, err := keyboard.GetKey()
		check(err)
		if key == keyboard.KeyEsc {
			break
		}
		instruction := string(rune)

		code := 0
		if instruction == "w" || instruction == "a" || instruction == "s" || instruction == "d" {
			*x0, *y0 = *x, *y
			switch instruction {
			case "w":
				*y--
				code = 1
			case "s":
				*y++
				code = 2
			case "d":
				*x++
				code = 3
			case "a":
				*x--
				code = 4
			}
		}

		if code != 0 {
			input <- code
		}
	}

	os.Exit(0)
}

func mappingInput(input chan int, x, y, x0, y0 *int) {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		rune, key, err := keyboard.GetKey()
		check(err)
		if key == keyboard.KeyEsc {
			break
		}
		instruction := string(rune)

		code := 0
		if instruction == "w" || instruction == "a" || instruction == "s" || instruction == "d" {
			*x0, *y0 = *x, *y
			switch instruction {
			case "w":
				*y--
				code = 1
			case "s":
				*y++
				code = 2
			case "d":
				*x++
				code = 3
			case "a":
				*x--
				code = 4
			}
		}

		if code != 0 {
			input <- code
		}
	}

	os.Exit(0)
}

func updateShip(ship *[][]string, x, y, x0, y0 *int, t string) {
	h := len(*ship)
	if *y > h-1 {
		newRow := make([]string, len((*ship)[0]))
		for i := range newRow {
			newRow[i] = "?"
		}
		*ship = append(*ship, newRow)
	} else if *y < 0 {
		newRow := make([]string, len((*ship)[0]))
		for i := range newRow {
			newRow[i] = "?"
		}
		*ship = append([][]string{newRow}, *ship...)
		*y = 0
		*y0 = 1
	}

	w := len((*ship)[*y])
	if *x > w-1 {
		for i := range *ship {
			newCol := make([]string, *x-w+1)
			for i := range newCol {
				newCol[i] = "?"
			}
			(*ship)[i] = append((*ship)[i], newCol...)
		}
	} else if *x < 0 {
		for i := range *ship {
			newCol := make([]string, 1)
			for i := range newCol {
				newCol[i] = "?"
			}
			(*ship)[i] = append(newCol, (*ship)[i]...)
		}
		*x = 0
		*x0 = 1
	}

	(*ship)[*y][*x] = t
}

// func shipOutput(x,y *int) func(int, bool) {
// 	ship := [][]string{[]string{"D"}}

// 	return func(output int, done bool) {
// 		switch output{
// 		case 0:
// 		case 0:
// 		case 0:
// 		}
// 		updateShip(&ship, *x, *y, update)
// 		fmt.Println("output:", output)
// 	}
// }
