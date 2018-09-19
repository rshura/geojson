package poly

// IntersectsLineString detect if a point intersects a linestring
func (p Point) IntersectsLineString(exterior Polygon) bool {
	for j := 0; j < len(exterior); j++ {
		if raycast(p, exterior[j], exterior[(j+1)%len(exterior)]).on {
			return true
		}
	}
	return false
}

// Intersects detects if a point intersects another polygon
func (p Point) Intersects(exterior Polygon, holes []Polygon) bool {
	return p.Inside(exterior, holes)
}

// Intersects detects if a rect intersects another polygon
func (r Rect) Intersects(exterior Polygon, holes []Polygon) bool {
	return r.Polygon().Intersects(exterior, holes)
}

// Intersects detects if a polygon intersects another polygon
func (shape Polygon) Intersects(exterior Polygon, holes []Polygon) bool {
	return shape.doesIntersects(false, exterior, holes)
}

// LineStringIntersectsLineString detects if a linestring intersects a
// linestring assume shape and exterior are actually linestrings
func (shape Polygon) LineStringIntersectsLineString(exterior Polygon) bool {
	for i := 0; i < len(shape)-1; i++ {
		for j := 0; j < len(exterior)-1; j++ {
			if lineintersects(
				shape[i], shape[i+1],
				exterior[i], exterior[i+1],
			) {
				return true
			}
		}
	}
	return false
}

// LineStringIntersectsPoint detects if a linestring intersects a point. The
// point will need to be exactly on a segment of the linestring
func (shape Polygon) LineStringIntersectsPoint(point Point) bool {
	return point.IntersectsLineString(shape)
}

// LineStringIntersectsRect detects if a linestring intersects a rect
func (shape Polygon) LineStringIntersectsRect(rect Rect) bool {
	return shape.LineStringIntersects(rect.Polygon(), nil)
}

// LineStringIntersects detects if a polygon intersects a linestring
// assume shape is a linestring
func (shape Polygon) LineStringIntersects(
	exterior Polygon, holes []Polygon,
) bool {
	return shape.doesIntersects(true, exterior, holes)
}
func (shape Polygon) doesIntersects(
	isLineString bool, exterior Polygon, holes []Polygon,
) bool {
	switch len(shape) {
	case 0:
		return false
	case 1:
		switch len(exterior) {
		case 0:
			return false
		case 1:
			return shape[0].X == exterior[0].X && shape[0].Y == shape[0].Y
		default:
			return shape[0].Inside(exterior, holes)
		}
	default:
		switch len(exterior) {
		case 0:
			return false
		case 1:
			return exterior[0].Inside(shape, holes)
		}
	}
	if !shape.Rect().IntersectsRect(exterior.Rect()) {
		return false
	}
	for i := 0; i < len(shape); i++ {
		for j := 0; j < len(exterior); j++ {
			if lineintersects(
				shape[i], shape[(i+1)%len(shape)],
				exterior[j], exterior[(j+1)%len(exterior)],
			) {
				return true
			}
		}
	}
	for _, hole := range holes {
		if shape.Inside(hole, nil) {
			return false
		}
	}
	if shape.Inside(exterior, nil) {
		return true
	}
	if !isLineString {
		if exterior.Inside(shape, nil) {
			return true
		}
	}
	return false
}

func lineintersects(
	a, b Point, // segment 1
	c, d Point, // segment 2
) bool {
	// do the bounding boxes intersect?
	// the following checks without swapping values.
	if a.Y > b.Y {
		if c.Y > d.Y {
			if b.Y > c.Y || a.Y < d.Y {
				return false
			}
		} else {
			if b.Y > d.Y || a.Y < c.Y {
				return false
			}
		}
	} else {
		if c.Y > d.Y {
			if a.Y > c.Y || b.Y < d.Y {
				return false
			}
		} else {
			if a.Y > d.Y || b.Y < c.Y {
				return false
			}
		}
	}
	if a.X > b.X {
		if c.X > d.X {
			if b.X > c.X || a.X < d.X {
				return false
			}
		} else {
			if b.X > d.X || a.X < c.X {
				return false
			}
		}
	} else {
		if c.X > d.X {
			if a.X > c.X || b.X < d.X {
				return false
			}
		} else {
			if a.X > d.X || b.X < c.X {
				return false
			}
		}
	}

	// the following code is from http://ideone.com/PnPJgb
	cmpx, cmpy := c.X-a.X, c.Y-a.Y
	rx, ry := b.X-a.X, b.Y-a.Y
	cmpxr := cmpx*ry - cmpy*rx
	if cmpxr == 0 {
		// Lines are collinear, and so intersect if they have any overlap
		if !(((c.X-a.X <= 0) != (c.X-b.X <= 0)) ||
			((c.Y-a.Y <= 0) != (c.Y-b.Y <= 0))) {
			return false
		}
		return true
	}
	sx, sy := d.X-c.X, d.Y-c.Y
	cmpxs := cmpx*sy - cmpy*sx
	rxs := rx*sy - ry*sx
	if rxs == 0 {
		return false // Lines are parallel.
	}
	rxsr := 1 / rxs
	t := cmpxs * rxsr
	u := cmpxr * rxsr
	if !((t >= 0) && (t <= 1) && (u >= 0) && (u <= 1)) {
		return false
	}
	return true
}
