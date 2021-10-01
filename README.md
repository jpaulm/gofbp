# gofbp 

This repo holds the beginning of an FBP implementation in Go

*Warning!  As of 30 Sept., 2021, this release is a bit buggy - please do not use for now!* 

Features include:

- delayed start of goroutines (FBP processes), unless `MustRun` attribute is specified or the process has no non-IIP inputs (same as JavaFBP delayed start feature) 
- the reason for `MustRun` is that components are not triggered if there is no data incoming on their input ports (apart from closing down downstream processes as appropriate;  some components however need to execute in spite of this, e.g. WriteFile, and counting components.
- optional output ports - see https://github.com/jpaulm/gofbp/blob/master/components/testrtn/writetoconsole.go


The following test cases are now working - thanks to Egon Elbre for all his help!

- 2 Senders, one Receiver - merging first come, first served

- 2 Senders, with outputs concatenated using ConcatStr

- stream of IPs being distributed among several Receivers using RoundRobinSender 

- file being written to console  

- file being copied             

- file records being selected    

- force deadlock (https://github.com/jpaulm/gofbp/blob/master/main11.go - change to main.go to run)
 

To run them, position to your `GitHub\gofbp` directory, and do any of the following:

- `go test -run Merge -count=1`
- `go test -run Concat -count=1`
- `go test -run RRDist -count=1`
- `go test -run CopyFile -count=1`
- `go test -run DoSelect -count=1`
- `go test -run WriteToConsUsingNL -count=1`  (note the activated/deactivated messages)

`go test` runs them all, in sequence - _note that the last test in this sequence is designed to crash!_

The following components are available:

testrtn folder:
- concatstr.go
- discard.go
- receiver.go
- roundrobinsender.go
- selector.go
- sender.go
- writetoconsole.go 
- writetoconsNL.go   (same, but written as a non-looper)

io folder:
- readfile.go
- writefile.go

Note: way too much logging - have to make that optional - use a JSON file...?  Issue raised for this...
