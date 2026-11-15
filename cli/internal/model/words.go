package model

// Words : 熟語用辞書の全体を表現するデータモデル
type Words struct {
	Entries []WordEntry `yaml:"entries"`
}

// WordEntry : 熟語用辞書の1単語を表現するデータモデル
type WordEntry struct {
	kanji string   `yaml:"漢字"`
	yomi  []string `yaml:"読み"`
}

// GetKanji : 熟語用辞書の漢字getter実装
func (e WordEntry) GetKanji() string {
	return e.kanji
}

// GetYomi : 熟語用辞書の読みgetter実装
func (e WordEntry) GetYomi() []string {
	return e.yomi
}

// GetIndexReading : 熟語用辞書の代表読みgetter実装
func (e WordEntry) GetIndexYomi() string {
	if len(e.yomi) == 0 {
		return ""
	}

	return e.yomi[0]
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
