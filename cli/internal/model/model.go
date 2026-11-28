// Package model はデータモデルを表現するための構造体を提供する
package model

// KanjiGetter : 漢字を取得できるインターフェース
type KanjiGetter interface {
	GetKanji() string
}

// YomisGetter : 全ての読みを取得できるインターフェース
type YomisGetter interface {
	GetYomis() []string
}

// IndexYomiGetter : 代表読みを取得できるインターフェース
type IndexYomiGetter interface {
	GetIndexYomi() string
}

// IndexKeyGetter : 索引文字を取得できるインターフェース
type IndexKeyGetter interface {
	GetIndexKey() string
}

// Entry : 辞書Entryのインターフェース
// 利用側は必要なメソッドだけを持つ小さなインターフェースを受け取る
type Entry interface {
	KanjiGetter
	YomisGetter
	IndexYomiGetter
	IndexKeyGetter
}
