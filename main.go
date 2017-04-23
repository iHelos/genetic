package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strconv"
//	"strings"
)

var path = flag.String("path", "data/data.csv", "path to dataset")
var firstTown = flag.Int("from", 4, "from town")
var lastTown = flag.Int("to", 7, "to town")

func getData(path string) (townmap [][]int) {
	f, _ := os.Open(path)
	r := csv.NewReader(bufio.NewReader(f))
	for {
		temprow := []int{}
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		row := record[0:]
		for _, val := range row {
			length, err := strconv.Atoi(val)
			if err != nil {
				fmt.Println("Not valid csv file")
				os.Exit(1)
			}
			temprow = append(temprow, length)
		}
		townmap = append(townmap, temprow)

	}
	return
}

type Population []TownPath

func (pop Population) Len() int {
	return len(pop)
}
func (pop Population) Less(i, j int) bool {
	return pop[i].Fitness() < pop[j].Fitness()
}
func (pop Population) Swap(i, j int) {
	pop[i], pop[j] = pop[j], pop[i]
}

type Generator func() TownPath

func GenericAlg(random_gen Generator, awesome_number, population_size, epochs, step int, mutation_prob float64, verbose int) {
	generation := make([]TownPath, population_size)
	all_generations := []TownPath{}
	all_generations_map := make(map[string]struct{})
	for i := 0; i < population_size; i++ {
		generation[i] = random_gen()
		if _, ok := all_generations_map[fmt.Sprint(generation[i])]; !ok {
			all_generations = append(all_generations, generation[i])
			all_generations_map[fmt.Sprint(generation[i].path)] = struct{}{}
		}
	}

	for i := 0; i < epochs; i++ {
		sort.Sort(Population(generation))
		//size := len(all_generations)
		for i := 0; i <= population_size / 2; i++ {
			generation_crossover := make([]TownPath, step)
			for j := 0; j < step; j++ {
				generation_crossover[j] = all_generations[int(rand.Float64()*float64(len(all_generations)))]
			}
			sort.Sort(Population(generation_crossover))
			new_gen1, new_gen2 := generation_crossover[0].crossover(&generation_crossover[1])
			if _, ok := all_generations_map[fmt.Sprint(new_gen1.path)]; !ok {
				all_generations = append(all_generations, new_gen1)
				all_generations_map[fmt.Sprint(new_gen1.path)] = struct{}{}
			}
			if _, ok := all_generations_map[fmt.Sprint(new_gen2.path)]; !ok {
				all_generations = append(all_generations, new_gen2)
				all_generations_map[fmt.Sprint(new_gen2.path)] = struct{}{}
			}
			//all_generations = append(all_generations, new_gen1, new_gen2)
		}
		//fmt.Print(new_generation[0])
		for i := range all_generations {
			if (rand.Float64() < mutation_prob) {
				gen := all_generations[i].mutate()
				if _, ok := all_generations_map[fmt.Sprint(gen.path)]; !ok {
					all_generations = append(all_generations, gen)
					all_generations_map[fmt.Sprint(gen.path)] = struct{}{}
				}
			}
		}
		//for i := 0; i < awesome_number; i++ {
		//	all_generations = append(all_generations, generation[i])
		//}
		fmt.Println(len(all_generations))

		sort.Sort(Population(all_generations))
		copy(generation, all_generations)
		//sort.Sort(Population(generation))
		if verbose == 1 {
			fmt.Printf("Эпоха %d: Лучший результат = %.2f %v\n", i+1, generation[0].Fitness(), generation[0])
		} else if verbose == 2 {
			fmt.Printf("Особи эпохи %d: %v\n", i+1, generation)
		}
		//generation
	}
	fmt.Println(all_generations_map)
}

func (tp TownPath) mutate() TownPath {

	index := int(rand.Float64() * float64(len(tp.path) - 1))
	//fmt.Print(index)
	tp.path[index] = tp.path[index + 1]

	return tp
}

func (tp1 *TownPath) crossover(tp2 *TownPath) (tp3 TownPath, tp4 TownPath) {
	index := int(rand.Float64() * float64(len(tp1.path)))
	tp3 = TownPath{
		townMap:   tp1.townMap,
		firstTown: tp1.firstTown,
		lastTown:  tp2.lastTown,
		path:      make([]int, len(tp1.path)),
	}
	tp4 = TownPath{
		townMap:   tp1.townMap,
		firstTown: tp1.firstTown,
		lastTown:  tp2.lastTown,
		path:      make([]int, len(tp1.path)),
	}
	copy(tp3.path, tp1.path)
	copy(tp4.path, tp2.path)
	for i := 0; i < index; i++ {
		tp3.path[i] = tp2.path[i]
		tp4.path[len(tp4.path)-1-i] = tp1.path[len(tp1.path)-1-i]
	}
	//fmt.Println(tp3, tp4)
	return
}

func (tp *TownPath) randomTown() int {
	size := len(tp.path)
	temp_town := int(rand.Float64() * float64(size))
	for temp_town == tp.lastTown {
		temp_town = int(rand.Float64() * float64(size))
	}
	return temp_town
}

type TownPath struct {
	townMap   *[][]int
	firstTown int
	lastTown  int
	path      []int
}

func (tp TownPath) Fitness() float64 {
	sum := 0
	tmap := *(tp.townMap)
	last_val := tp.firstTown
	for _, val := range tp.path {
		sum += tmap[last_val][val]
		last_val = val
	}
	sum += tmap[last_val][tp.lastTown]
	return float64(sum)
}

func RandomPath(first_town, last_town, size int, townMap *[][]int) Generator {
	return func() TownPath {
		townNums := []int{}
		for i := 0; i < size; i++ {
			temp_town := int(rand.Float64() * float64(size))
			for temp_town == last_town {
				temp_town = int(rand.Float64() * float64(size))
			}
			townNums = append(townNums, temp_town)
		}
		return TownPath{
			firstTown: first_town,
			lastTown:  last_town,
			path:      townNums,
			townMap:   townMap,
		}
	}
}

func main() {
	flag.Parse()
	townMap := getData(*path)
	GenericAlg(
		RandomPath(*firstTown, *lastTown, len(townMap[0]), &townMap),
		4,
		200,
		500,
		5,
		0.1,
		1,
	)
}
