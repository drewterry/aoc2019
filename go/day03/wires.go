package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
)

type step struct {
	dir  string
	dist int
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func readCsv(path string) ([][]step, error) {
	in, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	inString := string(in)
	// 	inString = `R8,U5,L5,D3
	// U7,R6,D4,L4`

	r := csv.NewReader(strings.NewReader(inString))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	data := make([][]step, len(records))

	for i, record := range records {
		data[i] = make([]step, len(record))

		for j, instruction := range record {
			dist, err := strconv.Atoi(instruction[1:])
			if err != nil {
				return nil, err
			}

			data[i][j] = step{string(instruction[0]), dist}
		}
	}

	return data, nil
}

func main() {
	input, err := readCsv("input.txt")
	check(err)

	// for i, wire := range input {
	// 	fmt.Println(i, wire)
	// }

	wires := make([][][2]int, len(input))

	for i, wire := range input {

		for j, instruction := range wire {

			for k := 0; k < instruction.dist; k++ {

				var lastInstruction [2]int
				if (j + k) > 0 {
					lastInstruction = wires[i][len(wires[i])-1]
				}

				dx, dy := 0, 0
				switch instruction.dir {
				case "U":
					dy = 1
				case "D":
					dy = -1
				case "R":
					dx = 1
				case "L":
					dx = -1
				}

				wires[i] = append(wires[i], [2]int{lastInstruction[0] + dx, lastInstruction[1] + dy})
			}
		}
	}

	for i, wire := range wires {
		fmt.Println(i, wire)
	}

	matches := make([][2]int, 0)
	steps := make([]int, 0)

	for i, pos1 := range wires[0] {
		for j, pos2 := range wires[1] {
			if pos1 == pos2 {
				matches = append(matches, pos1)
				steps = append(steps, i+j+2)
				break
			}
		}
	}

	fmt.Println(matches)

	distances := make([]int, 0)
	for _, match := range matches {
		distances = append(distances, int(math.Abs(float64(match[0]))+math.Abs(float64(match[1]))))
	}

	fmt.Println(distances, steps)

	var m int
	for i, e := range distances {
		if i == 0 || e < m {
			m = e
		}
	}
	fmt.Println("Shortest Distance: ", m)

	for i, e := range steps {
		if i == 0 || e < m {
			m = e
		}
	}

	fmt.Println("Shortest Steps: ", m)

}
