package core

import (
	"bytes"
	"fmt"
	"runtime/pprof"
	"sort"
	"strings"

	"github.com/google/pprof/profile"
)

// goroutineTrace returns summary of the goroutines.
func (n *Network) goroutineTrace() (string, error) {
	var pb bytes.Buffer
	profiler := pprof.Lookup("goroutine")
	if profiler == nil {
		return "", fmt.Errorf("unable to find profile")
	}
	err := profiler.WriteTo(&pb, 0)
	if err != nil {
		return "", fmt.Errorf("failed to write profile: %w", err)
	}

	p, err := profile.ParseData(pb.Bytes())
	if err != nil {
		return "", fmt.Errorf("failed to parse profile: %w", err)
	}

	return n.summarizeProfile(p, n.id())
}

func (n *Network) summarizeProfile(p *profile.Profile, networkID string) (string, error) {
	var b strings.Builder

	for _, sample := range p.Sample {
		if !existsInSlice(sample.Label["network"], networkID) {
			continue
		}

		fmt.Fprintf(&b, "count %d @", sample.Value[0])

		// stack trace summary

		if len(sample.Label)+len(sample.NumLabel) > 0 {
			if len(sample.Label) > 0 {
				keys := []string{}
				for k := range sample.Label {
					if k == "network" {
						continue
					}
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					values := sample.Label[k]
					fmt.Fprintf(&b, " %s:", k)
					switch len(values) {
					case 0:
					case 1:
						fmt.Fprintf(&b, "%q", values[0])
					default:
						fmt.Fprintf(&b, "%q", values)
					}
				}
			}
			if len(sample.NumLabel) > 0 {
				keys := []string{}
				for k := range sample.NumLabel {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintf(&b, "%s:%v", k, sample.NumLabel[k])
				}
			}
		}
		fmt.Fprintf(&b, "\n")

		// each line
		for _, loc := range sample.Location {
			for i, ln := range loc.Line {
				if i == 0 {
					fmt.Fprintf(&b, "#   %#8x", loc.Address)
					if loc.IsFolded {
						fmt.Fprint(&b, " [F]")
					}
				} else {
					fmt.Fprint(&b, "#           ")
				}
				if fn := ln.Function; fn != nil {
					fmt.Fprintf(&b, " %-50s %s:%d", fn.Name, fn.Filename, ln.Line)
				} else {
					fmt.Fprintf(&b, " ???")
				}
				fmt.Fprintf(&b, "\n")
			}
		}
		fmt.Fprintf(&b, "\n")
	}
	return b.String(), nil
}

func existsInSlice(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
