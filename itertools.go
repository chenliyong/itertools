package itertools

import "math/rand"

type Ranked struct {
	Index int
	Rank  float32
}

func SumRankedFits(ranked []Ranked) float32 {
	var f float32
	for _, v := range ranked {
		f = f + v.Rank
	}
	return f
}

func Range(n int) []int {
	result := []int{}
	for i := 0; i < n; i++ {
		result = append(result, i)
	}
	return result
}

func rangeAll(slice []int, start int, ch chan []int) {
	size := len(slice)
	if start == size-1 {
		//如果已经是最后位置了，直接将数组数据合并输出
		output := make([]int, len(slice))
		copy(output, slice[:])
		ch <- output
	}

	for i := start; i < size; i++ {
		//  1、abc输出abc，i=start时，输出自己
		//  2、如果i和start的值相同没有必要交换，交换后输出的是i=start时输出的内容，是重复的内容
		if i == start || slice[i] != slice[start] {
			//交换当前这个与后面的位置
			slice[i], slice[start] = slice[start], slice[i]
			//递归处理索引+1
			rangeAll(slice, start+1, ch)
			//换回来，因为是递归，如果不换回来会影响后面的操作，并且出现重复
			slice[i], slice[start] = slice[start], slice[i]
		}

	}
}

/**
 *  全排列
 */
func Permutations(n int) chan []int {
	ch := make(chan []int)
	go func() {
		defer close(ch)
		rangeAll(Range(n), 0, ch)
	}()
	return ch
}

func In(slice []int, c int) bool {
	for _, v := range slice {
		if v == c {
			return true
		}
	}
	return false
}

func RandomSample(slice []int) []int {
	for i := len(slice) - 1; i > 0; i-- {
		n := rand.Intn(i + 1)
		slice[i], slice[n] = slice[n], slice[i]
	}
	return slice
}
