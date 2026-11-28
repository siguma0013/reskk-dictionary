package model

import (
	"encoding/json"
	"io"
	"maps"
	"slices"
)

// Dictionary : JSON出力用辞書を表現するデータモデル
// 読みをキー、変換候補の漢字を値とする
type Dictionary map[string][]string

// NewDictionary : 空のDictionaryを生成する
func NewDictionary() Dictionary {
	return Dictionary{}
}

// Merge : 別のDictionaryの内容を、候補の順序を維持して追加する
func (d Dictionary) Merge(other Dictionary) {
	for _, yomi := range other.Yomis() {
		d[yomi] = append(d[yomi], other[yomi]...)
	}
}

// Yomis : 読みを昇順で返却する
func (d Dictionary) Yomis() []string {
	return slices.Sorted(maps.Keys(d))
}

// Encode : doc/merge-command-output-format.md の形式でJSONを書き込む
// キーは昇順、インデントなし、末尾に改行を付ける
func (d Dictionary) Encode(w io.Writer) error {
	if d == nil {
		d = Dictionary{}
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	return encoder.Encode(map[string][]string(d))
}
