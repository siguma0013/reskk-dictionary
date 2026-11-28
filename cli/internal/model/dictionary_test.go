package model

import "testing"

func TestDictionaryAddEntry(t *testing.T) {
	d := Dictionary{}
	d.AddEntry(CharacterEntry{Kanji: "愛", OnYomi: []string{"あい"}, KunYomi: []string{"いと"}})
	d.AddEntry(WordEntry{Kanji: "藍", Yomi: []string{"あい"}})

	if got := d["あい"]; len(got) != 2 || got[0] != "愛" || got[1] != "藍" {
		t.Errorf("AddEntry() あい = %v", got)
	}
	if got := d["いと"]; len(got) != 1 || got[0] != "愛" {
		t.Errorf("AddEntry() いと = %v", got)
	}
}
