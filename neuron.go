package main

import (
	"fmt"
	"math/rand"
	"sort"
	//	"strings"
	"math"
	"os"
	"encoding/csv"
	"bufio"
	"io"
	"strconv"
)

type NeuroPopulation []Neuron

func (pop NeuroPopulation) Len() int {
	return len(pop)
}
func (pop NeuroPopulation) Less(i, j int) bool {
	f1 := pop[i].Fitness()
	f2 := pop[j].Fitness()
	//fmt.Println(f1, f2)
	return f1 < f2
}
func (pop NeuroPopulation) Swap(i, j int) {
	pop[i], pop[j] = pop[j], pop[i]
}

func Sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func GenericAlgNeuron(data_x [][]float64, data_y []int, ind int, population_size, epochs, step int, mutation_prob float64, verbose int) {
	generation := make([]Neuron, population_size)
	all_generations := []Neuron{}
	//all_generations_map := make(map[string]struct{})
	for i := 0; i < population_size; i++ {
		generation[i] = NewRandomNeuron(len(data_x[0]), ind, Sigmoid, data_x, data_y)
		//if _, ok := all_generations_map[fmt.Sprint(generation[i])]; !ok {
		all_generations = append(all_generations, generation[i])
		//	all_generations_map[fmt.Sprint(generation[i].weights)] = struct{}{}
		//}
	}

	for i := 0; i < epochs; i++ {
		//sort.Sort(NeuroPopulation(generation))
		//size := len(all_generations)
		all_generations := []Neuron{}
		for i := 0; i <= population_size/2; i++ {
			generation_crossover := make([]Neuron, step)
			for j := 0; j < step; j++ {
				generation_crossover[j] = generation[int(rand.Float64()*float64(len(generation)))]
			}
			sort.Sort(NeuroPopulation(generation_crossover))
			new_gen1, new_gen2 := generation_crossover[0].crossover(generation_crossover[1])
			//if _, ok := all_generations_map[fmt.Sprint(new_gen1.weights)]; !ok {
			all_generations = append(all_generations, new_gen1)
			//	all_generations_map[fmt.Sprint(new_gen1.weights)] = struct{}{}
			//}
			//if _, ok := all_generations_map[fmt.Sprint(new_gen2.weights)]; !ok {
			all_generations = append(all_generations, new_gen2)
			//	all_generations_map[fmt.Sprint(new_gen2.weights)] = struct{}{}
			//}
			//all_generations = append(all_generations, new_gen1, new_gen2)
		}
		//fmt.Print(new_generation[0])
		for i := range all_generations {
			if (rand.Float64() < mutation_prob) {
				gen := all_generations[i].mutate()
				//if _, ok := all_generations_map[fmt.Sprint(gen.weights)]; !ok {
				all_generations = append(all_generations, gen)
				//	all_generations_map[fmt.Sprint(gen.weights)] = struct{}{}
				//}
			}
		}
		//for i := 0; i < awesome_number; i++ {
		//	all_generations = append(all_generations, generation[i])
		//}
		fmt.Println(len(all_generations))

		sort.Sort(NeuroPopulation(all_generations))
		copy(generation, all_generations)
		//sort.Sort(Population(generation))
		fit := generation[0].Fitness()
		if verbose == 1 {
			fmt.Printf("Эпоха %d: Лучший результат = %.2f; Совпавших: %d (%.2f%%)\n", i+1, fit, generation[0].count_valid, float64(generation[0].count_valid)/float64(len(data_x)) * 100)
		} else if verbose == 2 {
			fmt.Printf("Особи эпохи %d: %v\n", i+1, generation)
		}
		//generation
	}
	//fmt.Println(all_generations_map)
}

type Activation func(x float64) float64

type Neuron struct {
	input_size int
	ind        int
	weights    []float64
	activation Activation
	train_x    [][]float64
	train_y    []int

	count_valid int
}

func NewRandomNeuron(input_size, ind int, activate Activation, train_x [][]float64, train_y []int) (Neuron) {
	weights := make([]float64, input_size+1)
	for i := range weights {
		weights[i] = rand.Float64() * 1 * float64(rand.Intn(2) - 1)
	}
	return Neuron{input_size, ind, weights, activate, train_x, train_y, 0}
}

func (n *Neuron) getResult(data []float64) float64 {
	result := -n.weights[n.input_size]
	for i := 0; i < n.input_size; i++ {
		result += data[i] * n.weights[i]
	}
	r := n.activation(result)
	//fmt.Println(r)
	return r
}

func (n Neuron) mutate() Neuron {
	index := int(rand.Float64() * float64(len(n.weights)))
	//n.weights[index] *= rand.Float64() * 2 * float64(rand.Intn(2) - 1)
	n.weights[index] *= 1.5
	//n.weights[index + 1] *= 1.5
	return n
}

func (n1 Neuron) crossover(n2 Neuron) (n3 Neuron, n4 Neuron) {
	index := int(rand.Float64() * float64(len(n1.weights)))
	n3 = Neuron{n1.input_size, n1.ind, make([]float64, len(n1.weights)), n1.activation, n1.train_x, n1.train_y, 0}
	n4 = Neuron{n2.input_size, n2.ind, make([]float64, len(n2.weights)), n2.activation, n2.train_x, n2.train_y, 0}
	copy(n3.weights, n1.weights)
	copy(n4.weights, n2.weights)
	for i := 0; i < index; i++ {
		n3.weights[i] = n2.weights[i]
		n4.weights[len(n4.weights)-1-i] = n1.weights[len(n1.weights)-1-i]
	}
	//fmt.Println(tp3, tp4)
	return
}

func (n *Neuron) Fitness() float64 {
	sum := 0.
	n.count_valid = 0
	count_invalid := 0

	for i, value := range n.train_x {
		expected := 0.
		if(n.ind == n.train_y[i]) {
			expected = 1.
		}
		prob := n.getResult(value)
		if (prob >= 0.5 && expected == 1.) || (prob < 0.5 && expected == 0.) {
			n.count_valid += 1
		} else {
			count_invalid++
			//fmt.Println("Prob:", prob, n.train_y[i], expected)
		}
		sum += math.Abs((expected - prob))
	}
	//fmt.Println(count_invalid)
	//fmt.Println(n.count_valid)
	return float64(sum)
}

const MAX_DARKNESS = 255

func getTrainData(path string, header bool) (train_x [][]float64, train_y []int) {
	f, _ := os.Open(path)
	r := csv.NewReader(bufio.NewReader(f))
	if header {
		r.Read()
	}
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		value, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Println("Not valid csv file")
			os.Exit(1)
		}
		image := record[1:]
		image_normalized := make([]float64, len(image))
		for ind, pix := range image {
			pix, err := strconv.Atoi(pix)
			if err != nil {
				fmt.Println("Not valid csv file")
				os.Exit(1)
			}
			image_normalized[ind] = float64(pix) / MAX_DARKNESS
		}
		train_y = append(train_y, value)
		train_x = append(train_x, image_normalized)
	}
	return
}

func main() {
	train_x, train_y := getTrainData("data/train.csv", true)

	//fmt.Print(train_y)
	GenericAlgNeuron(
		train_x[:100],
		train_y[:100],
		1,
		250,
		10000,
		5,
		0.3,
		1,
	)
}
