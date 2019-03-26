package itertools

import "fmt"

type Point struct {
	x, y int
}

func NewPoint(x, y int) Point {
	return Point{
		x: x,
		y: y,
	}
}

type Points []Point

func (p Points) Reverse() {
	for i, j := 0, len(p)-1; i < j; i, j = i+1, j-1 {
		p[i], p[j] = p[j], p[i]
	}
}

type SquareGrid struct {
	width, height int
	walls         Points // 障碍物
	shelfs        Points // 货架
	weights       map[Point]int
}

func NewSquareGrid(w, h int) *SquareGrid {
	return &SquareGrid{
		width:  w,
		height: h,
	}
}

func (g *SquareGrid) AddWall(start Point, w, h int) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			g.walls = append(g.walls, NewPoint(start.x+i, start.y+j))
		}
	}
}

func (g *SquareGrid) AddShelf(start Point, width, quantity, corridorWidth, corridorNum int) {

	for n := 0; n < corridorNum; n++ {
		for y := 0; y < quantity; y++ {
			for x := 0; x < width; x++ {
				p := NewPoint(start.x+x, start.y+y)
				if g.InBounds(p) {
					g.shelfs = append(g.shelfs, p)
				}
			}
		}
		start = NewPoint(start.x+width+corridorWidth, start.y)
	}

}

func (g *SquareGrid) InBounds(p Point) bool {
	return 0 <= p.x && p.x < g.width && 0 <= p.y && p.y < g.height
}

func (g *SquareGrid) Passable(p Point) bool {
	block := append(g.walls, g.shelfs...)
	for _, w := range block {
		if w == p {
			return false
		}
	}
	return true
}

func (g *SquareGrid) Neighbors(p Point) Points {
	points := Points{
		Point{x: p.x + 1, y: p.y},
		Point{x: p.x, y: p.y - 1},
		Point{x: p.x - 1, y: p.y},
		Point{x: p.x, y: p.y + 1},
	}
	if (p.x+p.y)%2 == 0 {
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
	if v, ok := g.weights[to]; ok {
		return v
	}
	return 1
}

func (g *SquareGrid) getTile(id *Point, style map[string]Points) string {
	if v, ok := style["path"]; ok {
		if g.InPoints(id, v) {
			return "*"
		}
	}
	if v, ok := style["point"]; ok {
		if g.InPoints(id, v) {
			return "@"
		}
	}
	if g.InPoints(id, g.shelfs) {
		return "#"
	}
	if g.InPoints(id, g.walls) {
		return "x"
	}
	return "."
}

func (g *SquareGrid) InPoints(id *Point, points Points) bool {
	for _, v := range points {
		if v.x == id.x && v.y == id.y {
			return true
		}
	}
	return false
}

func (g *SquareGrid) Draw(style map[string]Points) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			fmt.Print(g.getTile(&Point{x: x, y: y}, style))
		}
		fmt.Print("\n")
	}
}
