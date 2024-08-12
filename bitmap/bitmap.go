package bitmap

import (
	"crypto/sha256"
	"encoding/binary"
	"hash/crc32"
	"hash/fnv"
	"math"
)

type BitMap struct {
	data []uint64
	size int
}

// NewBitMap 创建一个新的 BitMap
func NewBitMap(size int) *BitMap {
	return &BitMap{
		data: make([]uint64, (size+63)/64),
		size: size,
	}
}

// Set 将指定位置的位设置为1
func (bm *BitMap) Set(pos int) {
	idx, bit := pos/64, uint(pos%64)
	bm.data[idx] |= 1 << bit
}

// IsSet 检查指定位置的位是否为1
func (bm *BitMap) IsSet(pos int) bool {
	idx, bit := pos/64, uint(pos%64)
	return bm.data[idx]&(1<<bit) != 0
}

// hashFunctions 返回用于计算位图位置的多个哈希函数
func hashFunctions(data []byte, size int) []int {
	hashes := []int{}

	// FNV hash
	fnvHasher := fnv.New64()
	fnvHasher.Write(data)
	fnvHash := int(fnvHasher.Sum64() % uint64(size))
	hashes = append(hashes, fnvHash)

	// CRC32 hash
	crc32Hasher := crc32.NewIEEE()
	crc32Hasher.Write(data)
	crc32Hash := int(crc32Hasher.Sum32() % uint32(size))
	hashes = append(hashes, crc32Hash)

	// SHA-256 hash
	sha256Hash := sha256.Sum256(data)
	sha256HashInt := int(binary.BigEndian.Uint64(sha256Hash[:8]) % uint64(size))
	hashes = append(hashes, sha256HashInt)

	return hashes
}

// Add 添加值到位图
func (bm *BitMap) Add(v []byte) (collision bool) {
	hashes := hashFunctions(v, bm.size)
	i := 0
	for _, hash := range hashes {
		if bm.IsSet(hash) {
			i++
		}
		if i == len(hashes) {
			collision = true
			return
		}
		bm.Set(hash)
	}
	return
}

// Contains 检查值是否在位图中
func (bm *BitMap) Contains(v []byte) bool {
	for _, hash := range hashFunctions(v, bm.size) {
		if !bm.IsSet(hash) {
			return false
		}
	}
	return true
}

// CollisionRate 计算冲突率
/*
	位图不冲突的概率是
	P(no collision) ≈ (1−1/m)**kn

	其中：
	  m 是位图的长度。
	  k 是哈希函数的数量。
	  n 是插入的元素数量

	则，冲突的概率是
	P(collision) ≈ 1 − P(no collision)

	所以在使用位图时，一定要计算好冲突的概率
*/
func CollisionRate(m, k, n int64) float64 {
	if m <= 0 {
		return 0
	}
	if n <= 0 {
		return 0
	}
	if k <= 0 {
		k = 1
	}

	// P(no collision) ≈ (1−1/m)**kn
	// P(collision) ≈ 1 − P(no collision)
	return 1 - math.Pow(1-1/float64(m), float64(k*n))
}
