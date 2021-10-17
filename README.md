# gofbp 

This repo holds the beginning of an FBP implementation in Go

*As of 17 Oct., 2021, all test cases are working again.*

*Apologies: if you downloaded GoFBP before 17 Oct., please download it again, and use the current component definitions as a model for developing your own components*

There may well be further internal changes, but I am hoping that the "external" APIs (network and component definitions) are now firm. 

Features include:

- delayed start of goroutines (FBP processes), unless `MustRun` attribute is specified or the process has no non-IIP inputs (same as JavaFBP delayed start feature) 
- the reason for `MustRun` is that components are not triggered if there is no data incoming on their non-IIP input ports (apart from closing down downstream processes as appropriate);  some components however need to execute in spite of this, e.g. `components\io\writefile.go` (which must clear the output file), and counter-type components.
- optional output ports - see `components\testrtn\writetoconsole.go`

**Note:** the last test in `go test` is designed to crash.  To skip this test, remove the double slashes from the `t.Skip` statement at the beginning of `force_deadlock_test.go`.

The following test cases are now working - thanks to Egon Elbre for all his help!

- 2 Senders, one Receiver - merging first come, first served

- 2 Senders, with outputs concatenated using ConcatStr

- stream of IPs being distributed among several Receivers using RoundRobinSender 

- file being written to console  

- file being copied             

- file records being selected    

- force deadlock (last test in `gofbp_test.go`) - this is designed to crash!
 

To run them, position to your `GitHub\gofbp` directory, and do any of the following:

- `go test -run Merge -count=1`
- `go test -run Concat -count=1`
- `go test -run RRDist -count=1`
- `go test -run CopyFile -count=1`
- `go test -run DoSelect -count=1`
- `go test -run WriteToConsUsingNL -count=1`  (note the activated/deactivated messages)

To run ForceDeadlock, do the following:

- `go test -run ForceDeadlock -count=1 -timeout 0s` 

`go test -count=1` runs them all, including `ForceDeadlock`

Note that the last test in this sequence (`ForceDeadlock`) is designed to crash - you will probably want to run it by itself...

The following components are available:

testrtn folder:
- concatstr.go
- discard.go
- kick.go
- receiver.go
- roundrobinsender.go
- selector.go
- sender.go
- writetoconsole.go 
- writetoconsNL.go   (same, but written as a non-looper)

io folder:
- readfile.go
- writefile.go

**To dos**

- More and better documentation
- Convert `panic`s to more standard Go error handling
- Way too much logging - have to make that optional - use a JSON file...?  Issue raised for this...
- Add subnet handling
- Generate GoFBP networks from DrawFBP - https://github.com/jpaulm/drawfbp
