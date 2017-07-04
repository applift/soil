package cluster

import (
	"github.com/hashicorp/serf/serf"
	"sync"
)


// Synchronised storage
type Storage struct {
	mu *sync.Mutex
	state *StorageState
	kv map[string]StorageValue
	cache map[string]string
}

// Put value in storage
func (s *Storage) Put(key string, value StorageValue) {

}

// Merge remote state
func (s *Storage) Merge(state *StorageState, data map[string]StorageValue) {

}



type StorageValue struct {
	LTime serf.LamportTime
	Baseline int64
	Data []byte
}

// Storage state is advertised by nodes on join.
type StorageState struct {

	// Aggregated LTime from storage KV
	LTotal uint64

	// Minor baseline
	Major int64

	// Major baseline
	Minor int64
}

// Set new major
func (s *StorageState) SetMajor(v int64) {
	if s.Major < v {
		s.Major = v
	}
}

// Set new minor and increase major
func (s *StorageState) SetMinor(v, major int64) {
	s.Minor = v
	s.SetMajor(major)
}

// NeedMerge compares itself with candidate
func (s *StorageState) NeedMerge(candidate *StorageState) (res bool) {

	// fast checks
	res = s.Major == 0 || s.Major < candidate.Major || s.Minor > candidate.Minor
	if res {
		return
	}
	res = s.Major == candidate.Major && s.LTotal < candidate.LTotal
	return
}

