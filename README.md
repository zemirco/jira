
# jira

[![Build Status](https://travis-ci.org/zemirco/jira.svg)](https://travis-ci.org/zemirco/jira)
[![GoDoc](https://godoc.org/github.com/zemirco/jira?status.svg)](https://godoc.org/github.com/zemirco/jira)

JIRA REST API client in Go.

## Example

```go
package main

import "github.com/zemirco/jira"

func main() {
  jira := New("https://jira.atlassian.com/")
  rapidViews, err := jira.RapidViews()
  if err != nil {
    panic(err)
  }
  fmt.Printf("%+v", rapidViews)
}
```

## Test

`go test`

## License

MIT
