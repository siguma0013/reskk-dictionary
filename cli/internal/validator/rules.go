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
	"01-a.yaml": {"あ"},
	"02-i.yaml": {"い"},
	"03-u.yaml": {"う"},
	"04-e.yaml": {"え"},
	"05-o.yaml": {"お"},

	// か
	"11-ka.yaml": {"か", "が"},
	"12-ki.yaml": {"き", "ぎ"},
	"13-ku.yaml": {"く", "ぐ"},
	"14-ke.yaml": {"け", "げ"},
	"15-ko.yaml": {"こ", "ご"},

	// さ
	"21-sa.yaml": {"さ", "ざ"},
	"22-si.yaml": {"し", "じ"},
	"23-su.yaml": {"す", "ず"},
	"24-se.yaml": {"せ", "ぜ"},
	"25-so.yaml": {"そ", "ぞ"},

	// た
	"31-ta.yaml": {"た", "だ"},
	"32-ti.yaml": {"ち", "ぢ"},
	"33-tu.yaml": {"つ", "づ"},
	"34-te.yaml": {"て", "で"},
	"35-to.yaml": {"と", "ど"},

	// な
	"41-na.yaml": {"な"},
	"42-ni.yaml": {"に"},
	"43-nu.yaml": {"ぬ"},
	"44-ne.yaml": {"ね"},
	"45-no.yaml": {"の"},

	// は
	"51-ha.yaml": {"は", "ば", "ぱ"},
	"52-hi.yaml": {"ひ", "び", "ぴ"},
	"53-hu.yaml": {"ふ", "ぶ", "ぷ"},
	"54-he.yaml": {"へ", "べ", "ぺ"},
	"55-ho.yaml": {"ほ", "ぼ", "ぽ"},

	// ま
	"61-ma.yaml": {"ま"},
	"62-mi.yaml": {"み"},
	"63-mu.yaml": {"む"},
	"64-me.yaml": {"め"},
	"65-mo.yaml": {"も"},

	// や
	"71-ya.yaml": {"や"},
	"73-yu.yaml": {"ゆ"},
	"75-yo.yaml": {"よ"},

	// ら
	"81-ra.yaml": {"ら"},
	"82-ri.yaml": {"り"},
	"83-ru.yaml": {"る"},
	"84-re.yaml": {"れ"},
	"85-ro.yaml": {"ろ"},

	// わ
	"91-wa.yaml": {"わ"},
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
	"kstn" + // 子音(清音)
	"gzb" // 子音(濁音)

// isValidString : 文字列が許可された文字列で構成されているか判定する
func isValidString(value string, allowStrings string) bool {
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
