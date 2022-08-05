package stop

type triEntry[T, U, V any] struct{
	T T
	U U
	V V
}

type trimap[T, U, V comparable] map[triEntry[T, U, V]]triEntry[T, U, V]

func (m trimap[T, U, V]) Set(t T, u U, v V) {
	entry := triEntry[T, U, V]{t, u, v}

	m[triEntry[T, U, V]{T: t}] = entry
	m[triEntry[T, U, V]{U: u}] = entry
	m[triEntry[T, U, V]{V: v}] = entry
}

func (m trimap[T, U, V]) Del(t T, u U, v V) {
	delete(m, triEntry[T, U, V]{T: t})
	delete(m, triEntry[T, U, V]{U: u})
	delete(m, triEntry[T, U, V]{V: v})
}

func (m trimap[T, U, V]) DelT(t T) {
	u, v, _ := m.GetT(t)
	m.Del(t, u, v)
}

func (m trimap[T, U, V]) DelU(u U) {
	t, v, _ := m.GetU(u)
	m.Del(t, u, v)
}

func (m trimap[T, U, V]) DelV(v V) {
	t, u, _ := m.GetV(v)
	m.Del(t, u, v)
}

func (m trimap[T, U, V]) GetT(t T) (U, V, bool) {
	entry, ok := m[triEntry[T, U, V]{T: t}]
	return entry.U, entry.V, ok
}

func (m trimap[T, U, V]) GetU(u U) (T, V, bool) {
	entry, ok := m[triEntry[T, U, V]{U: u}]
	return entry.T, entry.V, ok
}

func (m trimap[T, U, V]) GetV(v V) (T, U, bool) {
	entry, ok := m[triEntry[T, U, V]{V: v}]
	return entry.T, entry.U, ok
}

