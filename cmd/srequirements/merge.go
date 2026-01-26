package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFile string
)

var mergeCmd = &cobra.Command{
	Use:   "merge [files...]",
	Short: "Merge multiple JSON files into one",
	Long:  `Merge multiple JSON files into one. The files are merged in the order they are provided.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := mergeFiles(args, outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully merged files into %s\n", outputFile)
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.Flags().StringVarP(&outputFile, "output", "o", "merged.json", "Output file name")
}

func mergeFiles(files []string, output string) error {
	var mergedData map[string]interface{}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading file %s: %w", file, err)
		}

		var currentData map[string]interface{}
		if err := json.Unmarshal(data, &currentData); err != nil {
			return fmt.Errorf("unmarshaling json from file %s: %w", file, err)
		}

		if mergedData == nil {
			mergedData = currentData
		} else {
			mergedData = deepMerge(mergedData, currentData)
		}
	}

	mergedJSON, err := json.MarshalIndent(mergedData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling merged data to json: %w", err)
	}

	if err := os.WriteFile(output, mergedJSON, 0600); err != nil {
		return fmt.Errorf("writing merged json to file %s: %w", output, err)
	}

	return nil
}

func deepMerge(a, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		if va, ok := a[k]; ok {
			if vaMap, ok := va.(map[string]interface{}); ok {
				if vMap, ok := v.(map[string]interface{}); ok {
					a[k] = deepMerge(vaMap, vMap)
					continue
				}
			}
		}
		a[k] = v
	}
	return a
}
