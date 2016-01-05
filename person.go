package main

import (
	"fmt"
	"math/rand"
)

type Person struct {
	Id      int
	Disease Disease // The disease this person is infected with.
	Status  int     // The infection Status for the individual.
	Clock   int     // The clock for handling incubation times.
	Heading int     // The positional heading for the person.
	X       int
	Y       int
}

func boundryCheck(n, max int) int {
	if n >= max {
		n -= max
	} else if n < 0 {
		n += max
	}
	return n
}

func (person *Person) Print() {
	fmt.Printf("%d, %d, %d, %d, %d\n", person.Id, person.Status, person.Heading, person.X, person.Y)

}

func (person *Person) Update(matrix [][]int) {

	width := len(matrix[0])
	height := len(matrix)

	movementArray := [][]int{[]int{-1, -1}, []int{-1, 0}, []int{-1, 1},
		[]int{0, 1}, []int{1, 1}, []int{1, 0},
		[]int{1, -1}, []int{0, -1}}

	list := rand.Perm(3)
	for i, _ := range list {
		newHeading := person.Heading + (list[i] - 1)
		newHeading = boundryCheck(newHeading, 8)

		newX := movementArray[newHeading][0] + person.X
		newY := movementArray[newHeading][1] + person.Y

		newX = boundryCheck(newX, width)
		newY = boundryCheck(newY, height)

		if matrix[newX][newY] == 0 {
			person.Heading = newHeading
			person.X = newX
			person.Y = newY
			break
		}

	}

	if person.Status == 1 && person.Clock < person.Disease.Timer[person.Status] {
		person.Clock++
	} else if person.Status == 2 && person.Clock < person.Disease.Timer[person.Status] {
		person.Clock++
	} else if person.Status == 2 {
		person.Status = 3
		person.Clock = 0
	} else if person.Status == 1 {
		person.Status = 2
		person.Clock = 0
	}

}

func (person *Person) CheckInfection(matrix [][]int) {

	if person.Status > 0 {
		return
	}

	width := len(matrix[0])
	height := len(matrix)

	movementArray := [][]int{[]int{-1, -1}, []int{-1, 0}, []int{-1, 1},
		[]int{0, 1}, []int{1, 1}, []int{1, 0},
		[]int{1, -1}, []int{0, -1}}

	for i := 0; i < 8; i++ {
		neighborX := person.X + movementArray[i][0]
		neighborY := person.Y + movementArray[i][1]

		neighborX = boundryCheck(neighborX, width)
		neighborY = boundryCheck(neighborY, height)

		if matrix[neighborX][neighborY] == 3 && rand.Float64() <= person.Disease.Virality {
			person.Status = 1
		}
	}
}
