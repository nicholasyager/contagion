package agents

import (
	"fmt"
	"math/rand"
)

type Person struct {
	Id             int
	Status         int     // The infection Status for the individual.
	StatusClock    int     // The clock for handling incubation times.
	IncubationTime int     // The amount of time it takes to incubate the infection.
	InfectiousTime int     // The amount of time it takes for recover from the infecton.
	Virality       float64 // The probability of getting infected if exposed.
	Heading        int     // The positional heading for the person.
	X              int
	Y              int
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
		if newHeading < 0 {
			newHeading += 8
		} else if newHeading > 7 {
			newHeading -= 8
		}

		newX := movementArray[newHeading][0] + person.X
		newY := movementArray[newHeading][1] + person.Y

		if newX >= width {
			newX -= width
		} else if newX < 0 {
			newX += width
		}

		if newY >= height {
			newY -= height
		} else if newY < 0 {
			newY += height
		}

		if matrix[newX][newY] == 0 {
			person.Heading = newHeading
			person.X = newX
			person.Y = newY
			break
		}

	}

	if person.Status == 1 && person.StatusClock < person.IncubationTime {
		person.StatusClock++
	} else if person.Status == 2 && person.StatusClock < person.InfectiousTime {
		person.StatusClock++
	} else if person.Status == 2 {
		person.Status = 3
		person.StatusClock = 0
	} else if person.Status == 1 {
		person.Status = 2
		person.StatusClock = 0
	}

}

func (person *Person) CheckInfection(matrix [][]int) {

	if person.Status > 0 {
		return
	}

	movementArray := [][]int{[]int{-1, -1}, []int{-1, 0}, []int{-1, 1},
		[]int{0, 1}, []int{1, 1}, []int{1, 0},
		[]int{1, -1}, []int{0, -1}}

	for i := 0; i < 8; i++ {
		neighborX := person.X + movementArray[i][0]
		neighborY := person.Y + movementArray[i][1]

		if neighborX >= len(matrix[0]) {
			neighborX -= len(matrix[0])
		} else if neighborX < 0 {
			neighborX += len(matrix[0])
		}

		if neighborY >= len(matrix) {
			neighborY -= len(matrix)
		} else if neighborY < 0 {
			neighborY += len(matrix)
		}

		if matrix[neighborX][neighborY] == 3 && rand.Float64() <= person.Virality {
			person.Status = 1
		}
	}
}
