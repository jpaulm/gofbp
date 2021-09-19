# gofbp 

This repo is for early testing of GoFBP ideas and trial balloons! 


Three test cases:

- 2 Senders, one Receiver - merging first come, first served

- 2 Senders, with outputs concatenated using ConcatStr

- stream of IPs being distributed among several Receivers using RoundRobinSender 

To run them:

- `go test -run Merge -count=1`
- `go test -run Concat -count=1`
- `go test -run RRDist -count=1`


Note: way too much logging - have to make that optional - use a JSON file...?  Issue raised for this...
