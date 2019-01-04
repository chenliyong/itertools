package itertools

import (
	"math/rand"
	"sort"
	"time"
)

type GA struct {
	mutationProb, elite          float32
	size, maxIter, maxIterRepeat int
	// 种群
	population [][]int
	// 适应度计算函数
	fitness func([]int) float32
}

func NewGA(length, popSize int, mutationProb, elite float32, maxIter, maxIterRepeat int, fitness func([]int) float32) *GA {
	rand.Seed(time.Now().Unix())
	population := make([][]int, popSize)
	for i := 0; i < popSize; i++ {
		population[i] = RandomSample(Range(length))
	}
	return &GA{
		size:          popSize,
		mutationProb:  mutationProb,
		elite:         elite,
		maxIter:       maxIter,
		maxIterRepeat: maxIterRepeat,
		population:    population,
		fitness:       fitness,
	}
}

func (g *GA) Run() {
	var n, iterRepeat int
	for {
		if n > g.maxIter {
			break
		}
		// 进行一次迭代
		newPop := g.bread()
		//fmt.Println("newPop", g.fitness(g.population[0]))
		if g.fitness(g.population[0]) == g.fitness(newPop[0]) {
			iterRepeat = iterRepeat + 1
			if iterRepeat >= g.maxIterRepeat {
				break
			}
		} else {
			iterRepeat = 0
		}
		// 改变当前种群
		copy(g.population, newPop)
		n = n + 1
	}
}

func (g *GA) Best() ([]int, float32) {
	result := g.population[0]
	rank := g.fitness(result)
	return result, rank
}

func (g *GA) ranked() []Ranked {
	popRanked := make([]Ranked, g.size)
	// 排序种群
	for i := 0; i < g.size; i++ {
		popRanked[i] = Ranked{
			Index: i,
			Rank:  g.fitness(g.population[i]),
		}
	}
	sort.Slice(popRanked, func(i, j int) bool {
		return popRanked[i].Rank > popRanked[j].Rank
	})
	return popRanked
}

// 适者生存
func (g *GA) selection() [][]int {
	var pick, cum_sum float32

	popRanked := g.ranked()
	results := [][]int{}
	sumFits := SumRankedFits(popRanked)

	//  精英选择
	eliteSize := int(float32(g.size) * g.elite)
	for i := 0; i < eliteSize; i++ {
		results = append(results, g.population[popRanked[i].Index])
	}

	// 轮盘选择
	for i := 0; i < len(popRanked)-eliteSize; i++ {
		pick = 100 * rand.Float32()
		cum_sum = 0
		for j := 0; j < len(popRanked); j++ {
			cum_sum = cum_sum + popRanked[j].Rank
			if pick <= 100*cum_sum/sumFits {
				results = append(results, g.population[popRanked[j].Index])
				break
			}
		}
	}
	return results
}

func (g *GA) bread() [][]int {
	children := [][]int{}

	pool := g.selection()

	// 精英直接进入下一种群
	eliteSize := int(float32(g.size) * g.elite)
	for i := 0; i < eliteSize; i++ {
		children = append(children, pool[i])
	}

	// 乱序
	pool = g.random(pool)
	for i := 0; i < len(pool)-eliteSize; i++ {
		c := g.crossover(pool[i], pool[len(pool)-1-i])
		c = g.mutate(c)
		children = append(children, c)
	}

	return children

}

// 乱序处理
func (g *GA) random(slice [][]int) [][]int {
	for i := len(slice) - 1; i > 0; i-- {
		n := rand.Intn(i + 1)
		slice[i], slice[n] = slice[n], slice[i]
	}
	return slice
}

// 变异
func (g *GA) mutate(p []int) []int {
	if rand.Float32() < g.mutationProb {
		newOne := make([]int, len(p))
		copy(newOne, p)
		geneA := rand.Intn(len(p) - 1)
		geneB := rand.Intn(len(p) - 1)

		newOne[geneA], newOne[geneB] = newOne[geneB], newOne[geneA]
		return newOne
	}
	return p
}

// 交叉
func (g *GA) crossover(p1 []int, p2 []int) []int {
	geneA := rand.Intn(len(p1) - 1)

	newOne := []int{}
	for i := 0; i < geneA; i++ {
		newOne = append(newOne, p1[i])
	}
	for _, v := range p2 {
		if !In(newOne, v) {
			newOne = append(newOne, v)
		}
	}
	return newOne
}
