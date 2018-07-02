package main

import "testing"

func TestURIGetHead(t *testing.T) {
	var (
		input  string = "abcdef12456/tail"
		expect string = "abcdef12456"
		head   string
		ok     bool
	)
	head, ok = URIGetHead(input)
	if !ok {
		t.Fatal("assertion failed, expected equal.")
	}
	if head != expect {
		t.Fatalf("inconsistent state, expected equal. Got (%s).\n", head)
	}
}

func TestParseURI(t *testing.T) {
	var (
		input      string   = "abcdef123456/test"
		expect     []string = []string{"abcdef123456", "test"}
		components []string
	)
	components = ParseURI(input)
	for i, c := range components {
		if c != expect[i] {
			t.Fatal("inconsistent state, expected equal.", c, expect[i])
		}
	}
}

func TestReadFile(t *testing.T) {
	var (
		path string = "/etc/hosts"
		err  error
	)
	_, err = readFile(path)
	if err != nil {
		t.Fatal("inconsistent state, expected equal.", err)
	}
}
