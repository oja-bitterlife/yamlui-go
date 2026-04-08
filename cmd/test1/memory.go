package main

import (
	"runtime"

	"github.com/oja-bitterlife/yamlui-go/script"
)

func PrintMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// HeapAlloc: 現在使用中のヒープ
	// Mallocs: 累計アロケーション回数 (これが増え続けていればGC予備軍)
	println("Alloc:", m.HeapAlloc, " Mallocs:", m.Mallocs, " Sys:", m.HeapSys)
}

var lastMallocs uint64

func TraceMemory() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 前回からの増加分を計算
	diff := m.Mallocs - lastMallocs
	lastMallocs = m.Mallocs

	if diff > 0 {
		script.Log("Alloc_KB %d, New_Mallocs %d, Total_Mallocs %d", m.HeapAlloc/1024, diff, m.Mallocs)
	}
}
