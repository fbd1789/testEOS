package parser

import (
    "gopkg.in/yaml.v3"
    "os"
)

type TestCase struct {
    Name     string                 `yaml:"name"`
    Command  string                 `yaml:"command"`
    Template string                 `yaml:"template"`
    Vars     map[string]interface{} `yaml:"vars"`
}

type testFile struct {
    Tests []TestCase `yaml:"tests"`
}

func LoadTests(path string) []TestCase {
    data, err := os.ReadFile(path)
    if err != nil {
        panic(err)
    }

    var tf testFile
    if err := yaml.Unmarshal(data, &tf); err != nil {
        panic(err)
    }

    return tf.Tests
}
