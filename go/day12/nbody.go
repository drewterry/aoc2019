// <x=-6, y=2, z=-9>
// <x=12, y=-14, z=-4>
// <x=9, y=5, z=-6>
// <x=-1, y=-4, z=9>

package main

import (
	"fmt"
	"math"
)

type axis struct {
	pos int
	vel int
}

type coord struct {
	x int
	y int
	z int
}

type moon struct {
	x axis
	y axis
	z axis
}

// <x=-1, y=0, z=2>
// <x=2, y=-10, z=-7>
// <x=4, y=-8, z=8>
// <x=3, y=5, z=-1>

// <x=-8, y=-10, z=0>
// <x=5, y=5, z=10>
// <x=2, y=-7, z=3>
// <x=9, y=-8, z=-3>

func main() {

	// moons := []moon{
	// 	{coord{-6, 2, -9}, coord{0, 0, 0}},
	// 	{coord{12, -14, -4}, coord{0, 0, 0}},
	// 	{coord{9, 5, -6}, coord{0, 0, 0}},
	// 	{coord{-1, -4, 9}, coord{0, 0, 0}},
	// }

	// moons := []moon{
	// 	{coord{-1, 0, 2}, coord{0, 0, 0}},
	// 	{coord{2, -10, -7}, coord{0, 0, 0}},
	// 	{coord{4, -8, 8}, coord{0, 0, 0}},
	// 	{coord{3, 5, -1}, coord{0, 0, 0}},
	// }

	moons := []moon{
		{axis{-1, 0}, axis{0, 0}, axis{2, 0}},
		{axis{2, 0}, axis{-10, 0}, axis{-7, 0}},
		{axis{4, 0}, axis{-8, 0}, axis{8, 0}},
		{axis{3, 0}, axis{5, 0}, axis{-1, 0}},
	}

	// moons := []moon{
	// 	{coord{-8, -10, 0}, coord{0, 0, 0}},
	// 	{coord{5, 5, 10}, coord{0, 0, 0}},
	// 	{coord{2, -7, 3}, coord{0, 0, 0}},
	// 	{coord{9, -8, -3}, coord{0, 0, 0}},
	// }

	// moons := []moon{
	// 	{axis{-8, 0}, axis{-10, 0}, axis{0, 0}},
	// 	{axis{5, 0}, axis{5, 0}, axis{10, 0}},
	// 	{axis{2, 0}, axis{-7, 0}, axis{3, 0}},
	// 	{axis{9, 0}, axis{-8, 0}, axis{-3, 0}},
	// }

	// initial := make([]moon, len(moons))
	// copy(initial, moons)

	// printMoons(moons)
	// for i := 0; i < 100; i++ {
	// 	moons = stepVelocity(moons)
	// 	moons = stepPosition(moons)
	// }
	// printMoons(moons)

	// totalEnergy := 0
	// for _, moon := range moons {
	// 	totalEnergy += calcEnergy(moon)
	// }
	// fmt.Println(totalEnergy)

	axesInitial := make([][]axis, 3)
	for i := range axesInitial {
		axesInitial[i] = make([]axis, len(moons))
	}
	for j, moon := range moons {
		axesInitial[0][j] = moon.x
		axesInitial[1][j] = moon.y
		axesInitial[2][j] = moon.z
	}

	fmt.Println(moons)
	fmt.Println(axesInitial)

	periods := []int{0, 0, 0}

	iter := 0
	for {
		// fmt.Println(iter + 1)
		iter++
		moons = stepVelocity(moons)
		moons = stepPosition(moons)

		axes := make([][]axis, 3)
		for i := range axes {
			axes[i] = make([]axis, len(moons))
		}
		for j, moon := range moons {
			axes[0][j] = moon.x
			axes[1][j] = moon.y
			axes[2][j] = moon.z
		}

		// fmt.Println("c:", axes)
		// fmt.Println("i:", axesInitial)
		// fmt.Println("p:", periods)

		for i := range axes {
			if equal(axes[i], axesInitial[i]) {
				periods[i] = iter + 1
				fmt.Println(axes[i])
				fmt.Println(axesInitial[i])
			}
		}

		if periods[0] != 0 && periods[1] != 0 && periods[2] != 0 {
			break
		}
	}

	fmt.Println(periods)
	fmt.Println(lcm(periods[0], periods[1], periods...))

}

func equal(a, b []axis) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// getAxis(moons []moon, axis string) {
// 	axisData := []
// 	if axis == "x" {

// 		return
// 	}
// }

func stepVelocity(moons []moon) []moon {
	for i := range moons {
		for _, moon := range moons {
			moons[i].x.vel += velocityCompare(moons[i].x.pos, moon.x.pos)
			moons[i].y.vel += velocityCompare(moons[i].y.pos, moon.y.pos)
			moons[i].z.vel += velocityCompare(moons[i].z.pos, moon.z.pos)
		}
	}

	return moons
}

func stepPosition(moons []moon) []moon {
	for i := range moons {
		//TODO define + operator on struct
		moons[i].x.pos = moons[i].x.pos + moons[i].x.vel
		moons[i].y.pos = moons[i].y.pos + moons[i].y.vel
		moons[i].z.pos = moons[i].z.pos + moons[i].z.vel
	}

	return moons
}

func calcEnergy(target moon) int {
	return int(
		(math.Abs(float64(target.x.pos)) +
			math.Abs(float64(target.y.pos)) +
			math.Abs(float64(target.z.pos))) *
			(math.Abs(float64(target.x.vel)) +
				math.Abs(float64(target.y.vel)) +
				math.Abs(float64(target.z.vel))))
}

func velocityCompare(a, b int) int {
	if a < b {
		return 1
	} else if a > b {
		return -1
	} else {
		return 0
	}
}

func printMoons(moons []moon) {
	for _, moon := range moons {
		fmt.Println(moon.x, moon.y, moon.z)
	}
}

// greatest common divisor (GCD) via Euclidean algorithm
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
