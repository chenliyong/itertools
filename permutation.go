package itertools

import (
	"sort"
)

// PermWithBest 全排列
type PermWithBest struct {
	// 全排列
	population [][]int
	// 适应度计算函数
	fitness func([]int) float32
}

// NewPermWithBest 返回全排列算法类
func NewPermWithBest(size int, fitness func([]int) float32) *PermWithBest {
	population := [][]int{}
	// 直接得到所有全排列
	for perm := range Permutations(size) {
		population = append(population, perm)
	}
	return &PermWithBest{
		population: population,
		fitness:    fitness,
	}
}

// Run 执行迭代
func (p *PermWithBest) Run() ([]int, float32) {
	best := p.population[p.ranked()[0].Index]
	fitness := p.fitness(best)
	return best, fitness
}

func (p *PermWithBest) ranked() []Ranked {
	size := len(p.population)
	popRanked := make([]Ranked, size)
	// 排序种群
	for i := 0; i < size; i++ {
		popRanked[i] = Ranked{
			Index: i,
			Rank:  p.fitness(p.population[i]),
		}
	}
	sort.Slice(popRanked, func(i, j int) bool {
		return popRanked[i].Rank > popRanked[j].Rank
	})
	return popRanked
}
