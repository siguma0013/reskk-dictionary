// Package model はデータモデルを表現するための構造体を提供する
package model

// Entry : 辞書Entryのインターフェース
type Entry interface {
	GetKanji() string
	GetIndexYomi() string
	GetIndexKey() string
}
