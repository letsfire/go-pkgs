package utils

import (
	"math/rand"
	"strconv"
)

// Atoi 字符串转数字
func Atoi(s string) int {
	if i, e := strconv.Atoi(s); e != nil {
		panic(e)
	} else {
		return i
	}
}

// MaxInt 获取最大整数
func MaxInt(vs ...int) int {
	var max = vs[0]
	for i := 1; i < len(vs); i++ {
		if max < vs[i] {
			max = vs[i]
		}
	}
	return max
}

// MinInt 获取最小整数
func MinInt(vs ...int) int {
	var min = vs[0]
	for i := 1; i < len(vs); i++ {
		if min > vs[i] {
			min = vs[i]
		}
	}
	return min
}

// SumInt 计算整数总和
func SumInt(args ...int) (v int) {
	for _, arg := range args {
		v += arg
	}
	return
}

// BitSplitInt 按位分解整数
func BitSplitInt(v int) (vs []int) {
	vs = make([]int, 0)
	for i := 1; i <= v; i *= 2 {
		if (i & v) == i {
			vs = append(vs, i)
		}
	}
	return vs
}

// InSliceInt 判断某值是否在切片中
func InSliceInt(v int, is []int) bool {
	for _, i := range is {
		if i == v {
			return true
		}
	}
	return false
}

// MaxInt32 获取最大整数
func MaxInt32(vs ...int32) int32 {
	var max = vs[0]
	for i := 1; i < len(vs); i++ {
		if max < vs[i] {
			max = vs[i]
		}
	}
	return max
}

// MinInt32 获取最小整数
func MinInt32(vs ...int32) int32 {
	var min = vs[0]
	for i := 1; i < len(vs); i++ {
		if min > vs[i] {
			min = vs[i]
		}
	}
	return min
}

// InSliceInt32 判断某值是否在切片中
func InSliceInt32(v int32, is []int32) bool {
	for _, i := range is {
		if i == v {
			return true
		}
	}
	return false
}

// MaxInt64 获取最大整数
func MaxInt64(vs ...int64) int64 {
	var max = vs[0]
	for i := 1; i < len(vs); i++ {
		if max < vs[i] {
			max = vs[i]
		}
	}
	return max
}

// MinInt64 获取最小整数
func MinInt64(vs ...int64) int64 {
	var min = vs[0]
	for i := 1; i < len(vs); i++ {
		if min > vs[i] {
			min = vs[i]
		}
	}
	return min
}

// SumInt64 计算64位整数总和
func SumInt64(args ...int64) (v int64) {
	for _, arg := range args {
		v += arg
	}
	return
}

// BitSplitInt64 按位分解64位整数
func BitSplitInt64(v int64) (vs []int64) {
	vs = make([]int64, 0)
	for i := int64(1); i <= v; i *= 2 {
		if (i & v) == i {
			vs = append(vs, i)
		}
	}
	return vs
}

// InSliceInt64 判断某值是否在切片中
func InSliceInt64(v int64, is []int64) bool {
	for _, i := range is {
		if i == v {
			return true
		}
	}
	return false
}

// RandIntn 屏蔽panic
func RandIntn(n int) int {
	if n <= 0 {
		return n
	}
	return rand.Intn(n)
}
