package validator

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"siguma0013/reskk-dictionary/internal/model"
)

// ファイル名と索引文字の関係を表すハッシュマップ
var indexKeysMatrix = map[string][]string{
	// あ
	"01-あ.yaml": {"あ"},
	"02-い.yaml": {"い"},
	"03-う.yaml": {"う"},
	"04-え.yaml": {"え"},
	"05-お.yaml": {"お"},

	// か
	"11-か.yaml": {"か", "が"},
	"12-き.yaml": {"き", "ぎ"},
	"13-く.yaml": {"く", "ぐ"},
	"14-け.yaml": {"け", "げ"},
	"15-こ.yaml": {"こ", "ご"},

	// さ
	"21-さ.yaml": {"さ", "ざ"},
	"22-し.yaml": {"し", "じ"},
	"23-す.yaml": {"す", "ず"},
	"24-せ.yaml": {"せ", "ぜ"},
	"25-そ.yaml": {"そ", "ぞ"},

	// た
	"31-た.yaml": {"た", "だ"},
	"32-ち.yaml": {"ち", "ぢ"},
	"33-つ.yaml": {"つ", "づ"},
	"34-て.yaml": {"て", "で"},
	"35-と.yaml": {"と", "ど"},

	// な
	"41-な.yaml": {"な"},
	"42-に.yaml": {"に"},
	"43-ぬ.yaml": {"ぬ"},
	"44-ね.yaml": {"ね"},
	"45-の.yaml": {"の"},

	// は
	"51-は.yaml": {"は", "ば", "ぱ"},
	"52-ひ.yaml": {"ひ", "び", "ぴ"},
	"53-ふ.yaml": {"ふ", "ぶ", "ぷ"},
	"54-へ.yaml": {"へ", "べ", "ぺ"},
	"55-ほ.yaml": {"ほ", "ぼ", "ぽ"},

	// ま
	"61-ま.yaml": {"ま"},
	"62-み.yaml": {"み"},
	"63-む.yaml": {"む"},
	"64-め.yaml": {"め"},
	"65-も.yaml": {"も"},

	// や
	"71-や.yaml": {"や"},
	"73-ゆ.yaml": {"ゆ"},
	"75-よ.yaml": {"よ"},

	// ら
	"81-ら.yaml": {"ら"},
	"82-り.yaml": {"り"},
	"83-る.yaml": {"る"},
	"84-れ.yaml": {"れ"},
	"85-ろ.yaml": {"ろ"},

	// わ
	"91-わ.yaml": {"わ"},
}

// getIndexKeys : ファイルパスから索引文字を取得
func getIndexKeys(filePath string) ([]string, error) {
	// ファイル名を取得
	fileName := filepath.Base(filePath)

	// ファイル名から代表読みを取得
	indexKeys, exists := indexKeysMatrix[fileName]

	// ファイル名チェック
	if !exists {
		return nil, fmt.Errorf("管理対象外のファイル名です")
	}

	return indexKeys, nil
}

// ディレクトリ名と漢字の文字列長の関係を表すハッシュマップ
var kanjiLengthMatrix = map[string]int{
	"常用漢字": 1,
	"二字熟語": 2,
	"三字熟語": 3,
	"四字熟語": 4,
}

// getKanjiLength : ファイルパスから漢字の文字列長を取得
func getKanjiLength(filePath string) (int, error) {
	// ディレクトリ名を取得
	dirName := filepath.Base(filepath.Dir(filePath))

	// ディレクトリ名から漢字の文字列長を取得
	length, exists := kanjiLengthMatrix[dirName]

	// ディレクトリ名チェック
	if !exists {
		return 0, fmt.Errorf("管理対象外のディレクトリ名です")
	}

	return length, nil
}

// validateKanji : 漢字のバリデーションロジック
func validateKanji(kanji string, length int, seenKanji map[string]struct{}) error {
	// 文字列長チェック
	if utf8.RuneCountInString(kanji) != length {
		return fmt.Errorf("漢字の文字列長が不正です")
	}

	// 重複チェック
	if _, exists := seenKanji[kanji]; exists {
		return fmt.Errorf("漢字が重複しています")
	}

	seenKanji[kanji] = struct{}{}

	return nil
}

// 単漢字の音読みに許可された文字列
const allowOnYomiString = "" +
	"あいうえお" +
	"かきくけこ" +
	"さしすせそ" +
	"なにぬねの" +
	"たちつてと" +
	"はひふへほ" +
	"まみむめも" +
	"やゆよ" +
	"らりるれろ" +
	"わ" +
	"ん" +
	"がぎぐげご" +
	"ざじずぜぞ" +
	"だぢづでど" +
	"ばびぶべぼ" +
	"ぱぴぷぺぽ" +
	"っ" +
	"ゃゅょ"

// 単漢字の訓読みと熟語の読みに許可された文字列
const allowYomiString = allowOnYomiString +
	"iue" + // 母音
	"kstnmry" + // 子音(清音)
	"gzb" // 子音(濁音)

// isValidString : 文字列が許可された文字列で構成されているか判定する
func isValidString(value string, allowStrings string) bool {
	// 空文字を許可しない
	if value == "" {
		return false
	}

	for _, rune := range value {
		if !strings.ContainsRune(allowStrings, rune) {
			return false
		}
	}

	return true
}

// validateYomi : 読みのバリデーションロジック
func validateYomi(values []string, allowStrings string) error {
	var invalid []string

	for _, value := range values {
		if !isValidString(value, allowStrings) {
			invalid = append(invalid, value)
		}
	}

	if len(invalid) > 0 {
		return fmt.Errorf("文字種が不正です [%s]", strings.Join(invalid, ", "))
	}

	return nil
}

// validateIndex : 代表読み・索引のバリデーションロジック
func validateIndex(entry model.Entry, profile []string) error {
	// 必須チェック
	if entry.GetIndexYomi() == "" {
		return fmt.Errorf("読みは1つ以上必要です")
	}

	// インデックス一致チェック
	if !slices.Contains(profile, entry.GetIndexKey()) {
		return fmt.Errorf("索引が正しくありません")
	}

	return nil
}
