package itertools

import "fmt"

type Point struct {
	X, Y int
}

func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

type Points []Point

func (p Points) Reverse() {
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
}

type SquareGrid struct {
	Width, Height int
	Walls         Points // 障碍物
	Shelfs        Points // 货架
	Weights       map[Point]int
}

func NewSquareGrid(w, h int) *SquareGrid {
	return &SquareGrid{
		Width:  w,
		Height: h,
	}
}

func (g *SquareGrid) AddWall(start Point, w, h int) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			g.Walls = append(g.Walls, NewPoint(start.X+i, start.Y+j))
		}
	}
}

func (g *SquareGrid) AddShelf(start Point, width, quantity, corridorWidth, corridorNum int) {

	for n := 0; n < corridorNum; n++ {
		for y := 0; y < quantity; y++ {
			for x := 0; x < width; x++ {
				p := NewPoint(start.X+x, start.Y+y)
				if g.InBounds(p) {
					g.Shelfs = append(g.Shelfs, p)
				}
			}
		}
		start = NewPoint(start.X+width+corridorWidth, start.Y)
	}

}

func (g *SquareGrid) InBounds(p Point) bool {
	return 0 <= p.X && p.X < g.Width && 0 <= p.Y && p.Y < g.Height
}

func (g *SquareGrid) Passable(p Point) bool {
	block := append(g.Walls, g.Shelfs...)
	for _, w := range block {
		if w == p {
			return false
		}
	}
	return true
}

func (g *SquareGrid) Neighbors(p Point) Points {
	points := Points{
		Point{X: p.X + 1, Y: p.Y},
		Point{X: p.X, Y: p.Y - 1},
		Point{X: p.X - 1, Y: p.Y},
		Point{X: p.X, Y: p.Y + 1},
	}
	if (p.X+p.Y)%2 == 0 {
		points.Reverse()
	}
	results := Points{}
	for _, i := range points {
		if g.InBounds(i) && g.Passable(i) {
			results = append(results, i)
		}
	}
	return results
}

func (g *SquareGrid) Cost(from, to Point) int {
	if v, ok := g.Weights[to]; ok {
		return v
	}
	return 1
}

func (g *SquareGrid) getTile(id *Point, style map[string]Points) string {
	if v, ok := style["point"]; ok {
		if g.InPoints(id, v) {
			return "@"
		}
	}
	if v, ok := style["path"]; ok {
		if g.InPoints(id, v) {
			return "^"
		}
	}
	if g.InPoints(id, g.Shelfs) {
		return "#"
	}
	if g.InPoints(id, g.Walls) {
		return "x"
	}
	return "."
}

func (g *SquareGrid) InPoints(id *Point, points Points) bool {
	for _, v := range points {
		if v.X == id.X && v.Y == id.Y {
			return true
		}
	}
	return false
}

func (g *SquareGrid) Draw(style map[string]Points) {
	for y := 0; y < g.Height; y++ {
		fmt.Printf("%02d", y)
		for x := 0; x < g.Width; x++ {
			fmt.Print(g.getTile(&Point{X: x, Y: y}, style) + " ")
		}
		fmt.Print("\n")
	}
	for i := 0; i < g.Width; i += 2 {
		fmt.Printf("  %02d", i)
	}
	fmt.Print("\n")
}
