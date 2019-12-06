package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func calcOrbits(orbits map[string]string, target string) int {
	orbitCount := 0

	for orbit := orbits[target]; ; {
		if orbit == "COM" {
			break
		}

		orbitCount++
		orbit = orbits[orbit]
	}
	orbitCount++

	return orbitCount
}

func getOrbitalSlice(orbits map[string]string, target string) []string {
	orbitSlice := make([]string, 0)

	for orbit := orbits[target]; ; {
		if orbit == "COM" {
			break
		}

		orbitSlice = append(orbitSlice, orbit)
		orbit = orbits[orbit]
	}

	return orbitSlice
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func main() {
	input, err := readLines("./input.txt")
	check(err)
	// input = []string{`COM)B`, `B)C`, `C)D`, `D)E`, `E)F`, `B)G`, `G)H`, `D)I`, `E)J`, `J)K`, `K)L`}
	// input = []string{"COM)B", "B)C", "C)D", "D)E", "E)F", "B)G", "G)H", "D)I", "E)J", "J)K", "K)L", "K)YOU", "I)SAN"}

	orbits := make(map[string]string)
	for _, inputString := range input {
		orbit := strings.Split(inputString, ")")
		orbits[orbit[1]] = orbit[0]
	}

	totalOrbits := 0
	for k := range orbits {
		totalOrbits += calcOrbits(orbits, k)
	}
	fmt.Println(totalOrbits)

	you := getOrbitalSlice(orbits, "YOU")
	santa := getOrbitalSlice(orbits, "SAN")

findCommonOrbit:
	for i, youOrbit := range you {
		for j, santaOrbit := range santa {
			if youOrbit == santaOrbit {
				fmt.Println(i + j)
				break findCommonOrbit
			}
		}
	}
}
