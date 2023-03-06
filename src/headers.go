package main

import (
    _ "text/template"
)

type TOCEntry struct {
    Number []byte
    Text []byte
    Level []int
}

type CurrentNumber struct {
    Number []int
    Prev int
}

func newCurrentNumber() *CurrentNumber {
    return &CurrentNumber{
        Number: []int{1,1,1,1,1},
        Prev: 0,
    }
}

func (current *CurrentNumber) increment(new_level int) []int {

    if new_level <= current.Prev {
        current.Number[new_level-1] += 1
    }

    for i := new_level; i < len(current.Number); i++ {
        current.Number[i] = 1
    }

    current.Prev = new_level

    return current.Number[:new_level]

}

func (storage HeaderStorage) generateTOC() ([]byte, error) {

    current_number := newCurrentNumber()

    processed_entries := make([]TOCEntry, len(storage))

    for i, elem := range storage {
        
        header_number := current_number.increment(elem.Level)

        processed_number, err := serve(header_number, "{{ range . }}{{ printf \"%d\" . }}.{{ end }}")
        if err != nil {
            return nil, err
        }

        processed_entries[i] = TOCEntry{
            Number: processed_number,
            Text: elem.Text,
        }
    }

    return serve(processed_entries, `<div class="toc">{{ range . }}<p>{{ printf "%s" .Number }} - {{ printf "%s" .Text }}</p>
    {{ end }}</div>`)
}
