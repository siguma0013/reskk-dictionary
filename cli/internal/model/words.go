package model

// Words : 熟語用辞書の全体を表現するデータモデル
type Words struct {
	Entries []WordEntry `yaml:"entries"`
}

// WordEntry : 熟語用辞書の1単語を表現するデータモデル
type WordEntry struct {
	Kanji string   `yaml:"漢字"`
	Yomi  []string `yaml:"読み"`
}

// GetKanji : 熟語用辞書の漢字getter実装
func (e WordEntry) GetKanji() string {
	return e.Kanji
}

// GetIndexYomi : 熟語用辞書の代表読みgetter実装
func (e WordEntry) GetIndexYomi() string {
	if len(e.Yomi) == 0 {
		return ""
	}

	return e.Yomi[0]
}

// GetIndexKey : 熟語用辞書の索引文字getter実装
func (e WordEntry) GetIndexKey() string {
	// 代表読みを取得
	reading := e.GetIndexYomi()

	// スライス操作によるpanicガード
	if reading == "" {
		return ""
	}

	// 代表読みの頭文字を索引文字として返却
	return string([]rune(reading)[0])
}
