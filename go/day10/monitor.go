package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

type dir struct {
	rise int
	run  int
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
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
	input, err := readLines("input.txt")
	check(err)

	// input = []string{`.A..B`, `.....`, `CDEFG`, `....H`, `...IJ`}
	// input = []string{`#.........`, `...A......`, `...B..a...`, `.EDCG....a`, `..F.c.b...`, `.....c....`, `..efd.c.gb`, `.......c..`, `....f...c.`, `...e..d..c`}
	// input = []string{`......#.#.`, `#..#.#....`, `..#######.`, `.#.#.###..`, `.#..#.....`, `..#....#.#`, `#..#....#.`, `.##.#..###`, `##...#..#.`, `.#....####`}
	// input = []string{`#.#...#.#.`, `.###....#.`, `.#....#...`, `##.#.#.#.#`, `....#.#.#.`, `.##..###.#`, `..#...##..`, `..##....##`, `......#...`, `.####.###.`}
	// input = []string{`.#..#..###`, `####.###.#`, `....###.#.`, `..###.##.#`, `##.##.#.#.`, `....###..#`, `..#.#..#.#`, `#..#.#.###`, `.##...##.#`, `.....#.#..`}
	// input = []string{`.#....#####...#..`, `##...##.#####..##`, `##...#...#.#####.`, `..#.....#...###..`, `..#.#.....#....##`}
	for _, row := range input {
		fmt.Println(strings.Split(row, ""))
	}

	// fmt.Println(countVisible(input, 1, 0))

	height, width := len(input), len(input[0])

	for _, dir := range getDirs(2*(height), 2*(width)) {
		fmt.Println(dir, getPosDegrees(dir))
	}

	// viz(input, 10, 4)

	counts := make([][]int, height)
	for i := 0; i < height; i++ {
		counts[i] = make([]int, width)
		for j := 0; j < width; j++ {
			if input[i][j] != '.' {
				count := countVisible(input, j, i)
				// fmt.Println(i, j, count)
				counts[i][j] = count
			}
		}
	}

	// for _, row := range counts {
	// 	fmt.Println(row)
	// }

	// viz(input, 0, 0)

	maxIndex := []int{0, 0}
	max := 0
	for i, row := range counts {
		for j, count := range row {
			if i == 0 && j == 0 || count > max {
				maxIndex = []int{i, j}
				max = count
			}
		}
	}
	fmt.Println(maxIndex, max)

	x, y := maxIndex[1], maxIndex[0]
	current := make([]string, height)
	copy(current, input)
	// coords := make([]int, 2)
	count := 0
	dirs := getDirs(2*(height), 2*(width))
	for _, dir := range dirs {
		// fmt.Println(dir, getPosDegrees(dir))

		for i, j := x+dir.run, y+dir.rise; i < width && i >= 0 && j < height && j >= 0; i, j = i+dir.run, j+dir.rise {
			// fmt.Println(i, j, string(current[j][i]))
			if string(current[j][i]) != "." {
				current[j] = replaceAtIndex(current[j], '.', i)
				count++
				fmt.Println(count, j, i, dir)
				// coords = []int{i, j}
				break
			}
		}
	}

	for _, row := range current {
		fmt.Println(strings.Split(row, ""))
	}

	for _, dir := range dirs {
		// fmt.Println(dir, getPosDegrees(dir))

		for i, j := x+dir.run, y+dir.rise; i < width && i >= 0 && j < height && j >= 0; i, j = i+dir.run, j+dir.rise {
			// fmt.Println(i, j, string(current[j][i]))
			if string(current[j][i]) != "." {
				current[j] = replaceAtIndex(current[j], '.', i)
				count++
				fmt.Println(count, j, i, dir)
				// coords = []int{i, j}
				break
			}
		}
	}

	for _, row := range current {
		fmt.Println(strings.Split(row, ""))
	}

	for _, dir := range dirs {
		// fmt.Println(dir, getPosDegrees(dir))

		for i, j := x+dir.run, y+dir.rise; i < width && i >= 0 && j < height && j >= 0; i, j = i+dir.run, j+dir.rise {
			// fmt.Println(i, j, string(current[j][i]))
			if string(current[j][i]) != "." {
				current[j] = replaceAtIndex(current[j], '.', i)
				count++
				fmt.Println(count, j, i, dir)
				// coords = []int{i, j}
				break
			}
		}
	}

	for _, row := range current {
		fmt.Println(strings.Split(row, ""))
	}

	for _, dir1 := range dirs {
		// fmt.Println(dir1, getPosDegrees(dir1))
		if (dir1 == dir{-1, -8}) {
			fmt.Println("!!!", x, y, x+dir1.run, y+dir1.rise, width, height)
		}
		for i, j := x+dir1.run, y+dir1.rise; i < width && i >= 0 && j < height && j >= 0; i, j = i+dir1.run, j+dir1.rise {
			if (dir1 == dir{-1, -8}) {
				fmt.Println(i, j)
			}
			// fmt.Println(i, j, string(current[j][i]))
			if string(current[j][i]) != "." {
				current[j] = replaceAtIndex(current[j], '.', i)
				count++
				fmt.Println(count, j, i, dir1)
				if (dir1 == dir{-8, -1}) {
					fmt.Println("***")
				}
				// coords = []int{i, j}
				break
			}
		}
	}

	// for _, row := range current {
	// 	fmt.Println(strings.Split(row, ""))
	// }

}

