package replaceholder

import (
	"strings"
)

func strReplace(s, old, new string) (string, error) {
	//newStrBytes, err := json.Marshal(new)
	//if err != nil {
	//	return "", err
	//}
	//// 去除左右两个"
	//newStr := string(newStrBytes[1 : len(newStrBytes)-1])

	rs := strings.Replace(s, old, new, -1)
	return rs, nil
}
