package main

import (
	"flag"
	"fmt"
	"github.com/mirrorweb/gzipped"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	f, err := makeConcatFile()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer f.Close()

	if err := program(f); err != nil {
		log.Fatalf("%v", err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

// makeConcatFile crates a file that acts like a concatenation of multiple gzip files.
func makeConcatFile() (*os.File, error) {
	f, err := os.Create("test-concat-files.gz")
	if err != nil {
		return nil, err
	}

	for i := 0; i < 1000; i++ {
		filler := strings.Repeat("a", 1024*1024)  // create data in the loop to cause some allocations for the profiler to pick up.
		fillerBytes := []byte(filler)
		if _, err := f.Write([]byte{0x00, 0x00, 0x1f, 0x8b}); err != nil {
			return nil, err
		}
		if _, err := f.Write(fillerBytes); err != nil {
			return nil, err
		}
	}
	f.Seek(0, 0)
	return f, nil
}

func program(f *os.File) error {
	scanner := gzipped.NewScanner(f, 1024)
	for i := 1; scanner.Scan(); i++ {
		scannedFile := scanner.FileBytes()
		fmt.Println(i, len(scannedFile))
	}
	if scanner.Err() != nil {
		return fmt.Errorf("scanning failed: %s", scanner.Err())
	}
	return nil
}
