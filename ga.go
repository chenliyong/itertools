package itertools

import (
	"math/rand"
	"sort"
	"time"
)

// FitnessFunc 适应度函数
type FitnessFunc func([]int) float32

// MutateFunc 变异函数
type MutateFunc func(p []int, prob float32) []int

// CrossoverFunc 交叉函数
type CrossoverFunc func(p1 []int, p2 []int) []int

// GA 遗传算法类
type GA struct {
	mutationProb, eliteProb         float32
	popSize, maxIter, maxIterRepeat int
	// 种群
	population [][]int
	// 适应度计算函数
	fitness FitnessFunc
	// 变异
	mutate MutateFunc
	// 交叉
	crossover CrossoverFunc
}

// GAOption call options
type GAOption func(*GA)

// WithMutationProb 变异概率
func WithMutationProb(prob float32) GAOption {
	return func(o *GA) {
		o.mutationProb = prob
	}
}

// WithEliteProb 精英概率
func WithEliteProb(prob float32) GAOption {
	return func(o *GA) {
		o.eliteProb = prob
	}
}

// WithMaxIter 最大迭代次数
func WithMaxIter(i int) GAOption {
	return func(o *GA) {
		o.maxIter = i
	}
}

// WithMaxIterRepeat 提前收敛的条件次数
func WithMaxIterRepeat(i int) GAOption {
	return func(o *GA) {
		o.maxIterRepeat = i
	}
}

// WithFiness 适应度函数
func WithFiness(f FitnessFunc) GAOption {
	return func(o *GA) {
		o.fitness = f
	}
}

// WithCrossover 交叉函数
func WithCrossover(f CrossoverFunc) GAOption {
	return func(o *GA) {
		o.crossover = f
	}
}

// WithMutate 变异函数
func WithMutate(f MutateFunc) GAOption {
	return func(o *GA) {
		o.mutate = f
	}
}

// NewGA 返回遗传算法类
func NewGA(population [][]int, opts ...GAOption) *GA {
	ga := &GA{
		population:    population,
		popSize:       len(population),
		crossover:     defaultCrossover,
		mutate:        defaultMutate,
		mutationProb:  0.2,
		eliteProb:     0.1,
		maxIter:       -1,
		maxIterRepeat: 10,
	}
	for _, o := range opts {
		o(ga)
	}
	// 防止无限迭代 maxIter和maxIterRepeat不能都为-1
	if ga.maxIter == -1 && ga.maxIterRepeat == -1 {
		ga.maxIterRepeat = 10
	}
	return ga
}

// Run 执行迭代命令
func (g *GA) Run() {
	var n, iterRepeat int
	for {
		if g.maxIter != -1 && n > g.maxIter {
			break
		}
		// 进行一次迭代
		newPop := g.bread()
		if g.fitness(g.population[0]) == g.fitness(newPop[0]) {
			iterRepeat = iterRepeat + 1
			if g.maxIterRepeat != -1 && iterRepeat >= g.maxIterRepeat {
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

// Best 返回算法最后迭代结果
func (g *GA) Best() ([]int, float32) {
	result := g.population[0]
	rank := g.fitness(result)
	return result, rank
}

func (g *GA) ranked() []Ranked {
	popRanked := make([]Ranked, g.popSize)
	// 排序种群
	for i := 0; i < g.popSize; i++ {
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
	var pick, cumSum float32

	popRanked := g.ranked()
	results := [][]int{}
	sumFits := SumRankedFits(popRanked)

	//  精英选择
	eliteSize := int(float32(g.popSize) * g.eliteProb)
	for i := 0; i < eliteSize; i++ {
		results = append(results, g.population[popRanked[i].Index])
	}

	// 轮盘选择
	// TODO 是否考虑使用退火算法接受新个体
	for i := 0; i < len(popRanked)-eliteSize; i++ {
		pick = 100 * rand.Float32()
		cumSum = 0
		for j := 0; j < len(popRanked); j++ {
			cumSum = cumSum + popRanked[j].Rank
			if pick <= 100*cumSum/sumFits {
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
	eliteSize := int(float32(g.popSize) * g.eliteProb)
	for i := 0; i < eliteSize; i++ {
		children = append(children, pool[i])
	}

	// 乱序
	pool = shuffle(pool)
	for i := 0; i < len(pool)-eliteSize; i++ {
		c := g.crossover(pool[i], pool[len(pool)-1-i])
		c = g.mutate(c, g.mutationProb)
		children = append(children, c)
	}

	return children

}

// MakePopulation 创建种群，给定种群大小popSize及染色体长度itemSize
func MakePopulation(popSize, itemSize int) [][]int {
	rand.Seed(time.Now().Unix())
	population := make([][]int, popSize)
	for i := 0; i < popSize; i++ {
		population[i] = RandomSample(Range(itemSize))
	}
	return population
}

// 种群乱序处理
func shuffle(slice [][]int) [][]int {
	for i := len(slice) - 1; i > 0; i-- {
		n := rand.Intn(i + 1)
		slice[i], slice[n] = slice[n], slice[i]
	}
	return slice
}

// 默认变异方法
func defaultMutate(p []int, prob float32) []int {
	if rand.Float32() < prob {
		newOne := make([]int, len(p))
		copy(newOne, p)
		geneA := rand.Intn(len(p) - 1)
		geneB := rand.Intn(len(p) - 1)

		newOne[geneA], newOne[geneB] = newOne[geneB], newOne[geneA]
		return newOne
	}
	return p
}

// 默认交叉方法
func defaultCrossover(p1 []int, p2 []int) []int {
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
