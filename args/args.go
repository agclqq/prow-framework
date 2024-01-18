package args

import "strings"

// TidyParmaNoPrefix 整理参数列表，只返回无前导符(-，--)的参数
func TidyParmaNoPrefix(param []string) []string {
	newParam := make([]string, 0)
	nextSkip := false
	for _, v := range param {
		if (strings.HasPrefix(v, "-") || strings.HasPrefix(v, "--")) && !strings.Contains(v, "=") {
			nextSkip = true
			continue
		}
		if nextSkip {
			continue
		}
		newParam = append(newParam, v)
	}
	return newParam
}

// TidyParmaWithPrefix 整理参数列表，只返回带前导符(-，--)的参数与值
// 写这个方法的原因是flag包不支持命令行中带前导符参数和不带前导符参数混用
func TidyParmaWithPrefix(param []string) map[string]string {
	newParam := make(map[string]string)
	key := ""
	nextSkip := true

	for _, v := range param {
		if len(v) >= 2 && strings.HasPrefix(v, "-") || len(v) >= 3 && strings.HasPrefix(v, "--") {
			v = strings.TrimLeft(v, "-")
			if len(v) == 0 {
				continue
			}

			if strings.Contains(v, "=") {
				vs := strings.SplitN(v, "=", 2)
				newParam[vs[0]] = vs[1]
				nextSkip = false
				continue
			}
			key = v
			newParam[key] = "" //为了防止只有前导符KEY，没有value的情况
			nextSkip = false
			continue
		}
		if !nextSkip {
			newParam[key] = v
			nextSkip = true
		}
	}
	return newParam
}
