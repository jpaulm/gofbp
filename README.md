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
- "subnets"- these are FBP networks where some of the connections are "sticky" - they can therefore act as (semi-) black box components

## Test Cases
The following test cases are now working - thanks to Egon Elbre for all his help!

- 2 Senders, one Receiver - merging first come, first served

- 2 Senders, with outputs concatenated using ConcatStr

- stream of IPs being distributed among several Receivers using RoundRobinSender 

- file being written to console  

- file being copied             

- file records being selected    

- test subnet (SubIn and SubOut)

- force deadlock (separate test file) - this is designed to crash, and in fact will give a message if it does *not* crash!
 

To run them, position to your `GitHub\gofbp` directory, and do any of the following:

- `go test -run Merge -count=1`
- `go test -run Concat -count=1`
- `go test -run RRDist -count=1`
- `go test -run CopyFile -count=1`
- `go test -run DoSelect1 -count=1`
- `go test -run DoSelect2 -count=1`  (with REJ not connected)
- `go test -run WriteToConsUsingNL -count=1`  (note the activated/deactivated messages)
- `go test -run ForceDeadlock -count=1`
- `go test -run TestSubnet -count=1`


**Note**: the last one in this list is constructed differently so that it can "crash" without disrupting the flow of tests: the network definition is in fact in `testdata`, while the test itself contains the code to compile and run the test.

You will occasionally see a message like `TempDir RemoveAll cleanup: remove ...\deadlock.exe: Access is denied.` - this is thought to be due to whatver AntiVirus software you are running.  I believe it can be ignored.

- `go test -count=1` runs them all, including `ForceDeadlock`

# Deadlocks

FBP deadlocks are well understood, and are handled well by other FBP implementations on https://github.com/jpaulm .  They also seem to be well detected by the Go scheduler - unfortunately, they are not so easy to troubleshoot, as Go deadlock detection is not "FBP-aware", and occurs before the GoFBP scheduler can analyze the process states to determine where the problem is occurring.  This has been raised as an issue - #28 .

To troubleshoot FBP deadlocks, look at the list of goroutines involved, and add the component names to your diagram, together with the "state".

As of this release (v2.2.1), a stand-alone program has been added, `analyze_deadlock.go`, which can be used to analyze the Go stack trace. Its `.exe` file can be found in the project `bin` directory.  Therefore, to analyze the deadlock, send the `go test` output for one test to `logfile`, i.e. `go test -run ForceDeadlock -count=1 > logfile`, then execute `bin\analyze_deadlock.exe`.  The output should be something like the following (based on running `go test -run ForceDeadlock -count=1`):

<pre>
Sender Goroutine no.: 19
Counter Goroutine no.: 20
Concat Goroutine no.: 21
Process: Sender, Status: Send
Process: Counter, Status: Send
Process: Concat, Status: Receive
</pre>

More information will be added, as we figured out how to extract it!

## Components

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
- subnet1.go   (this is a subnet, which can be treated as a component - we need to refine the packaging)

io folder:
- readfile.go
- writefile.go

**To dos**

- More and better documentation
- Convert `panic`s to more standard Go error handling
- Way too much logging - have to make that optional - use a JSON file...?  Issue raised for this... - done!
- Add subnet handling - done!
- Generate GoFBP networks from DrawFBP - https://github.com/jpaulm/drawfbp
- Add Load Balancing component
- Add sample code showing use of substreams
- "Automatic" ports
- Add Lua interface - see https://jpaulm.github.io/fbp/thlua.html
