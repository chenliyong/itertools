package main

import (
	"github.com/laoqiu/itertools"
)

func main() {
	grid := itertools.NewSquareGrid(62, 50)
	grid.AddShelf(itertools.NewPoint(0, 6), 2, 20, 2, 1)
	grid.AddShelf(itertools.NewPoint(4, 6), 4, 20, 2, 10)
	grid.AddShelf(itertools.NewPoint(0, 28), 2, 15, 2, 1)
	grid.AddShelf(itertools.NewPoint(4, 28), 4, 15, 2, 10)
	// points := itertools.AstarSearch(grid, itertools.NewPoint(3, 19), itertools.NewPoint(50, 45))
	points := itertools.Points{
		itertools.NewPoint(12, 22),
		itertools.NewPoint(14, 22),
	}
	style := map[string]itertools.Points{
		"point": points,
	}
	grid.Draw(style)
}
