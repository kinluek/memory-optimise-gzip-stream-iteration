# gzipped

Package used to read and split a stream of concatenated gzip files, no unzipping involved. 

## Run 

To proifle the example application, run `make`, this will build and run the program and then throw you into the pprof shell.

While in the `pprof` shell, run the command `list main.program` this will list the memory allocations for program function found in 
`cmd/profile/main.go` 

## Goals

 - To iterate through a stream of gzip files in a memory efficient manner.
 
 
## Current Benchmark

MacBook Pro (Retina, 15-inch, Mid 2015) - 2.2 GHz Quad-Core Intel Core i7 - 16 GB 1600 MHz DDR3

 -  Scanner allocates < 5MB to iterate through a stream of 1000 1MB+ files, one file at a time.