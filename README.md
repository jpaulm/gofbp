# gofbp 

This repo is for early testing of GoFBP ideas and trial balloons! 


Three test cases:

- 2 Senders, one Receiver - merging first come, first served

- 2 Senders, with outputs concatenated using ConcatStr

- stream of IPs being distributed among several Receivers using RoundRobinSender - **currently there is a deadlock in this test case**


Note: way too much logging - have to make that optional - use a JSON file...?
