package compression_tool

import "sort"

type FrequencyTree struct {
	letter    string
	frequency int
	left      *FrequencyTree
	right     *FrequencyTree
}

type FrequencyTrees []FrequencyTree

func (ft FrequencyTrees) Len() int {
	return len(ft)
}

func (ft FrequencyTrees) Less(i, j int) bool {
	return ft[i].frequency < ft[j].frequency
}

func (ft FrequencyTrees) Swap(i, j int) {
	ft[i], ft[j] = ft[j], ft[i]
}

func (ft FrequencyTrees) BuildHuffmanTree() FrequencyTrees {
	for len(ft) >= 2 {
		ft = ft.iterate()
	}
	return ft
}

func (ft FrequencyTrees) iterate() FrequencyTrees {
	if !(len(ft) >= 2) {
		return ft
	}
	sort.Sort(ft)
	t0, t1 := ft[0], ft[1]
	ft = ft[2:]
	t := merge(t0, t1)
	ft = append(ft, t)
	sort.Sort(ft)

	return ft
}

func merge(ft1, ft2 FrequencyTree) FrequencyTree {
	frequency := ft1.frequency + ft2.frequency
	left := FrequencyTree{}
	right := FrequencyTree{}
	if ft1.frequency < ft2.frequency {
		left, right = ft1, ft2
	} else {
		left, right = ft2, ft1
	}
	return FrequencyTree{frequency: frequency, right: &right, left: &left}
}
