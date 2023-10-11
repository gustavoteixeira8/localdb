package storagemgr

type StorageMgr[T any] interface {
	ReadFile(path string) (T, error)
	WriteFile(path string, data T) error
}
