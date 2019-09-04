[![Build Status](https://travis-ci.org/nishaad78/lwwset.svg?branch=master)](https://travis-ci.org/nishaad78/lwwset)
[![codecov](https://codecov.io/gh/nishaad78/lwwset/branch/master/graph/badge.svg)](https://codecov.io/gh/nishaad78/lwwset)
[![Go Report Card](https://goreportcard.com/badge/github.com/nishaad78/lwwset)](https://goreportcard.com/report/github.com/nishaad78/lwwset)
[![Go doc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square
)](https://godoc.org/github.com/nishaad78/lwwset)

# lwwset
This is a thread safe implementation of a [Last-Write-Wins-Element-Set](https://en.wikipedia.org/wiki/Conflict-free_replicated_data_type#LWW-Element-Set_(Last-Write-Wins-Element-Set)) with a bias towards removals.

# Example usage
```go
s := lwwset.NewLWW()

s.Add('a', time.Now())

t, ok := s.Lookup('a')
if ok {
    fmt.Println("we found it at ", t)
}
```