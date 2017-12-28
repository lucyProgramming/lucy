package jvm

func mkPath(path, name string) string {
	if path == "" {
		return name
	}
	return path + "$" + name
}
