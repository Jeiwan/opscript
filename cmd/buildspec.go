package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Jeiwan/opscript/spec"
	"github.com/spf13/cobra"
)

func newBuildSpecCmd() *cobra.Command {
	const outFile = "spec.json"

	cmd := &cobra.Command{
		Use: "buildspec",
		RunE: func(cmd *cobra.Command, args []string) error {
			spec, err := spec.NewFromBitcoinWiki()
			if err != nil {
				return err
			}

			specJSON, err := json.MarshalIndent(spec, "", "    ")
			if err != nil {
				return err
			}

			return ioutil.WriteFile(outFile, specJSON, 0644)
		},
	}
	return cmd
}
