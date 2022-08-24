package main

import (
    "encoding/csv"
    "fmt"
    "math"
    "os"
    "strconv"
)

var inputName = "test_task_data.csv"
var startVertex = 1909

type VertexPair struct {
    OutVertex int
    InVertex  int
}

type TrainCostPair struct {
    TrainID int
    Cost    float64
}

func ReadCSV() (map[VertexPair][]TrainCostPair, error) {
    inputMap := make(map[VertexPair][]TrainCostPair)

    f, err := os.Open(inputName)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    reader := csv.NewReader(f)
    line, _ := reader.Read()
    for line != nil {
        var trainData TrainCostPair
        trainData.TrainID, _ = strconv.Atoi(line[0])
        trainData.Cost, _ = strconv.ParseFloat(line[3], 64)
        outV, _ := strconv.Atoi(line[1])
        inV, _ := strconv.Atoi(line[2])
        vPair := VertexPair{outV, inV}
        if val, ok := inputMap[vPair]; ok {
            val = append(val, trainData)
            inputMap[vPair] = val
        } else {
            val = make([]TrainCostPair, 1)
            val[0] = trainData
            inputMap[vPair] = val
        }
        line, _ = reader.Read()
    }
    return inputMap, nil
}

func bubbleSort(inputMap map[VertexPair][]TrainCostPair, pair VertexPair) {
    input := inputMap[pair]
    for i := 0; i < len(input); i++ {
        for j := 0; j < len(input)-1; j++ {
            if input[j].Cost > input[j+1].Cost {
                input[j], input[j+1] = input[j+1], input[j]
            }
        }
    }
    inputMap[pair] = input
}

func DeleteTrains(input map[VertexPair][]TrainCostPair) map[VertexPair]TrainCostPair {
    output := make(map[VertexPair]TrainCostPair)
    for key := range input {
        bubbleSort(input, key)
        output[key] = input[key][0]
    }
    return output
}

func makeAdjacencyList(data map[VertexPair]TrainCostPair) map[int][]int {
    output := make(map[int][]int)
    for key := range data {
        if val, ok := output[key.OutVertex]; ok {
            val = append(val, key.InVertex)
            output[key.OutVertex] = val
        } else {
            val = make([]int, 1)
            val[0] = key.InVertex
            output[key.OutVertex] = val
        }
    }
    return output
}

func isThereUnvisitedVertex(input map[int]bool) bool {
    for _, v := range input {
        if v == false {
            return true
        }
    }
    return false
}

func Solver() ([]int, float64) {
    data, _ := ReadCSV()
    graph := DeleteTrains(data)
    adjList := makeAdjacencyList(graph)

    var PathCost float64
    var TakenTrains []int
    isVertexVisited := make(map[int]bool)

    for k := range adjList {
        isVertexVisited[k] = false
    }

    isVertexVisited[startVertex] = true
    curVertex := startVertex

    for isThereUnvisitedVertex(isVertexVisited) {
        adjacentVertices := adjList[curVertex]
        minCost := math.Inf(1)
        var nextVertex int
        var trainID int

        for _, tryVertex := range adjacentVertices {
            var vPair = VertexPair{curVertex, tryVertex}
            if minCost > graph[vPair].Cost {
                if isVertexVisited[tryVertex] == false {
                    nextVertex = tryVertex
                    minCost = graph[vPair].Cost
                    trainID = graph[vPair].TrainID
                }
            }
        }

        curVertex = nextVertex
        PathCost += minCost
        isVertexVisited[curVertex] = true
        TakenTrains = append(TakenTrains, trainID)
    }

    return TakenTrains, PathCost
}

func main() {
    path, cost := Solver()
    fmt.Println(path)
    fmt.Println(cost)
}
