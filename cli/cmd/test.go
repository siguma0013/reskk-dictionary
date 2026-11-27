/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"siguma0013/reskk-dictionary/internal/loader"
	"siguma0013/reskk-dictionary/internal/model"
	"siguma0013/reskk-dictionary/internal/sorter"
	"siguma0013/reskk-dictionary/internal/validator"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test <file>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// 存在確認
		if err := loader.FileExists(path); err != nil {
			return err
		}

		// 絶対パスに変換
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		// 辞書タイプを取得
		dicType, err := loader.GetDictionaryType(absPath)
		if err != nil {
			return err
		}

		if dicType == loader.CharacterDictionary {
			if err := processCaracterDictionary(absPath); err != nil {
				return err
			}

			cmd.Printf("テスト完了: %s\n", path)

			return nil
		}

		if dicType == loader.WordDictionary {
			if err := processWordDictionary(absPath); err != nil {
				return err
			}

			cmd.Printf("テスト完了: %s\n", path)

			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	rootCmd.SilenceUsage = true

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func processCaracterDictionary(path string) error {
	data, err := loader.LoadYaml[model.Characters](path)

	if err != nil {
		return err
	}

	if errs := validator.ValidateCharacters(data, path); len(errs) > 0 {
		errs = append([]error{fmt.Errorf("validation error: %q", path)}, errs...)
		return errors.Join(errs...)
	}

	return nil
}

func processWordDictionary(path string) error {
	data, err := loader.LoadYaml[model.Words](path)

	if err != nil {
		return err
	}

	if errs := validator.ValidateWords(data, path); len(errs) > 0 {
		errs = append([]error{fmt.Errorf("validation error: %q", path)}, errs...)
		return errors.Join(errs...)
	}

	if errs := sorter.IsSorted(data.Entries); len(errs) > 0 {
		errs = append([]error{fmt.Errorf("sort error: %q", path)}, errs...)
		return errors.Join(errs...)
	}

	return nil
}
