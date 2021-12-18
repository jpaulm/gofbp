package main

// Analyze deadlock output

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//var rec string
	procs := make(map[string]string)

	dat, err := os.ReadFile("logfile")
	check(err)
	//fmt.Print(string(dat))
	rec := string(dat[:])
	rec = strings.ReplaceAll(rec, string('\x00'), "")
	i := 0
	var gName string
	for {
		j := strings.Index(rec[i:], "Goroutine ")
		if j != -1 {
			// Goroutine Sender: no. 19
			i += j
			j := strings.Index(rec[i+9:], ":")
			s := rec[i+10 : i+10+j-1]
			i += 10 + j + 5
			j = strings.Index(rec[i:], "\r")
			if j == -1 {
				j = strings.Index(rec[i:], "\n")
			}
			procs[rec[i:i+j]] = s
			fmt.Println(s, "Goroutine no.:", rec[i:i+j])
			i += 16 + j
		} else {
			j = strings.Index(rec[i:], "goroutine ")
			if j == -1 {
				return
			}
			i += j
			// goroutine 19 [sync.Cond.Wait]:
			k := strings.Index(rec[i+10:], " ")
			//t := rec[i : i+30]
			//_ = t
			s := rec[i+10 : i+10+k]
			i += 10 + k
			gName = procs[s]
			if gName == "" {
				continue
			}
			j = strings.Index(rec[i:], "[sync.Cond.Wait]:")
			i += j
			// github.com/jpaulm/gofbp/core.(*Process).Send(...)
			k = strings.Index(rec[i:], "(*Process).")
			l := strings.Index(rec[i+k+11:], "(")
			fmt.Println("Process:", gName+",", "Status:", rec[i+k+11:i+k+11+l])
			i += k + 11
		}
	}
}
