package timeSolver

import (
	"encoding/csv"
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
	TrainID   int
	StartTime int
	StopTime  int
}

func stringToTimestamp(s string) int {
	hourS := s[:2]
	minuteS := s[3:5]
	hour, _ := strconv.Atoi(hourS)
	minute, _ := strconv.Atoi(minuteS)
	return hour*60 + minute
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
		trainData.StartTime = stringToTimestamp(line[4])
		trainData.StopTime = stringToTimestamp(line[5])
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

func MakeAdjacencyList(data map[VertexPair][]TrainCostPair) map[int][]int {
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

func travelTime(curTime int, train TrainCostPair) int {
	if curTime < train.StartTime { //We can take this train today
		if train.StopTime < train.StartTime { //The train comes the next day
			return (train.StartTime - curTime) + (24*60 - train.StartTime) + train.StopTime
		}
		return (train.StartTime - curTime) + (train.StopTime - train.StartTime)
	} else { //We'll have to wait until the next day to get on that train.
		if train.StopTime < train.StartTime {
			return (24*60 - curTime) + 24*60 + train.StopTime
		}
		return (24*60 - curTime) + train.StopTime
	}

}

func Solver() ([]int, int) {
	data, _ := ReadCSV()
	adjList := MakeAdjacencyList(data)

	var curTime = 0
	var OverallTime = 0
	var TakenTrains []int
	isVertexVisited := make(map[int]bool)

	for k := range adjList {
		isVertexVisited[k] = false
	}

	isVertexVisited[startVertex] = true
	curVertex := startVertex

	for isThereUnvisitedVertex(isVertexVisited) {
		adjacentVertices := adjList[curVertex]
		minTime := 1 << 32
		var nextVertex int
		var trainID int
		var timeInTravel = 0

		for _, tryVertex := range adjacentVertices {
			var vPair = VertexPair{curVertex, tryVertex}
			if isVertexVisited[tryVertex] == false {
				for _, train := range data[vPair] {
					timeInTravel = travelTime(curTime, train)
					if minTime > timeInTravel {
						nextVertex = tryVertex
						minTime = timeInTravel
						trainID = train.TrainID
					}
				}
			}
		}

		curTime = (curTime + timeInTravel) % (24 * 60)
		curVertex = nextVertex
		OverallTime += timeInTravel
		isVertexVisited[curVertex] = true
		TakenTrains = append(TakenTrains, trainID)
	}

	return TakenTrains, OverallTime
}
