package itertools

type Node struct {
	priority int
	point    Point
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func heuristic(p1, p2 Point) int {
	return abs(p1.x-p2.x) + abs(p1.y-p2.y)
}

func AstarSearch(graph *SquareGrid, start, end Point) Points {
	cameFrom := make(map[Point]*Point)
	costSoFar := make(map[Point]int)

	frontier := Queue{}
	frontier.put(&Node{point: start, priority: 0})

	cameFrom[start] = nil
	costSoFar[start] = 0

	for {
		if frontier.empty() {
			break
		}
		node := frontier.get().(*Node)
		current := node.point

		if current == end {
			break
		}

		for _, next := range graph.Neighbors(current) {

			newCost := costSoFar[current] + graph.Cost(current, next)

			if v, ok := costSoFar[next]; !ok || newCost < v {
				costSoFar[next] = newCost
				priority := newCost + heuristic(end, next)
				frontier.put(&Node{point: next, priority: priority})
				cameFrom[next] = &current
			}
		}
	}

	// reconstruct path
	current := end
	path := Points{}

	for {
		path = append(path, current)
		if current == start {
			break
		}
		current = *cameFrom[current]
	}

	path.Reverse()

	return path
}
