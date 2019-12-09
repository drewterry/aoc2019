package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	inputBytes, err := ioutil.ReadFile("input.txt")
	check(err)

	inputString := string(inputBytes)
	// inputString = "123456789012"
	// inputString = "0222112222120000"

	input := strings.Split(inputString, "")

	data := make([]int, len(input))
	for i, pixel := range input {
		temp, err := strconv.Atoi(pixel)
		check(err)

		data[i] = temp
	}
	// fmt.Println(data)

	width := 25
	height := 6
	depth := int(len(data) / width / height)

	image := make([][][]int, depth)
	for i := 0; i < depth; i++ {
		image[i] = make([][]int, height)
		for j := 0; j < height; j++ {
			image[i][j] = make([]int, width)
			for k := 0; k < width; k++ {
				fmt.Println(i, j, k, i*width*height+j*width+k)
				image[i][j][k] = data[i*width*height+j*width+k]
			}
		}
	}
	// fmt.Println(image)

	finalImage := image[0]
	for i := range image {
		// fmt.Println(i)
		// for _, row := range layer {
		// 	fmt.Println(row)
		// }

		for j := 0; j < height; j++ {
			for k := 0; k < width; k++ {
				if finalImage[j][k] == 2 {
					finalImage[j][k] = image[i][j][k]
				}
			}
		}
	}

	output := make([][]string, len(finalImage))
	for i, row := range finalImage {
		output[i] = make([]string, len(row))
		for j, el := range row {
			if el == 0 {
				output[i][j] = "."
			} else {
				output[i][j] = "X"
			}
		}
	}
	for _, row := range output {
		fmt.Println(row)
	}

	// min := layerZeroes[0]
	// minI := 0
	// for i, e := range layerZeroes {
	// 	if i == 0 || e < min {
	// 		min = e
	// 		minI = i
	// 	}
	// }
	// fmt.Println(min, minI)

	// layer := image[minI]
	// for _, row := range layer {
	// 	fmt.Println(row)
	// }
	// layerTwos := countInts2D(layer, 2)
	// layerOnes := countInts2D(layer, 1)
	// fmt.Println(layerTwos, layerOnes, layerOnes*layerTwos)

}

func countInts2D(layer [][]int, integer int) int {
	count := 0
	for i := 0; i < len(layer); i++ {
		for j := 0; j < len(layer[i]); j++ {
			if layer[i][j] == integer {
				count++
			}
		}
	}
	return count
}
