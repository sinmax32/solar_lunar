package main

import (
	"fmt"
	"testing"
)

func TestSolarLunar(t *testing.T) {
	sl := new(SolarLunar)

	y_1 := 2013
	m_1 := 3
	d_1 := 16

	y_2, m_2, d_2, l_2 := sl.Solar2Lunar(y_1, m_1, d_1)

	fmt.Println("公历", y_1, "年", m_1, "月", d_1, "日")

	if y_2 < 1900 {
		fmt.Println("该日期的转换不被支持 或 该日期无法识别")
	} else {
		if l_2 == 1 {
			fmt.Println("农历", y_2, "年 闰", m_2, "月", d_2, "日")
		} else {
			fmt.Println("农历", y_2, "年", m_2, "月", d_2, "日")
		}
	}

	fmt.Println("------------------------")

	y_1 = 2012
	m_1 = 8
	d_1 = 22
	l_2 = 0

	y_2, m_2, d_2 = sl.Lunar2Solar(y_1, m_1, d_1, l_2)

	if sl.leap == 1 {
		fmt.Println("农历", y_1, "年 闰", m_1, "月", d_1, "日")
	} else {
		fmt.Println("农历", y_1, "年", m_1, "月", d_1, "日")
	}

	fmt.Println("公历", y_2, "年", m_2, "月", d_2, "日")

	//后续考虑加入农历日的天干地支计算

}
