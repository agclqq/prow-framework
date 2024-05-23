package prettytable

import (
	"fmt"
	"strings"
)

// Markdown 格式化为 Markdown 表格
func Markdown(data [][]string) string {
	rs := ""
	for _, row := range data {
		for _, col := range row {
			// 将换行符替换为空格，以便在Markdown表格中换行显示
			col = strings.ReplaceAll(col, "\n", " ")
			rs += fmt.Sprintf("| %s ", col)
		}
		rs += fmt.Sprintln("|")
	}
	return rs
}

// PlainText 格式化为纯文本表格
func PlainText(data [][]string) string {
	if len(data) == 0 {
		return ""
	}
	colList := make([]int, len(data[0]))
	anySlice := make([][]interface{}, 0)
	// 遍历data数组，转换为 interface{} 切片
	for _, v := range data {
		tmpS := make([]interface{}, 0)
		for ii, vv := range colList {
			tmpS = append(tmpS, strings.ReplaceAll(v[ii], "\n", " "))
			if len(v[ii]) > vv {
				colList[ii] = len(v[ii])
			}
		}
		anySlice = append(anySlice, tmpS)
	}
	format := ""
	for _, v := range colList {
		format += fmt.Sprintf("|%%-%ds", v)
	}
	format += "|\n"
	// 转换为 interface{} 切片

	rs := ""
	for _, width := range colList {
		rs += " " + strings.Repeat("_", width)
	}
	rs += "\n"
	for _, v := range anySlice {
		rs += fmt.Sprintf(format, v...)
	}
	// 打印底部边框
	for _, width := range colList {
		rs += " " + strings.Repeat("‾", width)
	}
	rs += "\n"
	return rs
}
