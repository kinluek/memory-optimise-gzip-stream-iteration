# gzipped

Package used to read and split a stream of concatenated gzip files, no unzipping involved. 

## Goals

 - To iterate through a stream of gzip files in a memory efficient manner.
 
 
## Current Benchmark

 -  Scanner allocates < 2MB to iterate through a stream of 1000 1MB+ files, one file at a time.