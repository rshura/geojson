package geojson

import "github.com/tidwall/gjson"

type LineString struct {
	Coordinates []Position
	BBox        BBox
	Extra       *Extra
}

func (g LineString) HasBBox() bool {
	return g.BBox != nil && g.BBox.Defined()
}
func (g LineString) Empty() bool {
	return len(g.Coordinates) == 0
}
func (g LineString) Rect() Rect {
	if g.BBox != nil {
		return g.BBox.Rect()
	}
	var rect Rect
	for i := 0; i < len(g.Coordinates); i++ {
		if i == 0 {
			rect.Min = g.Coordinates[i]
			rect.Max = g.Coordinates[i]
		} else {
			rect = rect.Expand(g.Coordinates[i])
		}
	}
	return rect
}

func (g LineString) Center() Position {
	return g.Rect().Center()
}

func (g LineString) AppendJSON(dst []byte) []byte {
	dst = append(dst, `{"type":"LineString","coordinates":[`...)
	for i, p := range g.Coordinates {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = appendJSONPosition(dst, p, g.Extra, i)
	}
	dst = append(dst, ']')
	if g.BBox != nil && g.BBox.Defined() {
		dst = append(dst, `,"bbox":`...)
		dst = g.BBox.AppendJSON(dst)
	}
	dst = append(dst, '}')
	return dst
}
func (g LineString) ForEach(func(child Object) bool) {}

func (g LineString) Contains(other Object) bool {
	if g.HasBBox() {
		return g.Rect().Contains(other)
	}
	otherRect := other.Rect()
	if otherRect.Min != otherRect.Max {
		return false
	}
	return polyLine(g.Coordinates).LineStringIntersectsPoint(
		polyPoint(otherRect.Min),
	)
}

func (g LineString) Intersects(other Object) bool {
	if g.HasBBox() {
		return g.Rect().Intersects(other)
	} else if other.HasBBox() {
		return other.Rect().Intersects(g)
	}
	if !g.Rect().IntersectsRect(other.Rect()) {
		return false
	}
	switch other := other.(type) {
	case Position:
		return polyLine(g.Coordinates).LineStringIntersectsPoint(
			polyPoint(other),
		)
	case Rect:
		return polyLine(g.Coordinates).LineStringIntersectsRect(
			polyRect(other),
		)
	case Point:
		return polyLine(g.Coordinates).LineStringIntersectsPoint(
			polyPoint(other.Coordinates),
		)
	case LineString:
		return polyLine(g.Coordinates).LineStringIntersectsLineString(
			polyLine(other.Coordinates),
		)
	case Polygon:
		return polyLine(g.Coordinates).LineStringIntersects(
			polyPolygon(other.Coordinates),
		)
	}
	// check types with children
	var intersects bool
	other.ForEach(func(child Object) bool {
		if g.Intersects(child) {
			intersects = true
			return false
		}
		return true
	})
	return intersects

}

func loadJSONLineString(data string) (Object, error) {
	var g LineString
	var err error
	g.Coordinates, g.Extra, err = loadJSONLineStringCoords(data, gjson.Result{})
	if err != nil {
		return nil, err
	}
	g.BBox, err = loadBBox(data)
	if err != nil {
		return nil, err
	}
	if g.BBox == nil {
		g.BBox = bboxRect{g.Rect()}
	}
	return g, nil
}

func loadJSONLineStringCoords(data string, rcoords gjson.Result) (
	[]Position, *Extra, error,
) {
	var err error
	var coords []Position
	var ex *Extra
	var dims int
	if !rcoords.Exists() {
		rcoords = gjson.Get(data, "coordinates")
		if !rcoords.Exists() {
			return nil, nil, errCoordinatesMissing
		}
		if !rcoords.IsArray() {
			return nil, nil, errCoordinatesInvalid
		}
	}
	rcoords.ForEach(func(key, value gjson.Result) bool {
		if !value.IsArray() {
			err = errCoordinatesInvalid
			return false
		}
		var count int
		var nums [4]float64
		value.ForEach(func(key, value gjson.Result) bool {
			if count == 4 {
				return false
			}
			if value.Type != gjson.Number {
				err = errCoordinatesInvalid
				return false
			}
			nums[count] = value.Float()
			count++
			return true
		})
		if err != nil {
			return false
		}
		if count < 2 {
			err = errCoordinatesInvalid
			return false
		}
		coords = append(coords, Position{X: nums[0], Y: nums[1]})
		if ex == nil {
			if count > 2 {
				ex = new(Extra)
				if count > 3 {
					ex.Dims = DimsZM
				} else {
					ex.Dims = DimsZ
				}
				dims = int(ex.Dims)
			}
		}
		if ex != nil {
			for i := 0; i < dims; i++ {
				ex.Positions = append(ex.Positions, nums[2+i])
			}
		}
		return true
	})
	if err != nil {
		return nil, nil, err
	}
	return coords, ex, err
}