package main

import (
	_ "text/template"
)

type TOCEntry struct {
	Number []byte
	Text   []byte
	Level  int
	Link   string
}

type CurrentHeaderNumberStorage struct {
	Number []int
	Prev   int
}

func newCurrentNumber(max_levels int) *CurrentHeaderNumberStorage {

	initialized_storage := make([]int, max_levels)

	for i, _ := range initialized_storage {
		initialized_storage[i] = 1
	}

	return &CurrentHeaderNumberStorage{
		Number: initialized_storage,
		Prev:   0,
	}
}

func (current *CurrentHeaderNumberStorage) increment(new_level int) []int {

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

	current_number := newCurrentNumber(4)

	processed_entries := make([]TOCEntry, len(storage))

	for i, elem := range storage {

		header_number := current_number.increment(elem.Level)

		processed_number, err := serve(header_number, "{{ range . }}{{ printf \"%d\" . }}.{{ end }}")
		if err != nil {
			return nil, err
		}

		processed_entries[i] = TOCEntry{
			Level:  elem.Level,
			Number: processed_number,
			Text:   elem.Text,
			Link:   elem.ID,
		}
	}

	return serve(processed_entries, "<div class=\"toc\">{{ range . }}<p class=\"level{{ printf \"%d\" .Level }}\">{{ printf \"%s\" .Number }}\t\t<a href=\"#{{ .Link }}\">{{ printf \"%s\" .Text }}</a></p>\n{{ end }}</div>")
}
