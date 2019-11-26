package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

type opcodeSpec struct {
	Word   string `json:"word"`
	Opcode string `json:"opcode"`
	Input  string `json:"input"`
	Output string `json:"output"`
	Short  string `json:"short"`
}

func newBuildSpecCmd() *cobra.Command {
	const specURL = "https://en.bitcoin.it/wiki/Script"
	const outFile = "spec.json"
	const totalTables = 10

	cmd := &cobra.Command{
		Use: "buildspec",
		RunE: func(cmd *cobra.Command, args []string) error {
			spec := make(map[string]opcodeSpec)

			resp, err := http.Get(specURL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			z := html.NewTokenizer(resp.Body)

			var insideTable bool
			var insideRow bool
			var scrapedTables uint8

		Loop:
			for {
				tt := z.Next()

				switch {
				case tt == html.ErrorToken:
					break Loop

				case tt == html.StartTagToken:
					t := z.Token()

					if !insideTable {
						if isTable(t) {
							insideTable = true
						}
						continue
					}

					if !insideRow {
						if isTableRow(t) {
							insideRow = true
						}
						continue
					}

					if op := scrapeSpec(z); op != nil {
						spec[op.Word] = *op
					}

				case tt == html.EndTagToken:
					t := z.Token()

					if insideTable {
						scrapedTables++

						if scrapedTables >= totalTables {
							break Loop
						}

						if isTable(t) {
							insideTable = false
						}
					}

					if insideRow {
						if isTableRow(t) {
							insideRow = false
						}
					}
				}
			}

			specJSON, err := json.Marshal(spec)
			if err != nil {
				return err
			}

			return ioutil.WriteFile(outFile, specJSON, 0644)
		},
	}
	return cmd
}

func getClass(t html.Token) string {
	for _, a := range t.Attr {
		if a.Key == "class" {
			return a.Val
		}
	}

	return ""
}

func isTable(t html.Token) bool {
	return t.Data == "table" && getClass(t) == "wikitable"
}

func isTableRow(t html.Token) bool {
	return t.Data == "tr"
}

func isHeaderCell(t html.Token) bool {
	return t.Data == "th"
}

func isDataCell(t html.Token) bool {
	return t.Data == "td"
}

func scrapeSpec(z *html.Tokenizer) *opcodeSpec {
	var cells []string
	var spec opcodeSpec
	var insideDataCell bool

Loop:
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			break Loop

		case tt == html.StartTagToken:
			t := z.Token()

			if !insideDataCell {
				if isDataCell(t) {
					insideDataCell = true
				}
			}

		case tt == html.TextToken:
			if !insideDataCell {
				continue
			}

			t := z.Token()

			cells = append(cells, t.Data)

		case tt == html.EndTagToken:
			t := z.Token()

			if isTableRow(t) {
				break Loop
			}
		}
	}

	if len(cells) == 0 {
		return nil
	}

	spec.Word = strings.TrimSpace(cells[0])
	spec.Opcode = strings.TrimSpace(cells[1])
	if len(cells) == 6 {
		spec.Input = strings.TrimSpace(cells[3])
		spec.Output = strings.TrimSpace(cells[4])
	}
	spec.Short = strings.TrimSpace(cells[len(cells)-1])

	fmt.Printf("%s ", spec.Word)

	return &spec
}
