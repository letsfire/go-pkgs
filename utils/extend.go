package utils

type Extend map[string]interface{}

func (ex Extend) Int(key string) int {
	return int(ex.Float64(key))
}

func (ex Extend) Ints(key string) []int {
	iv := make([]int, 0)
	ex.sliceRange(key, func(i interface{}) {
		iv = append(iv, int(i.(float64)))
	})
	return iv
}

func (ex Extend) Int32(key string) int32 {
	return int32(ex.Float64(key))
}

func (ex Extend) Int32s(key string) []int32 {
	iv := make([]int32, 0)
	ex.sliceRange(key, func(i interface{}) {
		iv = append(iv, int32(i.(float64)))
	})
	return iv
}

func (ex Extend) Int64(key string) int64 {
	return int64(ex.Float64(key))
}

func (ex Extend) Int64s(key string) []int64 {
	iv := make([]int64, 0)
	ex.sliceRange(key, func(i interface{}) {
		iv = append(iv, int64(i.(float64)))
	})
	return iv
}

func (ex Extend) Float64(key string) float64 {
	if v, ok := ex[key]; ok {
		return v.(float64)
	}
	return 0
}

func (ex Extend) Float64s(key string) []float64 {
	fv := make([]float64, 0)
	ex.sliceRange(key, func(i interface{}) {
		fv = append(fv, i.(float64))
	})
	return fv
}

func (ex Extend) Extend(key string) Extend {
	if v, ok := ex[key]; ok {
		if v1, ok1 := v.(Extend); ok1 {
			return v1
		}
		if v2, ok2 := v.(map[string]interface{}); ok2 {
			return v2
		}
	}
	return Extend{}
}

func (ex Extend) Extends(key string) []Extend {
	ev := make([]Extend, 0)
	ex.sliceRange(key, func(i interface{}) {
		ev = append(ev, i.(map[string]interface{}))
	})
	return ev
}

func (ex Extend) String(key string) string {
	if v, ok := ex[key]; ok {
		return v.(string)
	}
	return ""
}

func (ex Extend) Strings(key string) []string {
	sv := make([]string, 0)
	ex.sliceRange(key, func(i interface{}) {
		sv = append(sv, i.(string))
	})
	return sv
}

func (ex Extend) Merge(ex2 Extend) Extend {
	for key, val := range ex2 {
		ex[key] = val
	}
	return ex
}

func (ex Extend) Exclude(keys ...string) Extend {
	ki := make(map[string]struct{})
	for _, key := range keys {
		ki[key] = struct{}{}
	}
	var nex = make(Extend)
	for k, v := range ex {
		if _, ok := ki[k]; ok {
			continue
		}
		nex[k] = v
	}
	return nex
}

func (ex Extend) sliceRange(key string, fn func(interface{})) {
	if v, ok := ex[key]; ok {
		vs := v.([]interface{})
		for i := range vs {
			fn(vs[i])
		}
	}
}
