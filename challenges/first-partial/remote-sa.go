package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"math"
)

type Vector struct {
	X, Y float64
}

func (a Vector) cross(b Vector) float64 {
	return ((a.X * b.Y) - (a.Y * b.X))
}

type Point struct {
	X, Y float64
}

func (p1 Point) toVector(p2 Point) Vector{
	return Vector{(p2.X - p1.X), (p2.Y - p1.Y)}
}

type Edge struct {
	a, b Point
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//generatePoints array
func generatePoints(s string) ([]Point, error) {

	points := []Point{}

	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	vals := strings.Split(s, ",")
	if len(vals) < 2 {
		return []Point{}, fmt.Errorf("Point [%v] was not well defined", s)
	}

	var x, y float64

	for idx, val := range vals {

		if idx%2 == 0 {
			x, _ = strconv.ParseFloat(val, 64)
		} else {
			y, _ = strconv.ParseFloat(val, 64)
			points = append(points, Point{x, y})
		}
	}
	return points, nil
}

// max returns the larger of x or y.
func max(x, y float64) float64 {
    if x < y {
        return y
    }
    return x
}

// min returns the smaller of x or y.
func min(x, y float64) float64 {
    if x > y {
        return y
    }
    return x
}

// insideBounds determines if q is on the pr segment boundaries
func insideBounds(p, q, r Point) bool {
	if (q.X <= max(p.X, r.X)) && (q.X >= min(p.X, r.X)) && (q.Y <= max(p.Y, r.Y)) && (q.Y >= min(p.Y, r.Y)){
		return true
	}
    return false
}

func getOrientation(p, q, r Point) uint8{
	//value2 := ((q.Y - p.Y) * (r.X - q.X)) - ((q.X - p.X) * (r.Y - q.Y)) 
	pq := p.toVector(q)
	qr := q.toVector(r)
	value := pq.cross(qr)
	
	// ccw and cw orientations
	if value > 0 {
		return 1
	}
	if value < 0 {
		return 2
	}

	// collinear case: cross product is zero (triangle area is zero)
	return 0
}

// sample test: curl http://localhost:8000\?vertices=\(1,1\),\(10,1\),\(1,2\),\(10,2\)
func verifyComplexPoly(points []Point) bool {
	// a complex polygon is one where non-consecutive sides collide 

	//var eps float64 = 0.00001
	var edges []Edge

	for i := 0; i < len(points); i++ {
		temp := Edge{points[i], points[(i + 1) % len(points)]}
		edges = append(edges, temp)
	}

	// checking if any pair of non-contiguous line segment (edge) intersects
	for i := 0; i < len(edges); i++ {
		next := (i + 2) % len(edges)
		curr_edge := edges[i]
		
		j := next
		times := 0
		for times < len(edges) - 3 {
			if times == len(edges) - 2 {
				break
			}
			o1 := getOrientation(curr_edge.a, curr_edge.b, edges[j].a) 
			o2 := getOrientation(curr_edge.a, curr_edge.b, edges[j].b) 
			o3 := getOrientation(edges[j].a, edges[j].b, curr_edge.a) 
			o4 := getOrientation(edges[j].a, edges[j].b, curr_edge.b) 
		
			// general case
			if (o1 != o2) && (o3 != o4) {
				return true
			}

			// special cases
			// p1 q1 p2 are collinear and p2 lies on segment p1q1 
			if ((o1 == 0) && insideBounds(curr_edge.a, edges[j].a, curr_edge.b)) { return true }
			// p1 q1 q2 are collinear and q2 lies on segment p1q1 
			if ((o2 == 0) && insideBounds(curr_edge.a, edges[j].b, curr_edge.b)) { return true }
			// p2 q2 p1 are collinear and p1 lies on segment p2q2 
			if ((o3 == 0) && insideBounds(edges[j].a, curr_edge.a, edges[j].b)) { return true }
			// p2 q2 q1 are collinear and q1 lies on segment p2q2 
			if ((o4 == 0) && insideBounds(edges[j].a, curr_edge.b, edges[j].b)) { return true } 

			j = (j + 1) % len(edges)
			times++
			
		}
		
	}

	return false
}

func getDistance(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow((p2.X - p1.X), 2) + math.Pow((p2.Y - p1.Y), 2))
}

// getArea gets the area inside from a given shape
func getArea(points []Point) float64 {
	// Your code goes here
	var area float64 = 0.0

	for i := 0; i < len(points); i++ {
		idx := (i + 1) % len(points)
		area += math.Abs((points[i].X * points[idx].Y) - (points[idx].X * points[i].Y))
	}

	a := area / 2.0

	return a
}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {
	// Your code goes here

	var perimeter float64 = 0.0

	for i := 0; i < len(points); i++ {
		p1 := points[i]
		p2 := points[(i + 1) % len(points)]
		perimeter += getDistance(p1, p2)
	}

	return perimeter
}

// handler handles the web request and reponds it
func handler(w http.ResponseWriter, r *http.Request) {

	var vertices []Point
	for k, v := range r.URL.Query() {
		if k == "vertices" {
			points, err := generatePoints(v[0])
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("error: %v", err))
				return
			}
			vertices = points
			break
		}
	}

	// Results gathering
	area := getArea(vertices)
	perimeter := getPerimeter(vertices)

	// Logging in the server side
	log.Printf("Received vertices array: %v", vertices)

	// Response construction
	response := fmt.Sprintf("Welcome to the Remote Shapes Analyzer\n")
	response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", len(vertices))
	if verifyComplexPoly(vertices) {
		response += fmt.Sprintf("ERROR - Your shape has self-intersections, its area cannot be computed with this program.\n")
	} else if len(vertices) < 3 {
		response += fmt.Sprintf("ERROR - Your shape is not compliying with the minimum number of vertices.\n")
	} else {
		response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
		response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
		response += fmt.Sprintf(" - Area            : %v\n", area)
	}
	
	// Send response to client
	fmt.Fprintf(w, response)
}
