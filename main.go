package main

import (
	"fmt"
	"testTrains/costSolver"
	"testTrains/timeSolver"
)

func main() {
	path, time := timeSolver.Solver()
	fmt.Println(path)
	fmt.Println(time)

	path, price := costSolver.Solver()
	fmt.Println(path)
	fmt.Println(price)

}
