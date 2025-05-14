package engine

import (
    "bytes"
    "text/template"
    // "os"
    "strings"

    "github.com/Masterminds/sprig/v3"
)

func Validate(jsonData map[string]interface{}, templateFile string, vars map[string]interface{}) (bool, string, error) {
    // Fusion des fonctions : sprig + custom
    funcMap := sprig.TxtFuncMap()
    funcMap["compareVersions"] = CompareVersions

    tmpl, err := template.New(templateFile).Funcs(funcMap).ParseFiles("templates/" + templateFile)
    if err != nil {
        return false, "", err
    }

    data := map[string]interface{}{
        "result": jsonData,
        "vars":   vars,
    }

    var output bytes.Buffer
    err = tmpl.Execute(&output, data)
    if err != nil {
        return false, "", err
    }

//     fmt.Println("Template:", templateFile)
// fmt.Println("Data keys passed:", data)
// fmt.Println("TEMPLATE OUTPUT:\n" + output.String()) // <== debug ici
    if strings.Contains(output.String(), "FAIL") {
        return false, output.String(), nil
    }
    return true, output.String(), nil
}

