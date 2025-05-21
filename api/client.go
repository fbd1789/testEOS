package api

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type EOSClient interface {
    Execute(command string) (map[string]interface{}, error)
}

type RealClient struct {
    Host     string
    Username string
    Password string
}

func NewRealClient(host, username, password string) *RealClient {
    return &RealClient{
        Host:     host,
        Username: username,
        Password: password,
    }
}

func (c *RealClient) Execute(command string) (map[string]interface{}, error) {
    url := fmt.Sprintf("https://%s/command-api", c.Host)

    payload := map[string]interface{}{
        "jsonrpc": "2.0",
        "method":  "runCmds",
        "params": map[string]interface{}{
            "version": 1,
            "cmds":    []string{command},
            "format":  "json",
        },
        "id": "1",
    }

    data, _ := json.Marshal(payload)

    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
    req.SetBasicAuth(c.Username, c.Password)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    if errData, ok := result["error"]; ok {
        return nil, fmt.Errorf("API error: %v", errData)
    }

    resList := result["result"].([]interface{})
    raw := resList[0].(map[string]interface{})
// 🔧 Corriger powerSupplies: convertir clés en string
if psRaw, ok := raw["powerSupplies"].(map[interface{}]interface{}); ok {
    fixed := make(map[string]interface{})
    for k, v := range psRaw {
        keyStr := fmt.Sprintf("%v", k)
        fixed[keyStr] = v
    }
    raw["powerSupplies"] = fixed
} else if psRaw, ok := raw["powerSupplies"].(map[interface{}]any); ok {
    fixed := make(map[string]interface{})
    for k, v := range psRaw {
        fixed[fmt.Sprintf("%v", k)] = v
    }
    raw["powerSupplies"] = fixed
} else if psRaw, ok := raw["powerSupplies"].(map[int]interface{}); ok {
    fixed := make(map[string]interface{})
    for k, v := range psRaw {
        fixed[fmt.Sprintf("%d", k)] = v
    }
    raw["powerSupplies"] = fixed
} else if psRaw, ok := raw["powerSupplies"].(map[float64]interface{}); ok {
    fixed := make(map[string]interface{})
    for k, v := range psRaw {
        fixed[fmt.Sprintf("%.0f", k)] = v
    }
    raw["powerSupplies"] = fixed
}
    return raw, nil
}
