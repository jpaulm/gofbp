# gofbp 

This repo holds the beginning of an FBP implementation in Go

*As of 17 Oct., 2021, all test cases are working again.*

*Apologies: if you downloaded GoFBP before 17 Oct., please download it again, and use the current component definitions as a model for developing your own components*

There may well be further internal changes, but I am hoping that the "external" APIs (network and component definitions) are now firm. 

## General
 
General web site for "classical" FBP: 
* https://jpaulm.github.io/fbp

In computer programming, flow-based programming (FBP) is a programming paradigm that defines applications as networks of "black box" processes, which exchange data across predefined connections by message passing, where the connections are specified externally to the processes. These black box processes can be reconnected endlessly to form different applications without having to be changed internally. FBP is thus naturally component-oriented.

FBP is a particular form of dataflow programming based on bounded buffers, information packets with defined lifetimes, named ports, and separate definition of connections.
 
GoFBP Network Definition Syntax and Component API:
* https://jpaulm.github.io/fbp/gosyntax.htm

## Features (these are common to all FBP implementations on GitHub/jpaulm):

- delayed start of goroutines (FBP processes), unless `MustRun` attribute is specified or the process has no non-IIP inputs (same as JavaFBP delayed start feature) 
- the reason for `MustRun` is that components are not triggered if there is no data incoming on their non-IIP input ports (apart from closing down downstream processes as appropriate);  some components however need to execute in spite of this, e.g. `components\io\writefile.go` (which must clear the output file), and counter-type components.
- optional output ports - see `components\testrtn\writetoconsole.go`

The following test cases are now working - thanks to Egon Elbre for all his help!

- 2 Senders, one Receiver - merging first come, first served

- 2 Senders, with outputs concatenated using ConcatStr

- stream of IPs being distributed among several Receivers using RoundRobinSender 

- file being written to console  

- file being copied             

- file records being selected    

- force deadlock (separate test file) - this is designed to crash, and in fact will give a message if it does *not* crash!
 

To run them, position to your `GitHub\gofbp` directory, and do any of the following:

- `go test -run Merge -count=1`
- `go test -run Concat -count=1`
- `go test -run RRDist -count=1`
- `go test -run CopyFile -count=1`
- `go test -run DoSelect -count=1`
- `go test -run WriteToConsUsingNL -count=1`  (note the activated/deactivated messages)
- `go test -run ForceDeadlock -count=1`

- `go test -count=1` runs them all, including `ForceDeadlock` (as the first one)

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
- Add Load Balancing component
- Add sample code showing use of substreams
