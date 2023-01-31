package main

import (
    "log"
    "os"
    "html/template"
)

func main() {
    var tmpl *template.Template
    var dest *os.File
    var err error

    tmpl, err = template.ParseFiles("./templates/main.html")
    if err != nil {
        log.Panicf("Can't read main template: %s", err)
    }

    dest, err = os.Create("./output/index.html")
    if err != nil {
        log.Panicf("Can't create output file: %s", err)
    }
    defer dest.Close()

    err = tmpl.Execute(dest, nil)
    if err != nil {
        log.Panicf("Can't write output file: %s", err)
    }

    log.Printf("Done")
}
