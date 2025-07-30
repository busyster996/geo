# geo

A Go library for 2D computational geometry, providing data structures and algorithms for working with coordinates, vectors, lines, polygons, circles, rectangles, triangles, and convex shapes. It is suitable for applications such as collision detection, geometric queries, and spatial reasoning.

## Features

- **Basic Structures**: 
  - `Coord`: 2D coordinate (X, Z).
  - `Vector`: 2D vector with operations (add, subtract, dot/cross product, rotation, length).
- **Geometric Shapes**:
  - `Segment`: Line segment with intersection and distance calculation methods.
  - `Rectangle`: With methods for point-in-rectangle, random point generation, and vector representation.
  - `Circle`: With methods for intersection with segments/polygons, point containment, and bounding rectangle.
  - `Triangle` and `Convex`: With point-in-shape tests, merging, and bounding box calculation.
  - `Polygon`: Interface for generic polygons.
- **Collision & Intersection**:
  - Check if shapes intersect (circle-polygon, segment-circle, segment-segment, etc).
  - Calculate intersection points between lines and shapes.
- **Spatial Queries**:
  - Determine location of points relative to borders or shapes.
  - Calculate distances between points and shapes.
- **Utilities**:
  - Random coordinate generation within rectangles.
  - Convex hull and convexity checks.
  - Edge and vertex management for complex shapes.

## Installation

```bash
go get github.com/busyster996/geo
```

## Usage

```go
import "github.com/busyster996/geo"

func main() {
    // Create coordinates
    a := geo.NewCoord(0, 0)
    b := geo.NewCoord(10, 0)
    c := geo.NewCoord(5, 5)

    // Create a segment and calculate its length
    seg := geo.NewSegment(a, b)
    length := seg.ToVector().Length()

    // Create a circle and check intersection with a segment
    circle := geo.NewCirCle(b, 5)
    intersects := circle.IsIntersect(&seg)

    // Create a rectangle and check if it contains a point
    rect := geo.NewRectangle(0, 0, 10, 10)
    inside := rect.IsCoordInside(c)
}
```

## Main Data Structures

- **Coord**: Represents a point in 2D.
- **Vector**: Direction and magnitude in 2D, with many vector operations.
- **Segment**: A line segment defined by two coordinates.
- **Rectangle**: Axis-aligned rectangle.
- **Circle**: Defined by a center and radius, includes intersection and containment logic.
- **Triangle/Convex**: For polygonal calculations, point inclusion, merging, and sorting.
- **Border**: Special rectangular border with location state queries.
- **Polygon (interface)**: For generic polygonal operations and queries.

## Example Capabilities

- Test if a point is inside a polygon/triangle/rectangle/circle.
- Find the intersection between a circle and a segment.
- Calculate the distance between points, or from point to segment.
- Sort vertices in counterclockwise order for convex polygons.
- Merge triangles into convex polygons and validate convexity.

## License

This project is licensed under the [Apache License 2.0](LICENSE).

---
**Repository:** https://github.com/busyster996/geo
