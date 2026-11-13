package geofence

import (
	"errors"
	"fmt"
	"math"
)

const (
	minLatitude  = -90.0
	maxLatitude  = 90.0
	minLongitude = -180.0
	maxLongitude = 180.0
	epsilon      = 1e-10
)

// ValidatePolygon validates a closed polygon whose coordinates use
// [latitude, longitude] order.
func ValidatePolygon(coordinates [][]float64) error {
	if len(coordinates) < 4 {
		return errors.New("polygon must contain at least four coordinate pairs")
	}

	for i, point := range coordinates {
		if len(point) != 2 {
			return fmt.Errorf("coordinate %d must contain latitude and longitude", i)
		}

		latitude := point[0]
		longitude := point[1]

		if math.IsNaN(latitude) || math.IsInf(latitude, 0) ||
			math.IsNaN(longitude) || math.IsInf(longitude, 0) {
			return fmt.Errorf("coordinate %d contains a non-finite value", i)
		}

		if latitude < minLatitude || latitude > maxLatitude {
			return fmt.Errorf("coordinate %d has invalid latitude", i)
		}

		if longitude < minLongitude || longitude > maxLongitude {
			return fmt.Errorf("coordinate %d has invalid longitude", i)
		}
	}

	first := coordinates[0]
	last := coordinates[len(coordinates)-1]

	if !samePoint(first, last) {
		return errors.New("polygon must be closed: first and last coordinates must match")
	}

	uniquePoints := make(map[[2]float64]struct{})

	for _, point := range coordinates[:len(coordinates)-1] {
		uniquePoints[[2]float64{point[0], point[1]}] = struct{}{}
	}

	if len(uniquePoints) < 3 {
		return errors.New("polygon must contain at least three unique vertices")
	}

	if math.Abs(polygonArea(coordinates)) < epsilon {
		return errors.New("polygon area must be greater than zero")
	}

	if hasSelfIntersection(coordinates) {
		return errors.New("polygon must not intersect itself")
	}

	return nil
}

// IsPointInPolygon checks whether a GPS point is inside or on the boundary
// of a polygon. Coordinates use [latitude, longitude] order.
func IsPointInPolygon(latitude, longitude float64, polygon [][]float64) bool {
	if len(polygon) < 4 {
		return false
	}

	if math.IsNaN(latitude) || math.IsInf(latitude, 0) ||
		math.IsNaN(longitude) || math.IsInf(longitude, 0) {
		return false
	}

	if latitude < minLatitude || latitude > maxLatitude ||
		longitude < minLongitude || longitude > maxLongitude {
		return false
	}

	inside := false

	// Convert to x = longitude and y = latitude for ray casting.
	x := longitude
	y := latitude

	for i, j := 0, len(polygon)-1; i < len(polygon); j, i = i, i+1 {
		if len(polygon[i]) != 2 || len(polygon[j]) != 2 {
			return false
		}

		xi, yi := polygon[i][1], polygon[i][0]
		xj, yj := polygon[j][1], polygon[j][0]

		if pointOnSegment(x, y, xi, yi, xj, yj) {
			return true
		}

		intersects := (yi > y) != (yj > y) &&
			x < (xj-xi)*(y-yi)/(yj-yi)+xi

		if intersects {
			inside = !inside
		}
	}

	return inside
}

// samePoint checks whether two coordinate pairs are equal.
func samePoint(a, b []float64) bool {
	return len(a) == 2 &&
		len(b) == 2 &&
		math.Abs(a[0]-b[0]) < epsilon &&
		math.Abs(a[1]-b[1]) < epsilon
}

// polygonArea calculates the signed polygon area using the shoelace formula.
// X is longitude and Y is latitude.
func polygonArea(polygon [][]float64) float64 {
	area := 0.0

	for i := 0; i < len(polygon)-1; i++ {
		x1, y1 := polygon[i][1], polygon[i][0]
		x2, y2 := polygon[i+1][1], polygon[i+1][0]

		area += x1*y2 - x2*y1
	}

	return area / 2
}

// pointOnSegment checks whether a point lies on a line segment.
func pointOnSegment(px, py, ax, ay, bx, by float64) bool {
	cross := (px-ax)*(by-ay) - (py-ay)*(bx-ax)

	if math.Abs(cross) > epsilon {
		return false
	}

	dot := (px-ax)*(bx-ax) + (py-ay)*(by-ay)

	if dot < -epsilon {
		return false
	}

	lengthSquared := (bx-ax)*(bx-ax) + (by-ay)*(by-ay)

	return dot <= lengthSquared+epsilon
}

// hasSelfIntersection checks whether non-adjacent polygon edges intersect.
func hasSelfIntersection(polygon [][]float64) bool {
	edgeCount := len(polygon) - 1

	for i := 0; i < edgeCount; i++ {
		a1 := polygon[i]
		a2 := polygon[i+1]

		for j := i + 1; j < edgeCount; j++ {
			// Adjacent edges share a legitimate endpoint.
			if j == i+1 || (i == 0 && j == edgeCount-1) {
				continue
			}

			b1 := polygon[j]
			b2 := polygon[j+1]

			if segmentsIntersect(
				a1[1], a1[0],
				a2[1], a2[0],
				b1[1], b1[0],
				b2[1], b2[0],
			) {
				return true
			}
		}
	}

	return false
}

// segmentsIntersect detects intersections between two line segments.
func segmentsIntersect(
	ax, ay, bx, by,
	cx, cy, dx, dy float64,
) bool {
	o1 := orientation(ax, ay, bx, by, cx, cy)
	o2 := orientation(ax, ay, bx, by, dx, dy)
	o3 := orientation(cx, cy, dx, dy, ax, ay)
	o4 := orientation(cx, cy, dx, dy, bx, by)

	if o1*o2 < 0 && o3*o4 < 0 {
		return true
	}

	if math.Abs(o1) <= epsilon && pointOnSegment(cx, cy, ax, ay, bx, by) {
		return true
	}

	if math.Abs(o2) <= epsilon && pointOnSegment(dx, dy, ax, ay, bx, by) {
		return true
	}

	if math.Abs(o3) <= epsilon && pointOnSegment(ax, ay, cx, cy, dx, dy) {
		return true
	}

	if math.Abs(o4) <= epsilon && pointOnSegment(bx, by, cx, cy, dx, dy) {
		return true
	}

	return false
}

// orientation returns the cross product used for segment intersection tests.
func orientation(ax, ay, bx, by, cx, cy float64) float64 {
	return (bx-ax)*(cy-ay) - (by-ay)*(cx-ax)
}
