package entry

// Entry 辞書ファイルの構造定義
type Entry struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}
