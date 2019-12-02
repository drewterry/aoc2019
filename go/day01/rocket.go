package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// func main() {
// 	dat, err := ioutil.ReadFile("input.txt")
// 	check(err)

// 	datString := string(dat)
// 	arr := strings.Split(datString, "\n")
// 	fmt.Print(arr)
// }

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []int
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return lines, err
		}
		lines = append(lines, x)
	}
	return lines, scanner.Err()
}

func calcFuel(mass int) int {
	return int(math.Floor(float64(mass)/3) - 2)
}

func fuelNeeded(mass int) int {
	fuel := calcFuel(mass)
	fmt.Println(fuel)
	if fuel < 0 {
		return 0
	}
	return fuel + fuelNeeded(fuel)
}

func main() {
	masses, err := readLines("input.txt")
	check(err)

	var total int

	for i, mass := range masses {
		total += fuelNeeded(mass)
		fmt.Println(i, mass, total)
	}

	fmt.Println(total)
}
