package entry

import "regexp"

// Entry 辞書ファイルの構造定義
type Entry struct {
	Key   string   `json:"key"`
	Value []string `json:"value"`
}

type FormatRule struct {
	Regexp *regexp.Regexp
	Message string
}

var FormatRules = []FormatRule{
	{regexp.MustCompile(`":\S`), "colon must be followed by one space"},
	{regexp.MustCompile(`":\s{2,}"`), "too many spaces after colon"},
	{regexp.MustCompile(`,\S`), "comma must be followed by one space"},
	{regexp.MustCompile(`,\s{2,}`), "too many spaces after comma"},
	{regexp.MustCompile(`\{\s+`), "no space after open brace ({)"},
	{regexp.MustCompile(`\s+\}`), "no space before close brace (})"},
}
