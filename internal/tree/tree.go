package tree

import (
	"fmt"
	"pattern-reduce/internal/distance"
	"pattern-reduce/internal/pattern"
)

type Node struct {
	Cluster *pattern.Cluster
	Left    *Node
	Right   *Node
	Depth   int
}

func NewRoot(c *pattern.Cluster) *Node {
	return &Node{
		Cluster: c,
		Left:    nil,
		Right:   nil,
		Depth:   0,
	}
}

func (n *Node) NewChild(c *pattern.Cluster) *Node {
	return &Node{
		Cluster: c,
		Left:    nil,
		Right:   nil,
		Depth:   n.Depth + 1,
	}
}

func (n *Node) Split(d *distance.EditDistanceCalculator, tid int, rem int) bool {

	left := n.Cluster.Patterns[0]

	excluded := make([]bool, len(n.Cluster.Patterns))
	excluded[0] = true
	leftDistances := make([]float32, len(n.Cluster.Patterns))

	maxDist := float32(0)
	maxDistIndex := 0
	for i, other := range n.Cluster.Patterns {
		if excluded[i] {
			continue
		}
		dist := distance.PatternDistance(left.Pattern, other.Pattern, d)
		leftDistances[i] = dist
		if dist == 0 {
			left.Frequency++
			excluded[i] = true
			continue
		}
		if dist > maxDist {
			maxDist = dist
			maxDistIndex = i
		}
	}

	right := n.Cluster.Patterns[maxDistIndex]
	excluded[maxDistIndex] = true

	leftPatterns := []*pattern.Tracker{left}
	rightPatterns := []*pattern.Tracker{right}
	hasSplit := false

	for i, other := range n.Cluster.Patterns {
		if excluded[i] {
			continue
		}
		leftDist := leftDistances[i]
		rightDist := distance.PatternDistance(right.Pattern, other.Pattern, d)
		if rightDist == 0 {
			right.Frequency++
			continue
		}
		if leftDist < rightDist {
			leftPatterns = append(leftPatterns, other)
			hasSplit = true
		} else if rightDist < leftDist {
			rightPatterns = append(rightPatterns, other)
			hasSplit = true
		} else {
			leftPatterns = append(leftPatterns, other)
			rightPatterns = append(rightPatterns, other)
		}
	}

	if hasSplit {
		fmt.Printf("[thread %d]\tsplit\t| depth: %d | dist: %f | left: %d | right: %d | buffer: %d\n", tid, n.Depth, maxDist, len(leftPatterns), len(rightPatterns), rem)
		//fmt.Println(left.Pattern.String())
		//fmt.Println(right.Pattern.String())
		n.Left = n.NewChild(pattern.NewCluster(leftPatterns))
		n.Right = n.NewChild(pattern.NewCluster(rightPatterns))
	} else {
		if len(n.Cluster.Patterns) > 100 {
			fmt.Printf("[thread %d]\tend\t| depth: %d | patterns: %d\n", tid, n.Depth, len(n.Cluster.Patterns))
			fmt.Println(n.Cluster.Patterns[0].Pattern.String())
		}
	}

	return hasSplit
}
