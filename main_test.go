package main

import "testing"

func BenchmarkImportWithReadAll(b *testing.B) {
	importFile(importReadAll)
	profileMemory("real-all")
}

func BenchmarkImportWithStream(b *testing.B) {
	importFile(importWithStream)
	profileMemory("stream")
}
