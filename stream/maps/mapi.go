package maps

type MapI[K comparable, V any] interface {
	Setter[K, V]
	Getter[K, V]
	Loader[K, V]
	Clear()
	Len() int
	Range(func(k K, v V) bool)
	Has(k K) bool
	Keys() []K
	Values() []V
	Merge(MapI[K, V])
	Equal(MapI[K, V]) bool
	Delete(k K)
}

type Setter[K comparable, V any] interface {
	Set(K, V)
}

type Getter[K comparable, V any] interface {
	Get(k K) (v V)
}

type Loader[K comparable, V any] interface {
	Load(k K) (v V)
}
