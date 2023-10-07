package filemgr

type FileMgr[T any] interface {
	ReadFile(path string) (T, error)
	WriteFile(path string, data T) error
}
