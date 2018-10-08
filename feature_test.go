package geojson

import "testing"

func TestFeatureParse(t *testing.T) {
	p := expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]}}`, nil)
	expect(t, p.Center() == P(1, 2))
	expectJSON(t, `{"type":"Feature"}`, errGeometryMissing)
	expectJSON(t, `{"type":"Feature","geometry":null}`, errDataInvalid)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":null}`, errBBoxInvalid)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"id":[4,true]}`, nil)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]},"id":"15","properties":{"a":"b"}}`, nil)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":[1,2,3,4]}`, nil)
	expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2],"bbox":[1,2,3,4]},"id":[4,true]}`, nil)
}

func TestFeatureVarious(t *testing.T) {
	var g = expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]}}`, nil)
	expect(t, string(g.AppendJSON(nil)) == `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]}}`)
	expect(t, g.Rect() == R(1, 2, 1, 2))
	expect(t, g.Center() == P(1, 2))
	expect(t, !g.Empty())
	g = expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"bbox":[1,2,3,4]}`, nil)
	expect(t, !g.Empty())
	expect(t, g.Rect() == R(1, 2, 3, 4))
	expect(t, g.Center() == R(1, 2, 3, 4).Center())
}

// func TestFeaturePoly(t *testing.T) {
// 	p := expectJSON(t, `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2]}}`, nil)
// 	expect(t, p.Intersects(PO(1, 2)))
// 	expect(t, p.Contains(PO(1, 2)))
// 	expect(t, p.Within(PO(1, 2)))

// }
