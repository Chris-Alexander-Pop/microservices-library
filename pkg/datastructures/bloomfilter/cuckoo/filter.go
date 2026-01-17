package cuckoo

import (
	"encoding/binary"
	"hash/fnv"
	"math/rand"
)

const (
	maxKicks = 500
)

// Filter is a simplified Cuckoo Filter.
// Supports Add, Contains, Delete.
type Filter struct {
	buckets    []bucket
	count      uint
	bucketSize uint // 4 is standard
	capacity   uint
}

type bucket [4]fingerprint
type fingerprint uint16 // 16-bit fingerprint

func New(capacity uint) *Filter {
	// Capacity needs to be power of 2 for fast modulo if possible, but standard modulo works.
	// We'll aim for load factor 95%.
	numBuckets := (capacity + 3) / 4 // ceil(capacity / 4)
	return &Filter{
		buckets:    make([]bucket, numBuckets),
		bucketSize: 4,
		capacity:   capacity,
	}
}

func (f *Filter) Add(data []byte) bool {
	if f.Contains(data) {
		return true
	}

	fp := getFingerprint(data)
	i1 := getIndex(data, uint(len(f.buckets)))
	i2 := getAltIndex(i1, fp, uint(len(f.buckets)))

	if f.insert(fp, i1) || f.insert(fp, i2) {
		f.count++
		return true
	}

	// Relocation (Kick)
	currIndex := i1
	if rand.Intn(2) == 1 {
		currIndex = i2
	}

	for k := 0; k < maxKicks; k++ {
		// kick random item from bucket
		randSlot := rand.Intn(4)
		oldFp := f.buckets[currIndex][randSlot]
		f.buckets[currIndex][randSlot] = fp

		fp = oldFp
		currIndex = getAltIndex(currIndex, fp, uint(len(f.buckets)))

		if f.insert(fp, currIndex) {
			f.count++
			return true
		}
	}

	return false // Filter full
}

func (f *Filter) Contains(data []byte) bool {
	fp := getFingerprint(data)
	i1 := getIndex(data, uint(len(f.buckets)))
	i2 := getAltIndex(i1, fp, uint(len(f.buckets)))

	return f.bucketContains(i1, fp) || f.bucketContains(i2, fp)
}

func (f *Filter) Delete(data []byte) bool {
	fp := getFingerprint(data)
	i1 := getIndex(data, uint(len(f.buckets)))
	if f.deleteFromBucket(i1, fp) {
		f.count--
		return true
	}
	i2 := getAltIndex(i1, fp, uint(len(f.buckets)))
	if f.deleteFromBucket(i2, fp) {
		f.count--
		return true
	}
	return false
}

func (f *Filter) insert(fp fingerprint, i uint) bool {
	for j := 0; j < 4; j++ {
		if f.buckets[i][j] == 0 {
			f.buckets[i][j] = fp
			return true
		}
	}
	return false
}

func (f *Filter) bucketContains(i uint, fp fingerprint) bool {
	for j := 0; j < 4; j++ {
		if f.buckets[i][j] == fp {
			return true
		}
	}
	return false
}

func (f *Filter) deleteFromBucket(i uint, fp fingerprint) bool {
	for j := 0; j < 4; j++ {
		if f.buckets[i][j] == fp {
			f.buckets[i][j] = 0
			return true
		}
	}
	return false
}

func getFingerprint(data []byte) fingerprint {
	h := fnv.New32a()
	h.Write(data)
	val := h.Sum32()
	fp := fingerprint(val & 0xFFFF)
	if fp == 0 {
		fp = 1 // 0 is reserved for empty
	}
	return fp
}

func getIndex(data []byte, numBuckets uint) uint {
	h := fnv.New64a()
	h.Write(data)
	return uint(h.Sum64() % uint64(numBuckets))
}

func getAltIndex(i uint, fp fingerprint, numBuckets uint) uint {
	// index2 = index1 ^ hash(fingerprint)
	// This property allows calculating alternate index knowing only current index and fingerprint

	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, uint16(fp))
	h := fnv.New64a()
	h.Write(buf)
	hashFp := uint(h.Sum64())

	return (i ^ hashFp) % numBuckets
}