func countVisible(input []string, x int, y int) int {
	height, width := len(input), len(input[0])
	dirs := getDirs(2*(height), 2*(width))

	dirArray := make([][]string, len(dirs))
	dirCount := 0

	for dirIndex, dir := range dirs {
		for i, j := x+dir.run, y+dir.rise; i < width && i >= 0 && j < height && j >= 0; i, j = i+dir.run, j+dir.rise {
			dirArray[dirIndex] = append(dirArray[dirIndex], string(input[j][i]))
			// fmt.Println(dir, j, i, string(input[j][i]), dirCount)
			if string(input[j][i]) != "." {
				dirCount++
				break
			}
		}
	}

	return dirCount
}

func getDirs(height int, width int) []dir {
	dirs := make(map[dir]bool)
	for i := -(height/2 + 1); i <= height/2+1; i++ {
		for j := -(width/2 + 1); j <= width/2+1; j++ {
			k, l := float64(i), float64(j)

			multiple := float64(GCD(i, j))
			if multiple != 0 {
				k, l = math.Copysign(k/multiple, k), math.Copysign(l/multiple, l)
			}

			dirs[dir{int(k), int(l)}] = true
		}
	}

	delete(dirs, dir{0, 0})

	keys := make([]dir, 0, len(dirs))
	for k := range dirs {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return getPosDegrees(keys[i]) < getPosDegrees(keys[j])
	})

	// fmt.Println((keys))

	return keys
}

func getSign(x int) int {
	if x < 0 {
		return -1
	}
	return 1
}

func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func viz(inputs []string, x int, y int) {
	height, width := len(inputs), len(inputs[0])

	viz := make([][]string, height)
	for i := range viz {
		viz[i] = make([]string, width)
		for j := range viz[i] {
			viz[i][j] = "-"
		}
	}

	for _, dir := range getDirs(2*height, 2*width) {
		viz[height-1+dir.rise][width-1+dir.run] = "#"
	}
	viz[height-1][width-1] = "*"
	for _, row := range viz {
		fmt.Println(row)
	}
}

func getPosDegrees(dir dir) float64 {
	deg := 180 / math.Pi * math.Atan2(float64(dir.rise), float64(dir.run))
	deg = 90 + deg
	if deg < 0 {
		deg = 360 + deg
	}
	return deg
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}
