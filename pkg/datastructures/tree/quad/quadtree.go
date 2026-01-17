package quad

import (
	"sync"
)

// Point represents a 2D coordinate.
type Point struct {
	X, Y float64
	Data interface{}
}

// Bounds represents a rectangular region.
type Bounds struct {
	X, Y, W, H float64 // Top-left X, Y and Width, Height
}

func (b Bounds) Contains(p Point) bool {
	return p.X >= b.X && p.X <= b.X+b.W && p.Y >= b.Y && p.Y <= b.Y+b.H
}

func (b Bounds) Intersects(other Bounds) bool {
	return !(other.X > b.X+b.W || other.X+other.W < b.X || other.Y > b.Y+b.H || other.Y+other.H < b.Y)
}

// Quadtree manages spatial indexing.
type Quadtree struct {
	boundary Bounds
	capacity int
	points   []Point
	divided  bool
	nw       *Quadtree
	ne       *Quadtree
	sw       *Quadtree
	se       *Quadtree
	mu       sync.RWMutex
}

func New(boundary Bounds, capacity int) *Quadtree {
	return &Quadtree{
		boundary: boundary,
		capacity: capacity,
		points:   make([]Point, 0, capacity),
	}
}

func (qt *Quadtree) Insert(p Point) bool {
	qt.mu.Lock()
	defer qt.mu.Unlock()
	return qt.insert(p)
}

func (qt *Quadtree) insert(p Point) bool {
	if !qt.boundary.Contains(p) {
		return false
	}

	if len(qt.points) < qt.capacity {
		qt.points = append(qt.points, p)
		return true
	}

	if !qt.divided {
		qt.subdivide()
	}

	if qt.nw.insert(p) {
		return true
	}
	if qt.ne.insert(p) {
		return true
	}
	if qt.sw.insert(p) {
		return true
	}
	if qt.se.insert(p) {
		return true
	}

	return false
}

func (qt *Quadtree) subdivide() {
	x := qt.boundary.X
	y := qt.boundary.Y
	w := qt.boundary.W / 2
	h := qt.boundary.H / 2

	qt.nw = New(Bounds{x, y, w, h}, qt.capacity)
	qt.ne = New(Bounds{x + w, y, w, h}, qt.capacity)
	qt.sw = New(Bounds{x, y + h, w, h}, qt.capacity)
	qt.se = New(Bounds{x + w, y + h, w, h}, qt.capacity)
	qt.divided = true
}

// Query returns points within range.
func (qt *Quadtree) Query(rangeBounds Bounds) []Point {
	qt.mu.RLock()
	defer qt.mu.RUnlock()

	var found []Point
	if !qt.boundary.Intersects(rangeBounds) {
		return found
	}

	for _, p := range qt.points {
		if rangeBounds.Contains(p) {
			found = append(found, p)
		}
	}

	if qt.divided {
		found = append(found, qt.nw.Query(rangeBounds)...)
		found = append(found, qt.ne.Query(rangeBounds)...)
		found = append(found, qt.sw.Query(rangeBounds)...)
		found = append(found, qt.se.Query(rangeBounds)...)
	}
	return found
}
