package compression_tool

import (
	"sort"
	"testing"
)

func TestFrequencyTreesSort(t *testing.T) {
	ft := testingFrequencies()
	sort.Sort(ft)

	if ft[0].letter != "Z" || ft[7].frequency != 120 {
		t.Errorf("sorting is not working properly")
	}
}

func TestBuildingHuffmanTree(t *testing.T) {
	ft := testingFrequencies()
	ft = ft.BuildHuffmanTree()
	ftTest := testingHuffmanTree()
	result := equalTrees(&ft[0], &ftTest)
	if !result {
		t.Errorf("expected trees to be equal")
	}
}

func equalTrees(a, b *FrequencyTree) bool {
	if a == nil && b == nil {
		return true
	}
	return a.frequency == b.frequency &&
		equalTrees(a.right, b.right) &&
		equalTrees(a.left, b.left)
}

// https://opendsa-server.cs.vt.edu/ODSA/Books/CS3/html/Huffman.html
func testingFrequencies() FrequencyTrees {
	ft1 := FrequencyTree{letter: "C", frequency: 32}
	ft2 := FrequencyTree{letter: "D", frequency: 42}
	ft3 := FrequencyTree{letter: "E", frequency: 120}
	ft4 := FrequencyTree{letter: "K", frequency: 7}
	ft5 := FrequencyTree{letter: "L", frequency: 42}
	ft6 := FrequencyTree{letter: "M", frequency: 24}
	ft7 := FrequencyTree{letter: "U", frequency: 37}
	ft8 := FrequencyTree{letter: "Z", frequency: 2}
	return FrequencyTrees{ft1, ft2, ft3, ft4, ft5, ft6, ft7, ft8}
}

func testingHuffmanTree() FrequencyTree {
	C := FrequencyTree{letter: "C", frequency: 32}
	D := FrequencyTree{letter: "D", frequency: 42}
	E := FrequencyTree{letter: "E", frequency: 120}
	K := FrequencyTree{letter: "K", frequency: 7}
	L := FrequencyTree{letter: "L", frequency: 42}
	M := FrequencyTree{letter: "M", frequency: 24}
	U := FrequencyTree{letter: "U", frequency: 37}
	Z := FrequencyTree{letter: "Z", frequency: 2}

	ft9 := FrequencyTree{letter: "", frequency: 9, left: &Z, right: &K}
	ft33 := FrequencyTree{letter: "", frequency: 33, left: &ft9, right: &M}
	ft65 := FrequencyTree{letter: "", frequency: 65, left: &C, right: &ft33}
	ft107 := FrequencyTree{letter: "", frequency: 107, left: &L, right: &ft65}
	ft79 := FrequencyTree{letter: "", frequency: 79, left: &U, right: &D}
	ft186 := FrequencyTree{letter: "", frequency: 186, left: &ft79, right: &ft107}
	return FrequencyTree{letter: "", frequency: 306, left: &E, right: &ft186}
}
