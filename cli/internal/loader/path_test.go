package loader

import "testing"

func TestDictionaryTypeFromPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    DictionaryType
		wantErr bool
	}{
		{
			name: "単漢字ディレクトリ",
			path: "/dictionary/単漢字/常用漢字/a.yaml",
			want: CharacterDictionary,
		},
		{
			name: "熟語ディレクトリ",
			path: "/dictionary/熟語/a.yaml",
			want: WordDictionary,
		},
		{
			name:    "対象外ディレクトリ",
			path:    "/dictionary/常用漢字/a.yaml",
			wantErr: true,
		},
		{
			name:    "名前の一部が一致するだけのディレクトリ",
			path:    "/dictionary/単漢字用/a.yaml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDictionaryType(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DictionaryTypeFromPath() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("DictionaryTypeFromPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
