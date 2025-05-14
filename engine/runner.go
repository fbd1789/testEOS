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
}) {
    var wg sync.WaitGroup
    results := make(chan TestResult, len(tests))

    for _, test := range tests {
        wg.Add(1)
        go func(t parser.TestCase) {
            defer wg.Done()
            jsonResp, err := client.Execute(t.Command)
            if err != nil {
                results <- TestResult{Name: t.Name, Success: false, Error: err.Error()}
                return
            }
            // fmt.Printf("DEBUG [%s] response JSON:\n%+v\n", t.Name, jsonResp) // Pour analyse du JSON venant du device
            ok, msg, err := Validate(jsonResp, t.Template, t.Vars)
            if err != nil {
                results <- TestResult{Name: t.Name, Success: false, Error: err.Error()}
            } else {
                results <- TestResult{Name: t.Name, Success: ok, Error: msg}
            }
        }(test)
    }

    wg.Wait()
    close(results)

    for r := range results {
        fmt.Println(r.Error)
    }
}
