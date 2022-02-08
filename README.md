# GoFBP 

This repo contains an FBP implementation in Go.  It conforms pretty closely to the scheduling logic in JavaFBP, C#FBP and C++FBP (all on https://github.com/jpaulm ).

There may well be further internal changes, but I am hoping that the "external" APIs (network and component definitions) are now firm. 

Current tag: **v1.0.8** .

## General
 
General web site for "classical" FBP: 
* https://jpaulm.github.io/fbp

In computer programming, flow-based programming (FBP) is a programming paradigm that defines applications as networks of "black box" processes, which exchange data across predefined connections by message passing, where the connections are specified externally to the processes. These black box processes can be reconnected endlessly to form different applications without having to be changed internally. FBP is thus naturally component-oriented.

FBP is a particular form of dataflow programming based on bounded buffers, information packets with defined lifetimes, named ports, and separate definition of connections.

DrawFBP (https://github.com/jpaulm/drawfbp) can now generate working GoFBP network and subnet definitions from a network diagram (as well as JavaFBP, C#FBP, JSON and free-form notations).
 
## GoFBP Network Definition Syntax and Component API:
* https://jpaulm.github.io/fbp/gosyntax.htm

Please note the addition of 

```
params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
 ```
 
 and the associated ` net.SetParams(params)` after the call to `NewNetwork`.

## Features (these are common to all FBP implementations on GitHub/jpaulm):

- delayed start of goroutines (FBP processes), unless `MustRun` attribute is specified or the process has no non-IIP inputs (same as JavaFBP delayed start feature) 
- the reason for `MustRun` is that components are not triggered if there is no data incoming on their non-IIP input ports (apart from closing down downstream processes as appropriate);  some components however need to execute in spite of this, e.g. `components\io\writefile.go` (which must clear the output file), and counter-type components.
- optional output ports - see `components\testrtn\writetoconsole.go`
- "subnets"- these are FBP networks where some of the connections are "sticky" - they can therefore act as (semi-) black box components
- "automatic" in- and out-ports - notation is port name = "*"
- GoFBP (and the other FBP implementations on https://github.com/jpaulm distingish between "invocation" and "activation" of processes: a process is invoked once, but may be activated multiple times (if it does a return before all its input ports have been drained) 

## Running your app with GoFBP

In the `go.mod` file in your root, add the statement

```
require github.com/jpaulm/gofbp latest
```
and then run the command `go mod tidy` - this will change the word `latest` to the latest version, and store it back in your `go.mod` file.

If you need parameter values other than `false`, you will have to access a `params.xml` file, as shown below.

Command to run your network, e.g.:

```
go run Merge.go
```

A number of test cases are in the `testing` folder, and can be run as follows:

```
cd testing
go test
```
To run an individual network, e.g. `TestIntSender` in `testing\testintsender_test.go`,
run
```
cd testing
go test -run IntSender
```

### Notes:

- One comment about running with subnets: the subnet or subnets should be in a different folder from the code that invokes it.  If they are in the same folder, apparently both the invoking code and the subnet need to be named individually in a `go run` command... this is not required when running under VSCode, or from an `.exe` file.  

## Tracing

An XML file has been provided in the root, called `params.xml`.  If you need values other than all `false`, you will need to reference this using the `LoadXMLParams` and `SetParams` methods - see https://jpaulm.github.io/fbp/gosyntax.htm .

Format of the tracing definitions file:

```
<?xml version="1.0"?> 
<runparams>
<tracing>false</tracing>
<tracelocks>false</tracelocks>
<tracepkts>false</tracepkts>     (traces packet creates, createbrackets and discards)
<generate-gIds>true</generate-gIds> 
</runparams>
```
The `generate-gIds` parameter is only used to assist in debugging deadlocks, so can default most of the time - see below.

## Websockets support

The GoFBP WebSockets support works very similarly to the corresponding `JavaFBP` facilities - see https://github.com/jpaulm/javafbp-websockets#readme - esp. the diagrams at the end.

For information about how to use the `GoFBP` Websockets components, see https://jpaulm.github.io/fbp/gofbp_websockets.md .  In this test case, we use HTML and JavaScript as the client. 

#### Testing WebSockets support:

- position to your `GitHub\gofbp\testing` directory in DOS
- update `params.xml` to set tracing parameters as desired
- run `go test -run TestWebSocket -count=1`
- in File Explorer, locate `GitHub\gofbp\scripts\chat2.html`
- open with Firefox or Chrome (make sure you don't open the browser before starting the `gofbp` app)
- enter `namelist` in the Command box, hit enter or Send
- you should see 
```
    Server: Line1
    Server: Line2
    Server: Line3
```
show up in the box below "Send"
- click on `Stop WS` - you will see `End of dialog` next to `Status`

- you're done... 

## Subnets
These are described in Chap. 7 in "Flow-Based Programming": Composite Components - https://jpaulmorrison.com/fbp/compos.shtml , although we haven't implemented dynamic subnets (yet)!

**Note:** If a subnet `.go` file is in the same folder as the network `.go` file invoking it, the `NewProc` call should be written without the folder - e.g. `net.NewProc("Run___Subnet", &Subnet1{})`, but currently DrawFBP cannot detect this situation, so, when using notation `GoFBP`, this call will not be generated correctly. 


## Other Test Cases
The following test cases are now working - thanks to Egon Elbre and Emil Valeev for all their help!

- 2 Senders, one Receiver - merging first come, first served

- 2 Senders, with outputs concatenated using ConcatStr

- stream of IPs being distributed among several Receivers using RoundRobinSender 

- file being written to console  

- file being copied             

- file records being selected    

- test subnet (SubIn and SubOut)

- force deadlock (separate test file) - this is designed to crash, and in fact will give a message if it does *not* crash!

- simple web sockets test
 

To run them, position to your `GitHub\gofbp\testing` directory (note the additional folder, as of tag `v1.0.3` plus), and do any of the following:

- `go test -run Merge -count=1`
- `go test -run Concat -count=1`
- `go test -run RRDist -count=1`
- `go test -run CopyFile -count=1`
- `go test -run DoSelect1 -count=1`
- `go test -run DoSelect2 -count=1`  (with REJ not connected)
- `go test -run WriteToConsUsingNL -count=1`  (note the activated/deactivated messages)
- `go test -run ForceDeadlock -count=1`
- `go test -run InfQueueAsMain -count=1` (note the "automatic" ports between WriteFile and ReadFile)
- `go test -run LoadBal -count=1` (this does load balancing, and uses two DelayedReceiver processes)
- `go test -run Subnet1 -count=1` 
- `go test -run Subnet2 -count=1` 
- `go test -run Subnet3 -count=1` 
- `go test -run TestWebSocket -count=1`  - for more instructions, see # Test WebSockets (below)


**Note**: ForceDeadlock is constructed differently so that it can "crash" without disrupting the flow of tests: the network definition has to be compiled "on the fly", so it is actually in `testdata`, while the test itself contains the code to compile and run the test.

You will occasionally see a message like `TempDir RemoveAll cleanup: remove ...\deadlock.exe: Access is denied.` - this is thought to be due to whatver AntiVirus software you are running.  I believe it can be ignored.

- `go test -count=1` runs them all, including `ForceDeadlock`


## Deadlocks

FBP deadlocks are well understood, and are handled well by other FBP implementations on https://github.com/jpaulm .  They also seem to be well detected by the Go scheduler - unfortunately, they are not so easy to troubleshoot, as Go deadlock detection is not "FBP-aware", and occurs before the GoFBP scheduler can analyze the process states to determine where the problem is occurring.  This has been raised as an issue - #28 .

As of release v2.2.1, a stand-alone program has been added, `analyze_deadlock.go`, which can be used to analyze the Go stack trace. Its `.exe` file can be found in the assets of the latest release.  To analyze the deadlock, change the `generate-gIds` parameter in `params.xml` to `true`, then send the `go test` output to `logfile`, i.e. `go test -run ForceDeadlock -timeout 0ms > logfile`.  Now download `analyze_deadlock.exe` to some convenient folder, then execute it.  The output should be something like the following (based on running `go test -run ForceDeadlock -count=1`):

<pre>
Sender Goroutine no.: 19
Counter Goroutine no.: 20
Concat Goroutine no.: 21
Process: Sender, Status: Send
Process: Counter, Status: Send
Process: Concat, Status: Receive
</pre>

Now look at the list of goroutines involved, and add the component names to your diagram, together with the displayed "status".  Typically the deadlock will be "between" some goroutines waiting to Send and some waiting to Receive.


In MS-DOS, you can do the above all on one line, as follows:

<pre>
go test -run ForceDeadlock -count=1 > logfile & analyze_deadlock.exe
</pre>

(Not sure if you can do this with PowerShell...?)

## Components

The following components are available:

"testrtn" folder:
- `concatstr.go`
- `discard.go`
- `kick.go`        
- `receiver.go`
- `roundrobinsender.go`
- `selector.go`
- `sender.go`
- `writetoconsole.go` 
- `writetoconsNL.go`   (same, but written as a non-looper)

"subnets" folder:
- `subnet1.go`   (this is a subnet, i.e. a "network" with "sticky" connections - this can be treated as a component)
- `sssubnet1.go`   (this is similar to subnet1.go, but with a substream-sensitive front-end, and a substream delimiter generating back-end)
- `sssubnet2.go`   (this is the same as sssubnet1, but with a counter emitting the size of each incoming substream)

"io" folder:
- `readfile.go`
- `writefile.go`

"websocket" folder (for instructions, see below):
- `ws_request.go`
- `ws_respond.go`
- `ws_ans_req.go`  (sample component - properly belongs in test suite)

**To dos**

- Convert `panic`s to more standard Go error handling
- Way too much logging - have to make that optional - put remaining logging under switch control - *done!*
- Add subnet handling - *done!*
- Generate GoFBP networks from DrawFBP - https://github.com/jpaulm/drawfbp - *done!*
- Add Load Balancing component - *done!*
- Add sample code showing use of substreams - *done!*
- "Automatic" ports - *done!*
- Add Lua interface - similar to https://jpaulm.github.io/fbp/thlua.html 

