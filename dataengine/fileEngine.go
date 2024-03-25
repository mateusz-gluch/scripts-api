package dataengine

type FileEngine interface {
	Read(name string) ([]map[string]any, error)
	List(name string) ([]string, error)
}
