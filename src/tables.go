package main

import (
	"fmt"
)

type processedTable struct {
	Title []byte
	// TODO: add table description header row
	// Description []byte
	ID   []byte
	Rows []*processedTableRow
}

type processedTableRow struct {
	Images     []*processedTableImage
	Paragraphs [][]byte
}

type processedTableImage struct {
	Path   []byte
	IsWide bool
}

func (table Table) Serve(images *ImageStorage, ids *IdStorage, links *LinkStorage) ([]byte, error) {

	processed_table := new(processedTable)
	processed_table.Rows = make([]*processedTableRow, len(table.Rows))

	id, err := ids.push(table.ID)
	if err != nil {
		return nil, fmt.Errorf("Table serving error: %s", err)
	}

	processed_table.ID = id

	for i, row := range table.Rows {

		processed_table.Rows[i] = &processedTableRow{}

		processed_table.Rows[i].Images = make([]*processedTableImage, len(row.Paths))

		for j, path := range row.Paths {

			content, err := path.Serve(images)
			if err != nil {
				return nil, err
			}

			is_wide, err := isWide(content)
			if err != nil {
				return nil, fmt.Errorf("Invalid image %s. Error: %s\n", content, err)
			}

			processed_table.Rows[i].Images[j] = &processedTableImage{
				Path:   content,
				IsWide: is_wide,
			}

		}

		processed_table.Rows[i].Paragraphs = make([][]byte, len(row.Paragraphs))

		for j, paragraph := range row.Paragraphs {
			content, err := paragraph.Serve(links)
			if err != nil {
				return nil, err
			}
			processed_table.Rows[i].Paragraphs[j] = content
		}
	}

	processed_title, err := serve(table.Title, `{{ .Text }}`)
	if err != nil {
		return nil, err
	}

	processed_table.Title = processed_title

	// TODO: add description row
	// processed_description, err := serve(table.Description, `{{ .Text }}`)
	// if err != nil {
	// 	return nil, err
	// }

	return serve(processed_table, "<div class=\"wrapper\"><table class=\"buttons\" id=\"{{ printf \"%s\" .ID }}\"><colgroup><col style=\"width:22%;\"><col style=\"width:78%;\"></colgroup><thead><tr><th colspan=\"2\">{{ printf \"%s\" .Title }}</th></tr><tr><th>Элемент</th><th>Функция</th></tr></thead><tbody>{{ range .Rows }}<tr><th><div class=\"grid\">{{ range .Images }}<div{{ if .IsWide }} class=\"wide-element\"{{ end }}><img src=\"{{ printf \"%s\" .Path }}\" class=\"grid-image\"></div>{{ end }}</div></th><th>{{ range .Paragraphs }}{{ printf \"%s\" . }}{{ end }}</th></tr>{{ end }}</tbody></table></div>")
}
