package keks

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
)

var path = flag.String("path", "data/data.csv", "path to dataset")
var firstTown = flag.Int("from", 0, "from town")
var lastTown = flag.Int("to", 9, "to town")

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

type Genome interface {
	Fitness() float64
}
type Population []Genome

func (pop Population) Len() int {
	return len(pop)
}
func (pop Population) Less(i, j int) bool {
	return pop[i].Fitness() < pop[j].Fitness()
}
func (pop Population) Swap(i, j int) {
	pop[i], pop[j] = pop[j], pop[i]
}

type Generator func() Genome
type Mutator func(Genome) Genome
type Childor func(Genome, Genome) Genome

func Mutation(population []Genome, mutator Mutator, awesome int) {
	for _, v := range population[awesome:] {
		if rand.Float64() > 0.4 {
			v = mutator(v)
		}

	}
}

func Selection([]Genome) {

}

func NewGeneration([]Genome) {

}

func Result() int {
	return 0
}

func GenericAlg(random_gen Generator, mutate Mutator, awesome_number, population_size, epochs, verbose int) {
	generation := make([]Genome, population_size)
	for i := 0; i < population_size; i++ {
		generation[i] = random_gen()
	}
	sort.Sort(Population(generation))
	for i := 0; i < epochs; i++ {
		Mutation(generation, mutate, awesome_number)
		Selection(generation)
		//NewGeneration(generation)
		if verbose > 0 {
			fmt.Printf("Эпоха %d: Лучший результат = %d", i+1, Result())
		}
	}
}

func mutate(gen Genome) Genome{
	return gen
}

type TownPath struct {
	first_town *int
	last_town  *int
	townMap    *[][]int
	townNums   []int
}

func (tp TownPath) Fitness() float64 {
	sum := 0
	last_val := *tp.first_town
	tmap := *(tp.townMap)
	for _, val := range tp.townNums {
		sum += tmap[last_val][val]
	}
	sum += tmap[last_val][*tp.last_town]
	return float64(sum)
}

func RandomPath(first_town, last_town, size int, townMap *[][]int) Generator {
	return func() Genome {
		townNums := make([]int, size)
		for i := 0; i < size; i++ {
			townNums[i] = int(rand.Float64() * float64(size))
		}
		return TownPath{
			first_town: &first_town,
			last_town:  &last_town,
			townNums:   townNums,
			townMap:    townMap,
		}
	}
}

func main() {
	flag.Parse()
	townMap := getData(*path)
	GenericAlg(RandomPath(*firstTown, *lastTown, len(townMap[0]), &townMap), mutate, 0, 20, 1, 1)
}
