package helputils

import (
	"testing"
)

func TestCpuCount(t *testing.T) {
	cpus := CpuCount()
	if cpus > 0 {
		t.Logf("Found %d cpus", cpus)
	} else {
		t.Fatal("Failed to count cpus")
	}
}

func TestResolvedHost(t *testing.T) {
	for test, pass := range map[string]bool{
		"Loader-1-207.lab.il.elastifile.com": true,
		"8.8.8.8": true,
		"bla":     false,
	} {
		result := ResolvedHost(test)
		if pass != bool(result != "") {
			t.Fatalf("Failed test: %s, expected: %v, result: %v", test, pass, result)
		} else {
			t.Logf("Passed test: %s", test)
		}
	}
}
