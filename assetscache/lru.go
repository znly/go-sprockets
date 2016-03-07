package assetscache

// assetLru is an LRU cache, safe for concurrent access.
type assetLru struct {
	maxEntries int

	ll    *assetLruList
	cache map[int64]*assetLruListElement
}

// *entry is the type stored in each *list.Element.
type assetLruListEntry struct {
	key   int64
	value *assetCache
}

// New returns a new cache with the provided maximum items.
func newAssetLru(maxEntries int) *assetLru {
	return &assetLru{
		maxEntries: maxEntries,
		ll:         newAssetLruList(),
		cache:      make(map[int64]*assetLruListElement),
	}
}

// Add adds the provided key and value to the cache, evicting
// an old item if necessary.
func (al *assetLru) Add(key int64, value *assetCache) {
	// Already in cache?
	if ee, ok := al.cache[key]; ok {
		al.ll.MoveToFront(ee)
		ee.Value.value = value
		return
	}

	// Add to cache if not present
	ele := al.ll.PushFront(&assetLruListEntry{key, value})
	al.cache[key] = ele

	if al.ll.Len() > al.maxEntries {
		al.removeOldest()
	}
}

// Get fetches the key's value from the cache.
// The ok result will be true if the item was found.
func (al *assetLru) Get(key int64) (value *assetCache, ok bool) {
	if ele, hit := al.cache[key]; hit {
		al.ll.MoveToFront(ele)
		return ele.Value.value, true
	}
	return
}

// RemoveOldest removes the oldest item in the cache and returns its key and value.
// If the cache is empty, the empty string and nil are returned.
func (al *assetLru) RemoveOldest() (key int64, value *assetCache) {
	return al.removeOldest()
}

// note: must hold al.mu
func (al *assetLru) removeOldest() (key int64, value *assetCache) {
	ele := al.ll.Back()
	if ele == nil {
		return
	}
	al.ll.Remove(ele)
	ent := ele.Value
	delete(al.cache, ent.key)
	return ent.key, ent.value

}

// Len returns the number of items in the cache.
func (al *assetLru) Len() int {
	return al.ll.Len()
}
