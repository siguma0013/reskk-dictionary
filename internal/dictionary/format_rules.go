package dictionary

import "regexp"

// FormatRule は辞書ファイルのフォーマットルールです
// パースに関してはjson.Decoderなどに責務を置き、スペースの数などを正規表現で定義する
type FormatRule struct {
	Regexp *regexp.Regexp
	Message string
}

var FormatRules = []FormatRule{
	{regexp.MustCompile(`:\S`), "colon must be followed by one space"},
	{regexp.MustCompile(`:\s{2,}`), "too many spaces after colon"},
	{regexp.MustCompile(`,\S`), "comma must be followed by one space"},
	{regexp.MustCompile(`,\s{2,}`), "too many spaces after comma"},
	{regexp.MustCompile(`\{\s+`), "no space after open brace ({)"},
	{regexp.MustCompile(`\s+\}`), "no space before close brace (})"},
	{regexp.MustCompile(`"\s+`), "no space after double quotation"},
}
