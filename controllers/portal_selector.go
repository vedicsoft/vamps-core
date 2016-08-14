package controllers

import (
	"sort"
	"strconv"
)

type Tuple struct {
	key    string
	value  string
	weight float64
}

type Bucket struct {
	id     string
	tuples []Tuple
}

type Pair struct {
	Key   string
	Value float64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p PairList) Len() int {
	return len(p)
}
func (p PairList) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]float64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

func main() {
	referenceTuple1 :=
		Bucket{
			id: "abc",
			tuples: []Tuple{
				Tuple{key: "location", value: "A", weight: 100},
				Tuple{key: "age", value: "55", weight: 75},
				Tuple{key: "gender", value: "m", weight: 75},
			},
		}

	referenceTuple2 :=
		Bucket{
			id: "pqr",
			tuples: []Tuple{
				Tuple{key: "location", value: "B", weight: 100},
				Tuple{key: "age", value: "45", weight: 75},
				Tuple{key: "gender", value: "f", weight: 25},
			},
		}

	buckets := []Bucket{
		referenceTuple1,
		referenceTuple2,
	}

	inputTuples := []Tuple{
		Tuple{key: "location", value: "A"},
		Tuple{key: "age", value: "45"},
		Tuple{key: "gender", value: "f"},
	}

	scores := make(map[string]float64, len(buckets))
	for _, bucket := range buckets {
		scores[bucket.id] = getWeightedAverage(bucket, inputTuples)
	}

	sortMapByValue(scores)

	for k, v := range scores {
		println(k + " " + strconv.FormatFloat(v, 'f', 6, 64))
	}

}

func getWeightedAverage(bucket Bucket, inputTuples []Tuple) float64 {
	var score float64
	var devider float64
	for _, tuple := range bucket.tuples {
		for _, inputTuple := range inputTuples {
			if inputTuple.key == tuple.key {
				if inputTuple.value == tuple.value {
					score += tuple.weight
				}
			}
		}
		devider += tuple.weight
	}
	score /= devider
	return score
}
