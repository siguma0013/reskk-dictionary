package model

// Dictionary : JSON出力用辞書を表現するデータモデル
// 読みをキー、変換候補の漢字を値とする
// 制約: ValidateCharacters/ValidateWordsを通過した単漢字・熟語は読みを1つ以上持つため、AddEntryでnil候補は生成されない
type Dictionary map[string][]string

// AddEntry : 辞書Entryの全ての読みに対して漢字を出現順で追加する
func (d Dictionary) AddEntry(entry Entry) {
	for _, yomi := range entry.GetYomi() {
		d[yomi] = append(d[yomi], entry.GetKanji())
	}
}
