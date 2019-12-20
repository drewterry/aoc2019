package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type formula map[string]int

type reaction struct {
	output  int
	formula formula
}

type factory map[string]reaction

func readCsv(path string) factory {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	input := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	// input = []string{`10 ORE => 10 A`,
	// 	`1 ORE => 1 B`,
	// 	`7 A, 1 B => 1 C`,
	// 	`7 A, 1 C => 1 D`,
	// 	`7 A, 1 D => 1 E`,
	// 	`7 A, 1 E => 1 FUEL`}

	// input = []string{`9 ORE => 2 A`,
	// 	`8 ORE => 3 B`,
	// 	`7 ORE => 5 C`,
	// 	`3 A, 4 B => 1 AB`,
	// 	`5 B, 7 C => 1 BC`,
	// 	`4 C, 1 A => 1 CA`,
	// 	`2 AB, 3 BC, 4 CA => 1 FUEL`,
	// }

	// input = []string{`157 ORE => 5 NZVS`,
	// 	`165 ORE => 6 DCFZ`,
	// 	`44 XJWVT, 5 KHKGT, 1 QDVJ, 29 NZVS, 9 GPVTF, 48 HKGWZ => 1 FUEL`,
	// 	`12 HKGWZ, 1 GPVTF, 8 PSHF => 9 QDVJ`,
	// 	`179 ORE => 7 PSHF`,
	// 	`177 ORE => 5 HKGWZ`,
	// 	`7 DCFZ, 7 PSHF => 2 XJWVT`,
	// 	`165 ORE => 2 GPVTF`,
	// 	`3 DCFZ, 7 NZVS, 5 HKGWZ, 10 PSHF => 8 KHKGT`}

	// input = []string{`2 VPVL, 7 FWMGM, 2 CXFTF, 11 MNCFX => 1 STKFG`,
	// 	`17 NVRVD, 3 JNWZP => 8 VPVL`,
	// 	`53 STKFG, 6 MNCFX, 46 VJHF, 81 HVMC, 68 CXFTF, 25 GNMV => 1 FUEL`,
	// 	`22 VJHF, 37 MNCFX => 5 FWMGM`,
	// 	`139 ORE => 4 NVRVD`,
	// 	`144 ORE => 7 JNWZP`,
	// 	`5 MNCFX, 7 RFSQX, 2 FWMGM, 2 VPVL, 19 CXFTF => 3 HVMC`,
	// 	`5 VJHF, 7 MNCFX, 9 VPVL, 37 CXFTF => 6 GNMV`,
	// 	`145 ORE => 6 MNCFX`,
	// 	`1 NVRVD => 8 CXFTF`,
	// 	`1 VJHF, 6 MNCFX => 4 RFSQX`,
	// 	`176 ORE => 6 VJHF`}

	input = []string{`171 ORE => 8 CNZTR`,
		`7 ZLQW, 3 BMBT, 9 XCVML, 26 XMNCP, 1 WPTQ, 2 MZWV, 1 RJRHP => 4 PLWSL`,
		`114 ORE => 4 BHXH`,
		`14 VRPVC => 6 BMBT`,
		`6 BHXH, 18 KTJDG, 12 WPTQ, 7 PLWSL, 31 FHTLT, 37 ZDVW => 1 FUEL`,
		`6 WPTQ, 2 BMBT, 8 ZLQW, 18 KTJDG, 1 XMNCP, 6 MZWV, 1 RJRHP => 6 FHTLT`,
		`15 XDBXC, 2 LTCX, 1 VRPVC => 6 ZLQW`,
		`13 WPTQ, 10 LTCX, 3 RJRHP, 14 XMNCP, 2 MZWV, 1 ZLQW => 1 ZDVW`,
		`5 BMBT => 4 WPTQ`,
		`189 ORE => 9 KTJDG`,
		`1 MZWV, 17 XDBXC, 3 XCVML => 2 XMNCP`,
		`12 VRPVC, 27 CNZTR => 2 XDBXC`,
		`15 KTJDG, 12 BHXH => 5 XCVML`,
		`3 BHXH, 2 VRPVC => 7 MZWV`,
		`121 ORE => 7 VRPVC`,
		`7 XCVML => 6 RJRHP`,
		`5 BHXH, 4 VRPVC => 5 LTCX`}

	output := factory{}

	for _, line := range input {
		parts := strings.Split(line, " => ")

		reagents := formula{}
		for _, reagent := range strings.Split(parts[0], ",") {
			var (
				quant int
				name  string
			)
			fmt.Sscanf(reagent, "%d %s", &quant, &name)
			reagents[name] = quant

		}

		var (
			quant  int
			result string
		)
		fmt.Sscanf(parts[1], "%d %s", &quant, &result)

		output[result] = reaction{quant, reagents}
	}

	return output
}

func main() {
	reactor := readCsv("input.txt")

	// for k, v := range reactor {
	// 	fmt.Println(k, v)
	// }

	extra := formula{}
	ore := 1000000000000
	fuel := 0
	for ore > 0 {
		ore -= make1Fuel(reactor, &extra)
		fuel++
		if fuel%10000 == 0 {
			fmt.Println(fuel, ore)
		}
	}
	fmt.Println(extra)
	fmt.Println(fuel - 1)

}

func make1Fuel(reactor factory, extra *formula) int {
	reaction := make([]formula, 13)
	for i := range reaction {
		reaction[i] = formula{}
	}
	reaction[0] = formula{"FUEL": 1}

	// fmt.Println("reaction:", reaction[0])
	// fmt.Println("extra:", *extra)
	// fmt.Println()

	for i := 1; i < 13; i++ {
		for k1, v1 := range reaction[i-1] {
			k1Reaction, k1Exists := reactor[k1]
			_, ore := k1Reaction.formula["ORE"]
			if k1Exists && !ore {
				// fmt.Println(k1, reactor[k1].formula)
				needed := v1
				if (*extra)[k1] >= needed {
					(*extra)[k1] -= needed
					needed = 0
				} else {
					// fmt.Println()
					needed -= (*extra)[k1]
					(*extra)[k1] = 0
				}
				reactions := int(math.Ceil(float64(needed) / float64(reactor[k1].output)))
				produced := reactions * reactor[k1].output
				// fmt.Println(k1, reactions, produced, needed)
				if produced > needed {
					(*extra)[k1] += (produced - needed)
				}
				for k2, v2 := range reactor[k1].formula {
					reaction[i][k2] += reactions * v2
				}
				// fmt.Println(k1, reaction[i])
				// fmt.Println(k1, *extra)
			} else {
				reaction[i][k1] += v1
			}
		}
	}

	ore := 0
	for k, v := range reaction[len(reaction)-1] {
		kReaction := reactor[k]
		// fmt.Println(k, v, kReaction, int(math.Ceil(float64(v)/float64(kReaction.output)))*kReaction.formula["ORE"])
		ore += int(math.Ceil(float64(v)/float64(kReaction.output))) * kReaction.formula["ORE"]
	}
	return ore
}
