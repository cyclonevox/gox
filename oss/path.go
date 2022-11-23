package oss

import (
	`fmt`
	`strings`
)

type Path []string

func NewPath(paths ...any) Path {
	p := make([]string, len(paths))
	for i := range paths {
		p[i] = fmt.Sprintf("%v", paths[i])
	}

	return p
}

func (p Path) Key() string {
	return strings.Join(p, "/")
}
