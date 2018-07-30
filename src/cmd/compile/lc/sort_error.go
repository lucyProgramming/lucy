package lc

//import (
//	"strconv"
//	"strings"
//)
//
//type SortErrors []error
//
//func (errs SortErrors) Len() int {
//	return len(errs)
//}
//func (errs SortErrors) Less(i, j int) bool {
//	e1, e2 := errs[i].Error(), errs[j].Error()
//	e1s, e2s := strings.Split(e1, ":"), strings.Split(e2, ":")
//	if string(e1s[0]) < string(e2s[0]) {
//		return true
//	}
//	if string(e1s[0]) > string(e2s[0]) {
//		return false
//	}
//	line1, _ := strconv.Atoi(e1s[1])
//	line2, _ := strconv.Atoi(e2s[1])
//	if line1 < line2 {
//		return true
//	}
//	if line1 > line2 {
//		return false
//	}
//	// line1 == line2
//	c1 := errs.parseColumn(e1s[2])
//	c2 := errs.parseColumn(e2s[2])
//	return c1 < c2
//}
//
//func (errs SortErrors) parseColumn(s string) int {
//	var ret int
//	for _, v := range []byte(s) {
//		if v >= '0' && v <= '9' {
//			ret = ret*10 + int(v-'0')
//		} else {
//			break
//		}
//	}
//	return ret
//}
//
//func (errs SortErrors) Swap(i, j int) {
//	errs[i], errs[j] = errs[j], errs[i]
//}
