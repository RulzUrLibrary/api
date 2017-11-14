package test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
)

func helperBodyCompare(t *testing.T, resp *http.Response, name string) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	expected := helperLoadBytes(t, name)
	if !bytes.Equal(body, expected) {
		t.Errorf("compare body failed against %s:\n%s", name, body)
	}
}

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes[:len(bytes)-1]
}
