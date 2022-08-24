package main

import (
	"fmt"
	"testTrains/costSolver"
)

func main() {
	path, cost := costSolver.Solver()
	fmt.Println(path)
	fmt.Println(cost)
}
