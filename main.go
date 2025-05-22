package main

import (
    "bufio"
    "eosTester/api"
    "eosTester/engine"
    "eosTester/parser"
    "flag"
    "fmt"
    "os"
    "strings"
    "sync"
)

func main() {
    hostFile := flag.String("hosts", "hosts.txt", "Path to file with EOS device list")
    user := flag.String("user", "", "Username for EOS")
    pass := flag.String("pass", "", "Password for EOS")
    yml := flag.String("file", "config/test.yml", "Path to test YAML file")
    selected := flag.String("s", "", "Name of a specific test to run (optional)")
    flag.Parse()

    hosts, err := loadHosts(*hostFile)
    if err != nil {
        fmt.Printf("Failed to load hosts: %v\n", err)
        return
    }

    allTests := parser.LoadTests(*yml)
    var tests []parser.TestCase
    if *selected != "" {
        for _, t := range allTests {
            if t.Name == *selected {
                tests = append(tests, t)
            }
        }
        if len(tests) == 0 {
            fmt.Printf("No test found with name '%s'\n", *selected)
            return
        }
    } else {
        tests = allTests
    }

    var wg sync.WaitGroup
    for _, host := range hosts {
        wg.Add(1)
        go func(h string) {
            defer wg.Done()
            fmt.Printf("🔍 Running tests on host: %s\n", h)
            client := api.NewRealClient(h, *user, *pass)
            engine.RunTests(tests, client, host)
        }(host)
    }
    wg.Wait()
}

func loadHosts(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var hosts []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line != "" && !strings.HasPrefix(line, "#") {
            hosts = append(hosts, line)
        }
    }
    return hosts, scanner.Err()
}
