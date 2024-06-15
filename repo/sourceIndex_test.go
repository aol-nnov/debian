package repo_test

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/aol-nnov/debian/internal/universalreader"
	"github.com/aol-nnov/debian/repo"
	"golang.org/x/exp/slices"
)

// This benchmark shows near-zero difference between "Adaptive Hashing" and "Full Hashing" beforehand
// So, "Full Hashing" won and has been chosen
func BenchmarkSearchVariants(b *testing.B) {
	packagessToSearch := []string{
		"yabar",
		"0ad",
		"apt",
		"dpkg",
		"mc",
		"autoclass",
		"nano",
		"qbittorrent",
		"zzzeeksphinx",
		"xz",
	}
	// src := "http://ftp.debian.org/debian/dists/bullseye/main/source/Sources.gz"
	src := "./testdata/Sources"

	start := time.Now()
	si, err := repo.NewSourceIndex(universalreader.New(src))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("parsed %d src packages in %v\n", len(si.Packages), time.Since(start))

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	b.Run("LinearFixed", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			for _, pkg := range packagessToSearch {
				slices.IndexFunc(si.Packages, func(item repo.SourceIndexItem) bool {
					return item.Name == pkg
				})
			}
		}
	})

	b.Run("LinearRandom", func(b *testing.B) {
		// 10 packages for each run
		packagessToSearch := make([]string, 10*b.N)

		for i := 0; i < len(packagessToSearch); i++ {
			packagessToSearch[i] = si.Packages[r.Intn(len(si.Packages))].Name
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, pkg := range packagessToSearch[i*10 : (i*10)+10] {
				slices.IndexFunc(si.Packages, func(item repo.SourceIndexItem) bool {
					return item.Name == pkg
				})
			}
		}
	})

	b.Run("AdaptiveFixed", func(b *testing.B) {
		byNameAdaptive := make(map[string]int, len(si.Packages))

		for i := 0; i < b.N; i++ {
			for _, pkg := range packagessToSearch {
				if _, found := byNameAdaptive[pkg]; found {
					continue
				}

				idx := slices.IndexFunc(si.Packages, func(item repo.SourceIndexItem) bool {
					return item.Name == pkg
				})

				if idx != -1 {
					byNameAdaptive[pkg] = idx
				}
			}
		}
	})

	b.Run("AdaptiveRandom", func(b *testing.B) {
		byNameAdaptive := make(map[string]int, len(si.Packages))
		// 10 packages for each run
		packagessToSearch := make([]string, 10*b.N)

		for i := 0; i < len(packagessToSearch); i++ {
			packagessToSearch[i] = si.Packages[r.Intn(len(si.Packages))].Name
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, pkg := range packagessToSearch[i*10 : (i*10)+10] {
				if _, found := byNameAdaptive[pkg]; found {
					continue
				}

				idx := slices.IndexFunc(si.Packages, func(item repo.SourceIndexItem) bool {
					return item.Name == pkg
				})

				if idx != -1 {
					byNameAdaptive[pkg] = idx
				}
			}
		}
	})

	b.Run("HashingFixed", func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			for _, pkg := range packagessToSearch {
				si.FindByName(pkg)
			}
		}
	})

	b.Run("HashingRandom", func(b *testing.B) {

		// 10 packages for each run
		packagessToSearch := make([]string, 10*b.N)

		for i := 0; i < len(packagessToSearch); i++ {
			packagessToSearch[i] = si.Packages[r.Intn(len(si.Packages))].Name
		}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, pkg := range packagessToSearch[i*10 : (i*10)+10] {
				si.FindByName(pkg)
			}
		}
	})
}

func TestTopoSort(t *testing.T) {
	si, err := repo.NewSourceIndex(universalreader.New(""))

	if err != nil {
		t.Fatal(err)
	}

	topo, _, err := si.BuildOrder()

	if err != nil {
		t.Fatal(err)
	}

	for _, p := range topo {
		fmt.Println(p)
	}
}
