package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"math"
)

type Point struct {
	X, Y float64
}

type Vector struct {
	X, Y float64
}

type Edge struct {
	a, b Point
}

func (p1 Point) toVector(p2 Point) Vector{
	return Vector{(p2.X - p1.X), (p2.Y - p1.Y)}
}

func (a Vector) cross(b Vector) float64 {
	return ((a.X * b.Y) - (a.Y * b.X))
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

func getOrientation(p, q, r Point) uint8{
	value := ((q.Y - p.Y) * (r.X - q.X)) - ((q.X - p.X) * (r.Y - q.Y)) 
	pq := p.toVector(q)
	qr := q.toVector(r)
	value2 := pq.cross(qr)
	fmt.Printf("link %v me %v\n", value, value2)
	if value > 0 {
		return 1
	}
	if value < 0 {
		return 2
	}
	return 0
}

func verifyComplexPoly(points []Point) bool {
	// a complex polygon is one where non-consecutive sides collide 

	//var eps float64 = 0.00001
	var edges []Edge

	for i := 0; i < len(points); i++ {
		//p1p2 := points[i].toVector(points[(i + 1) % len(points)])

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
			//fmt.Printf("c: %v n: %v\n", curr_edge, edges[j])
			fmt.Println("-------------")
			o1 := getOrientation(curr_edge.a, curr_edge.b, edges[j].a) 
			o2 := getOrientation(curr_edge.a, curr_edge.b, edges[j].b) 
			o3 := getOrientation(edges[j].a, edges[j].b, curr_edge.a) 
			o4 := getOrientation(edges[j].a, edges[j].b, curr_edge.b) 
		
			if (o1 != o2) && (o3 != o4) {
				return true
			}

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
	response += fmt.Sprintf(" - Complex			: %v\n", verifyComplexPoly(vertices))
	response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", len(vertices))
	response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
	response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
	response += fmt.Sprintf(" - Area            : %v\n", area)

	// Send response to client
	fmt.Fprintf(w, response)
}
