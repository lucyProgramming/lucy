package lc

import (
	"strconv"
	"strings"
)

type Errors []error

func (e Errors) Len() int {
	return len(e)
}
func (e Errors) Less(i, j int) bool {
	e1, e2 := e[i].Error(), e[j].Error()
	e1s, e2s := strings.Split(e1, ":"), strings.Split(e2, ":")
	if string(e1s[0]) < string(e2s[0]) {
		return true
	}
	if string(e1s[0]) > string(e2s[0]) {
		return false
	}
	line1, _ := strconv.Atoi(e1s[1])
	line2, _ := strconv.Atoi(e2s[1])
	if line1 < line2 {
		return true
	}
	if line1 > line2 {
		return false
	}
	// line1 == line2
	c1 := e.parseColumn(e1s[2])
	c2 := e.parseColumn(e2s[2])
	return c1 < c2
}

func (e Errors) parseColumn(s string) int {
	var ret int
	for _, v := range []byte(s) {
		if v >= '0' && v <= '9' {
			ret = ret*10 + int((v - '0'))
		}
	}
	return ret
}

func (e Errors) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
