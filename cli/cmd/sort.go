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

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort <file>",
	Short: "辞書ファイルを代表読み順にソートして上書き保存する",
	Args:  cobra.ExactArgs(1),
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
			if err := sortCharacterDictionary(absPath); err != nil {
				return err
			}

			cmd.Printf("ソート完了: %s\n", path)

			return nil
		}

		if dicType == loader.WordDictionary {
			if err := sortWordDictionary(absPath); err != nil {
				return err
			}

			cmd.Printf("ソート完了: %s\n", path)

			return nil
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sortCmd)
}

func sortCharacterDictionary(path string) error {
	data, err := loader.LoadYaml[model.Characters](path)
	if err != nil {
		return err
	}

	if errs := validator.ValidateCharacters(data, path); len(errs) > 0 {
		errs = append([]error{fmt.Errorf("validation error: %q", path)}, errs...)
		return errors.Join(errs...)
	}

	if filepath.Base(filepath.Dir(path)) == "常用漢字" {
		return fmt.Errorf("常用漢字はソート対象外です")
	}

	data.Entries = sorter.SortEntries(data.Entries)

	return loader.SaveYaml(path, data)
}

func sortWordDictionary(path string) error {
	data, err := loader.LoadYaml[model.Words](path)
	if err != nil {
		return err
	}

	if errs := validator.ValidateWords(data, path); len(errs) > 0 {
		errs = append([]error{fmt.Errorf("validation error: %q", path)}, errs...)
		return errors.Join(errs...)
	}

	data.Entries = sorter.SortEntries(data.Entries)

	return loader.SaveYaml(path, data)
}
