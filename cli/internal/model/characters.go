package model

// Characters : 単漢字用辞書の全体を表現するデータモデル
type Characters struct {
	Entries []CharacterEntry `yaml:"entries"`
}

// CharacterEntry : 単漢字用辞書の1漢字を表現するデータモデル
type CharacterEntry struct {
	Kanji   string   `yaml:"漢字"`
	OnYomi  []string `yaml:"音読み"`
	KunYomi []string `yaml:"訓読み"`
}

// GetKanji : 単漢字用辞書の漢字getter実装
func (e CharacterEntry) GetKanji() string {
	return e.Kanji
}

// GetIndexYomi : 単漢字用辞書の代表読みgetter実装
func (e CharacterEntry) GetIndexYomi() string {
	if len(e.OnYomi) > 0 {
		return e.OnYomi[0]
	}

	if len(e.KunYomi) > 0 {
		return e.KunYomi[0]
	}

	return ""
}

// GetIndexKey : 単漢字用辞書の索引文字getter実装
func (e CharacterEntry) GetIndexKey() string {
	// 代表読みを取得
	reading := e.GetIndexYomi()

	// スライス操作によるpanicガード
	if reading == "" {
		return ""
	}

	// 代表読みの頭文字を索引文字として返却
	return string([]rune(reading)[0])
}
