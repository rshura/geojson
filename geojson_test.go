package geojson

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/tidwall/pretty"
)

func expectJSON(t *testing.T, data string, exp error) Object {
	t.Helper()
	obj, err := Load(data)
	if err != exp {
		t.Fatalf("expected '%v', got '%v'", exp, err)
	}
	return obj
}

func cleanJSON(data string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(data), &v); err != nil {
		panic(err)
	}
	dst, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	opts := *pretty.DefaultOptions
	opts.Width = 99999999
	return string(pretty.PrettyOptions(dst, &opts))
}

func testGeoJSONFile(t *testing.T, path string) Object {
	t.Helper()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	obj, err := Load(string(data))
	if err != nil {
		t.Fatal(err)
	}
	orgJSON := cleanJSON(string(data))
	newJSON := cleanJSON(string(obj.AppendJSON(nil)))
	if orgJSON != newJSON {
		var ln int
		var col int
		for i := 0; i < len(orgJSON) && i < len(newJSON); i++ {
			if orgJSON[i] != newJSON[i] {
				break
			}
			if orgJSON[i] == '\n' {
				ln++
				col = 0
			} else {
				col++
			}
		}
		tpath1 := "/tmp/org.json"
		tpath2 := "/tmp/new.json"
		ioutil.WriteFile(tpath1, []byte(orgJSON), 0666)
		ioutil.WriteFile(tpath2, []byte(newJSON), 0666)
		t.Fatalf("%v (ln: %d, col: %d)\nfile://%s\nfile://%s",
			filepath.Base(path), ln, col, tpath1, tpath2)
	}
	return obj
}

func expect(t *testing.T, v bool) {
	t.Helper()
	if !v {
		t.Fatal("invalid expectation")
	}
}

func TestGeoJSON(t *testing.T) {
	fis, err := ioutil.ReadDir("test_files")
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range fis {
		testGeoJSONFile(t, filepath.Join("test_files", fi.Name()))
	}
}