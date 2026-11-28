package model

// Dictionary : JSON出力用辞書を表現するデータモデル
// 読みをキー、変換候補の漢字を値とする
type Dictionary map[string][]string

// AddEntry : 辞書Entryの全ての読みに対して漢字を出現順で追加する
func (d Dictionary) AddEntry(entry Entry) {
	for _, yomi := range entry.GetYomi() {
		d[yomi] = append(d[yomi], entry.GetKanji())
	}
}
