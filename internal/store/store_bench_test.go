package store

import (
	"path/filepath"
	"testing"

	
)

func BenchmarkSet(b *testing.B) {
	//create temp store
	dir := b.TempDir()
	path := filepath.Join(dir, "BenchmarkingSet.log")

	s, err := NewStore(path)

	if err != nil {
		b.Fatal(err)
	}

	defer s.Close()

	//reset timer
	b.ResetTimer()

	//call set b.N times
	for i := 0; i < b.N; i++ {
		if err := s.Set("Name", "Anshit");err!=nil{
			b.Fatal(err)
		}
	}
}
