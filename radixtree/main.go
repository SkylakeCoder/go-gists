package main

import (
	"fmt"
	"time"
)

var testCount = 10000000

func testRadixTree() {
	rt := NewRadixTree()
	rt.Insert("/golang")
	rt.Insert("/hello")
	rt.Insert("/hehe")
	rt.Insert("/hey")
	rt.Insert("/hell")
	rt.Insert("/fuck")
	rt.Insert("/htest")
	rt.Dump()

	successCount := 0
	st := time.Now()
	for i := 0; i < testCount; i++ {
		exist := rt.Find("/hello")
		if !exist {
			panic("impossible...radix-tree")
		}
		successCount++
	}
	et := time.Since(st).Nanoseconds()
	fmt.Printf("radix-tree cost: %v, successCount: %v\n", et, successCount)
}

func testMap() {
	m := make(map[string]bool)
	m["/golang"] = true
	m["/hello"] = true
	m["/hehe"] = true
	m["/hey"] = true
	m["/hell"] = true
	m["/fuck"] = true
	m["/htest"] = true

	successCount := 0
	st := time.Now()
	for i := 0; i < testCount; i++ {
		_, exist := m["/hello"]
		if !exist {
			panic("impossible...map")
		}
		successCount++
	}
	et := time.Since(st).Nanoseconds()
	fmt.Printf("map cost: %v, successCount: %v\n", et, successCount)
}

func main() {
	testRadixTree()
	testMap()
}
