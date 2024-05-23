package pretty

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar 进度条，在当前行展示进度条
// col 进度条长度
// interval 进度条刷新间隔
func ProgressBar(col int, interval time.Duration) {
	bar := fmt.Sprintf("\033[2K\033[1G[%%-%vs]", col) //\033[2K 清除当前行； \033[1G 将光标移动到行首。\033等同于\u001B
	baseBar(col, interval, bar)
}

// ProgressBarWithNewLine 进度条，在新行展示进度条
// col 进度条长度
// interval 进度条刷新间隔
func ProgressBarWithNewLine(col int, interval time.Duration) {
	bar := fmt.Sprintf("[%%-%vs]\n", col) // \x0d 回车
	baseBar(col, interval, bar)
}

// ProgressBarWithClearScreen 进度条，清屏后展示进度条
// col 进度条长度
// interval 进度条刷新间隔
func ProgressBarWithClearScreen(col int, interval time.Duration) {
	bar := fmt.Sprintf("\033[2J\033[H[%%-%vs]", col) // \x0c 清屏  \033[2J
	baseBar(col, interval, bar)
}

// ProgressBarRate 进度条，展示进度条
// col 进度条长度
// rate 进度条进度 0-1之间的小数
func ProgressBarRate(col int, rate float64) {
	bar := fmt.Sprintf("\u001B[2K\u001B[1G[%%-%vs]", col) //\033[2K 清除当前行； \033[1G 将光标移动到行首
	baseBarRate(col, rate, bar)
}

// ProgressBarRateWithNewLine 进度条，新行展示进度条
// col 进度条长度
// rate 进度条进度 0-1之间的小数
func ProgressBarRateWithNewLine(col int, rate float64) {
	bar := fmt.Sprintf("[%%-%vs]\n", col) // \x0d 回车
	baseBarRate(col, rate, bar)
}

// ProgressBarRateWithClearScreen 进度条，清屏后展示进度条
// col 进度条长度
// rate 进度条进度 0-1之间的小数
func ProgressBarRateWithClearScreen(col int, rate float64) {
	bar := fmt.Sprintf("\u001B[2J\u001B[H[%%-%vs]", col) // \x0c 清屏
	baseBarRate(col, rate, bar)
}

func baseBar(col int, interval time.Duration, bar string) {
	if col < 0 {
		return
	}
	for i := 0; i < col; i++ {
		fmt.Printf(bar, strings.Repeat("=", i)+">")
		time.Sleep(interval)
	}
	fmt.Printf(bar, strings.Repeat("=", col))
}
func baseBarRate(col int, rate float64, bar string) {
	if col < 0 || rate < 0 || rate > 1 {
		return
	}
	fCol := float64(col)
	f := rate * fCol
	fmt.Printf(bar, strings.Repeat("=", int(f))+">")

	if fCol-f <= fCol/100000 { // 误差小于0.001%,认为相等
		fmt.Printf(bar, strings.Repeat("=", col))
	}
}
