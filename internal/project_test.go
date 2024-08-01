package internal

import (
	// "fmt"
	// "time"
	"testing"
)

func BenchmarkNewEntrytoYamlBytesBuf(b *testing.B) {
	proj := Project{
		name:     "Sample Project",
		desc:     "This is a description",
		path:     "/path/to/project",
		language: "Go",
		active:   true,
	}
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = proj.newEntrytoYamlBytes()
	}
}

func BenchmarkRetrieveData(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = retrieveData()
	}
}
