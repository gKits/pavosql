package vcache

type VCache[I comparable, T any] struct {
	cache map[uint64]map[I]T
	ver   uint64
}

func New[I comparable, T any](ver uint64) VCache[I, T] {
	return VCache[I, T]{
		cache: map[uint64]map[I]T{},
		ver:   ver,
	}
}

func (vc *VCache[I, T]) Get(id I, ver uint64) (T, bool) {
	for ; ver <= vc.ver; ver++ {
		d, ok := vc.cache[ver][id]
		if ok {
			return d, true
		}
	}
	return *new(T), false
}

func (vc *VCache[I, T]) Cache(i I, d T) {
	vc.cache[vc.ver][i] = d
}

func (vc *VCache[I, T]) NewVersion() uint64 {
	vc.ver++
	return vc.ver
}

func (vc *VCache[I, T]) DeleteVersion(ver uint64) {
	delete(vc.cache, ver)
}
