package main

import (
    "log"
    "os"
    "text/template"
    "bytes"
    "path/filepath"
)

type TemplateContents struct {
    Stylesheet []byte
    Title string
    Content string
}

func main() {
    var tmpl *template.Template
    var dest *os.File
    var err error

    tmpl, err = template.ParseFiles("./templates/main.html")
    if err != nil {
        log.Panicf("Can't read main template: %s", err)
    }

    var stylesheets []os.DirEntry
    var stylesheet_combined [][]byte = [][]byte{}
    stylesheets, err = os.ReadDir("./css")
    for _, stylesheet := range stylesheets {
        var stylesheet_contents []byte
        stylesheet_contents, err = os.ReadFile(filepath.Join("./css",stylesheet.Name()))
        if err != nil {
            log.Printf("Error with %s: %s", stylesheet.Name(), err)
            continue
        }
        stylesheet_combined = append(stylesheet_combined, stylesheet_contents)
    }

    var contents *TemplateContents = &TemplateContents{
        Stylesheet: bytes.Join(stylesheet_combined, []byte{byte('\n')}),
        Title: "Penis",
        Content: "Quality content",
    }
    log.Printf("%s", bytes.Join(stylesheet_combined, []byte{byte('\n')}))

    dest, err = os.Create("./output/index.html")
    if err != nil {
        log.Panicf("Can't create output file: %s", err)
    }
    defer dest.Close()

    err = tmpl.Execute(dest, contents)
    if err != nil {
        log.Panicf("Can't write output file: %s", err)
    }

    log.Printf("Done")
}
