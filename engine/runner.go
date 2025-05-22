package engine

import (
    "eosTester/parser"
    "fmt"
    "sync"
)

type TestResult struct {
    Name    string
    Success bool
    Error   string
}

func RunTests(tests []parser.TestCase, client interface {
    Execute(string) (map[string]interface{}, error)
}, host string) {
    var wg sync.WaitGroup
    results := make(chan TestResult, len(tests))

    for _, test := range tests {
        wg.Add(1)
        go func(t parser.TestCase, h string) {
            defer wg.Done()
            jsonResp, err := client.Execute(t.Command)
            if err != nil {
                results <- TestResult{Name: t.Name, Success: false, Error: err.Error()}
                return
            }
            ok, msg, err := Validate(jsonResp, t.Template, t.Vars, h)
            if err != nil {
                results <- TestResult{Name: t.Name, Success: false, Error: err.Error()}
            } else {
                results <- TestResult{Name: t.Name, Success: ok, Error: msg}
            }
        }(test, host)
    }

    wg.Wait()
    close(results)

    for r := range results {
        fmt.Println(r.Error)
    }
}
