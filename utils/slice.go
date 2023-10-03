package utils

import (
	"strings"
)

func SSliceFormat(ss interface{}) []string {
	switch v := ss.(type) {
	case string:
		return SSliceFilter(strings.Split(v, ","), nil)
	case []string:
		if len(v) == 1 {
			v = strings.Split(v[0], ",")
		}
		return SSliceFilter(v, nil)
	default:
		panic("unsupported param type")
	}
}

func SSliceFilter(ss []string, fn func(string) bool) []string {
	if fn == nil {
		fn = func(s string) bool { return strings.TrimSpace(s) != "" }
	}
	var ns = make([]string, 0)
	for i := range ss {
		if fn(ss[i]) {
			ns = append(ns, ss[i])
		}
	}
	return ns
}

func SSliceMerge(sss ...[]string) []string {
	if len(sss) == 0 {
		return []string{}
	}
	var ss = sss[0]
	for _, vs := range sss[1:] {
		ss = append(ss, vs...)
	}
	return ss
}

func SSliceContain(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func SSliceDiff(ss1, ss2 []string) []string {
	idx := make(map[string]struct{})
	for _, s1 := range ss1 {
		idx[s1] = struct{}{}
	}
	for _, s2 := range ss2 {
		if _, ok := idx[s2]; !ok {
			idx[s2] = struct{}{}
		} else {
			delete(idx, s2)
		}
	}
	res := make([]string, 0)
	for s := range idx {
		res = append(res, s)
	}
	return res
}

func SSliceSame(ss1, ss2 []string) []string {
	idx := make(map[string]struct{})
	for _, s1 := range ss1 {
		idx[s1] = struct{}{}
	}
	res := make([]string, 0)
	for _, s2 := range ss2 {
		if _, ok := idx[s2]; ok {
			res = append(res, s2)
		}
	}
	return res
}
