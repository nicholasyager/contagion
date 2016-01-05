package main

import (
	"flag"
	"fmt"
	"github.com/cheggaaa/pb"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"os"
	"runtime"
)

var black = color.RGBA{0x00, 0x00, 0x00, 0xff}
var blue = color.RGBA{0x00, 0x00, 0xff, 0xff}
var red = color.RGBA{0xff, 0x00, 0x00, 0xff}
var green = color.RGBA{0x00, 0xff, 0x00, 0xff}
var grey = color.RGBA{0x33, 0x33, 0x33, 0xff}

var populationStatistics = []int{0, 0, 0, 0}
var states = []string{"Succeptable", "Exposed", "Infectious", "Removed"}

func generateMatrix(width, height int) [][]int {
	// Allocate the top-level slice.
	matrix := make([][]int, width) // One row per unit of y.
	// Loop over the rows, allocating the slice for each row.
	for i := range matrix {
		matrix[i] = make([]int, height)
	}
	return matrix
}

func render(RGBAImage *image.RGBA, matrix [][]int, time int) {

	for x := 0; x < len(matrix[0]); x++ {
		for y := 0; y < len(matrix); y++ {
			if matrix[x][y] == 1 {
				RGBAImage.Set(x, y, red)
			} else if matrix[x][y] == 2 {
				RGBAImage.Set(x, y, blue)
			} else if matrix[x][y] == 3 {
				RGBAImage.Set(x, y, green)
			} else if matrix[x][y] == 4 {
				RGBAImage.Set(x, y, grey)

			} else {
				RGBAImage.Set(x, y, black)
			}
		}
	}

	imgFile, _ := os.Create(fmt.Sprintf("images/%06d.png", time))

	png.Encode(imgFile, RGBAImage)
	imgFile.Close()

}

func generatePeople(numPeople int, width, height int, disease Disease) []Person {
	people := make([]Person, numPeople)
	for i, person := range people {
		person.Id = i
		person.Status = 0
		person.Heading = rand.Intn(8)
		person.Disease = disease
		person.Clock = 0
		person.X = rand.Intn(width)
		person.Y = rand.Intn(height)
		people[i] = person
	}
	return people
}

func updatePeople(people []Person, matrix [][]int, c chan int) {
	for i, person := range people {
		matrix[person.X][person.Y] = 0
		// Update the location of each person
		person.Update(matrix)
		people[i] = person
	}
	c <- 1
}

func updateInfections(people []Person, matrix [][]int, c chan int) {
	for i, person := range people {
		// Check for an infection
		person.CheckInfection(matrix)
		matrix[person.X][person.Y] = person.Status + 1

		people[i] = person
	}
	c <- 1
}

var maxTime = flag.Int("time", 3000, "The number of steps to simulate.")
var width = flag.Int("width", 500, "The width of the world.")
var height = flag.Int("height", 500, "The height of the world.")
var density = flag.Float64("density", 0.1, "The density of the agents in the world.")
var virality = flag.Float64("virality", 0.5, "The proportion of individuals that are infected when exposed.")

func main() {

	numCPU := runtime.NumCPU()
	flag.Parse()

	numPeople := int(float64(*width**height) * *density)

	fmt.Printf("Duration: %d ticks.\n", *maxTime)
	fmt.Printf("World Size: %dx%d.\n", *width, *height)
	fmt.Printf("Population: %d.\n", numPeople)
	fmt.Printf("Density: %f.\n", *density)
	fmt.Printf("Using %d CPU cores.\n", numCPU)
	fmt.Printf("Simulation started.\n")

	os.Mkdir("images", 0777)
	file, _ := os.Create("simulation.csv")

	// Create a new progress bar.
	bar := pb.StartNew(*maxTime)

	disease := NewDisease(*virality, []int{0, 150, 50, 0}, SEIRMatrix)

	people := generatePeople(numPeople, *width, *height, *disease)
	people[rand.Intn(numPeople)].Status = 1

	matrix := generateMatrix(*width, *height)

	rectangle := image.Rect(0, 0, len(matrix[0]), len(matrix))
	RGBAImage := image.NewRGBA(rectangle)

	for time := 0; time < *maxTime; time++ {

		for i, _ := range populationStatistics {
			populationStatistics[i] = 0
		}

		c := make(chan int, numCPU) // Buffering optional but sensible.
		for i := 0; i < numCPU; i++ {
			go updatePeople(people[i*(numPeople/numCPU):(i+1)*(numPeople/numCPU)], matrix, c)
		}
		// Drain the channel.
		for i := 0; i < numCPU; i++ {
			<-c // wait for one task to complete
		}
		for i := 0; i < numCPU; i++ {
			go updateInfections(people[i*(numPeople/numCPU):(i+1)*(numPeople/numCPU)], matrix, c)
		}
		// Drain the channel.
		for i := 0; i < numCPU; i++ {
			<-c // wait for one task to complete
		}

		for _, person := range people {
			populationStatistics[person.Status]++
		}

		render(RGBAImage, matrix, time)
		bar.Increment()

		specialPop := 0

		for i := 0; i < 4; i++ {
			io.WriteString(file, fmt.Sprintf("%d, %s, %d\n", time, states[i],
				populationStatistics[i]))
		}

		for i := 1; i < 3; i++ {
			specialPop += populationStatistics[i]

		}
		if specialPop == 0 {
			fmt.Println("No infectious or exposed individuals remaining.")
			break
		}
	}
	file.Close()
	bar.FinishPrint("Simulation complete!")
}
