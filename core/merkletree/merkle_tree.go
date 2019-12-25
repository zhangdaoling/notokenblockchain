package merkletree

import (
	"math"

	"github.com/zhangdaoling/simplechain/common"
	"github.com/zhangdaoling/simplechain/pb"
)

func Build(m *pb.MerkleTree, data [][]byte) {
	if len(data) == 0 {
		m.HashList = make([][]byte, 1)
		m.HashList[0] = nil
		return
	}
	n := int32(math.Exp2(math.Ceil(math.Log2(float64(len(data))))))
	m.LeafNum = n
	m.HashList = make([][]byte, 2*n)
	copy(m.HashList[n-1:n+int32(len(data))-1], data)
	start := n - 1
	end := n + int32(len(data)) - 2
	for {
		for i := start; i <= end; i = i + 2 {
			p := (i - 1) / 2
			if m.HashList[i+1] == nil {
				m.HashList[p] = common.NewSha3(append(m.HashList[i], m.HashList[i]...))
			} else {
				m.HashList[p] = common.NewSha3(append(m.HashList[i], m.HashList[i+1]...))
			}
		}
		start = (start - 1) / 2
		end = (end - 1) / 2
		if start == end {
			break
		}
	}
}

func RootHash(m *pb.MerkleTree) []byte {
	return m.HashList[0]
}
