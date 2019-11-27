package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Jeiwan/opscript/spec"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newBuildSpecCmd() *cobra.Command {
	const specURL = "https://en.bitcoin.it/wiki/Script"
	const outFile = "spec.json"
	const totalTables = 10

	cmd := &cobra.Command{
		Use: "buildspec",
		RunE: func(cmd *cobra.Command, args []string) error {
			spec := make(spec.Script)

			resp, err := http.Get(specURL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return err
				}

				return fmt.Errorf("spec page is not available: %d, %s", resp.StatusCode, b)
			}

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				return err
			}

			var scrapedTables uint8
			doc.Find("table.wikitable").Each(func(i int, s *goquery.Selection) {
				if scrapedTables >= totalTables {
					return
				}

				s.Find("tr").Each(func(j int, s *goquery.Selection) {
					cells := s.Find("td").Map(func(k int, s *goquery.Selection) string {
						return s.Text()
					})

					if len(cells) == 0 {
						return
					}

					ops := cellsToOps(cells)

					for _, op := range ops {
						spec[op.Word] = op
					}
				})

				scrapedTables++
			})

			specJSON, err := json.MarshalIndent(spec, "", "    ")
			if err != nil {
				return err
			}

			return ioutil.WriteFile(outFile, specJSON, 0644)
		},
	}
	return cmd
}

func cellsToOps(cells []string) []spec.Opcode {
	if len(cells) == 0 {
		return nil
	}

	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}

	var input, output string
	word := cells[0]
	opcodeDec := cells[1]
	opcode := cells[2]
	if len(cells) == 6 {
		input = cells[3]
		output = cells[4]
	}
	short := cells[len(cells)-1]

	if !strings.HasPrefix(word, "OP_") {
		logrus.Infof("skipping %s\n", word)
		return nil
	}

	var ops []spec.Opcode
	alts := strings.Split(word, ", ")
	rnge := strings.Split(word, "-")

	if len(alts) == 2 && len(rnge) == 2 {
		rnge = strings.Split(alts[1], "-")
		alts = []string{alts[0]}
		opcodeDec = strings.Split(opcodeDec, ", ")[1]
	}

	cleanWord := regexp.MustCompile(`OP_[\w_-]+`)

	if len(alts) == 1 {

		ops = append(ops, spec.Opcode{
			Word:   cleanWord.FindString(alts[0]),
			Opcode: opcode,
			Input:  input,
			Output: output,
			Short:  short,
		})
	}

	if len(alts) == 2 {
		ops = append(ops, []spec.Opcode{
			{
				Word:    cleanWord.FindString(alts[0]),
				WordAlt: cleanWord.FindString(alts[1]),
				Opcode:  opcode,
				Input:   input,
				Output:  output,
				Short:   short,
			},
			{
				Word:    cleanWord.FindString(alts[1]),
				WordAlt: cleanWord.FindString(alts[0]),
				Opcode:  opcode,
				Input:   input,
				Output:  output,
				Short:   short,
			},
		}...)
	}

	if len(rnge) == 2 {
		leftOp := rnge[0]
		rightOp := rnge[1]

		opNumber := regexp.MustCompile(`(OP_[A-Z]*)(\d+)`)

		leftMatches := opNumber.FindStringSubmatch(leftOp)
		rightMatches := opNumber.FindStringSubmatch(rightOp)
		if len(leftMatches) != 3 || len(rightMatches) != 3 {
			logrus.Errorf("skipping %s: invalid opcodes range %q, %q\n", word, leftMatches, rightMatches)
			return ops
		}

		leftN, err := strconv.Atoi(leftMatches[2])
		if err != nil {
			logrus.Errorf("skipping %s: %+v\n", word, err)
			return ops
		}

		rightN, err := strconv.Atoi(rightMatches[2])
		if err != nil {
			logrus.Errorf("skipping %s: %+v\n", word, err)
			return ops
		}

		opcodes := strings.Split(opcodeDec, "-")
		if len(opcodes) != 2 {
			logrus.Errorf("skipping %s: invalid opcode %+v\n", word, opcodeDec)
			return ops
		}

		opcodeLeft, err := strconv.Atoi(opcodes[0])
		if err != nil {
			logrus.Errorf("skipping %s: %+v\n", word, err)
			return ops
		}

		var outputLeft int

		isEmptyOutput := len(output) == 0

		if !isEmptyOutput {
			outputs := strings.Split(output, "-")
			if len(outputs) != 2 {
				logrus.Errorf("skipping %s: invalid output %+v\n", word, output)
				return ops
			}

			outputLeft, err = strconv.Atoi(outputs[0])
			if err != nil {
				logrus.Errorf("skipping %s: %+v\n", word, err)
				return ops
			}
		}

		for i := leftN; i <= rightN; i++ {
			op := spec.Opcode{
				Word:   fmt.Sprintf("%s%d", leftMatches[1], i),
				Opcode: fmt.Sprintf("0x%x", opcodeLeft+i-leftN),
				Input:  input,
				Short:  short,
			}

			if !isEmptyOutput {
				op.Output = strconv.Itoa(outputLeft + i - leftN)
			}

			ops = append(ops, op)
		}
	}

	return ops
}
