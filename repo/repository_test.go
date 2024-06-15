package repo_test

import (
	"fmt"
	"testing"

	"github.com/aol-nnov/debian/repo"
)

func TestNewRepo(t *testing.T) {
	repo, err := repo.New("https://mirror.yandex.ru/debian/", "bookworm")

	if err != nil {
		t.Fatal(err)
	}

	si, _ := repo.Component("main").SourceIndex()
	fmt.Println(si.FindByName("doxygen"))

	// bi, _ := repo.Component("main").BinaryIndex("amd64")
	// fmt.Println(bi.FindByName("mc"))
}
