package main

import (
	"time"
)

var (
	//农历基本参照日
	BASE_LUNAR_DATE time.Time = time.Date(1900, time.Month(1), 31, 0, 0, 0, 0, time.UTC)
	//农历大小月数据
	LUNAR_DATA = []int{19416, 19168, 42352, 21717, 53856, 55632, 25940, 22176, 39632, 21970, 19168, 42422, 42192, 53840, 53845, 46400, 54944, 44450, 38320, 18807, 18800, 42160, 46261, 27216, 27968, 43860, 11104, 38256, 21234, 18800, 25958, 54432, 59984, 28309, 23248, 11104, 34531, 37600, 51415, 51536, 54432, 55462, 46416, 22176, 42420, 9680, 37584, 53938, 43344, 46423, 27808, 46416, 21333, 19872, 42416, 17779, 21168, 43432, 59728, 27296, 44710, 43856, 19296, 43748, 42352, 21088, 62051, 55632, 23383, 22176, 38608, 19925, 19152, 42192, 54484, 53840, 54616, 46400, 46752, 38310, 38320, 18864, 43380, 42160, 45690, 27216, 27968, 44870, 43872, 38256, 19189, 18800, 25776, 29859, 59984, 27480, 21952, 43872, 38613, 37600, 51552, 55636, 54432, 55888, 30034, 22176, 43959, 9680, 37584, 51893, 43344, 46240, 47780, 44368, 21977, 19360, 42416, 20854, 21168, 43312, 31060, 27296, 44368, 23378, 19296, 42726, 42208, 53856, 60005, 54576, 23200, 30371, 38608, 19195, 19152, 42192, 53430, 53840, 54560, 56645, 46496, 22224, 21938, 18864, 42359, 42160, 43600, 45653, 27936, 44448, 19299, 37744, 18936, 18800, 25776, 26790, 59984, 27424, 42692, 43744, 41696, 53987, 51552, 54615, 54432, 55888, 23893, 22176, 42704, 21972, 21200, 43448, 43344, 46240, 46758, 44368, 21920, 43940, 42416, 21168, 45683, 26928, 29495, 27296, 44368, 19285, 19296, 42352, 21732, 53600, 59752, 54560, 55968, 27302, 22224, 19168, 43476, 41680, 53584, 62034, 54560}
	//闰月数据
	LUNAR_LEAP_DAY = []int32{67379328, 344986112, 131072, 2048, 16786468, 128, 65792}
)

type SolarLunar struct {
	solarYear  int //公历年
	solarMonth int //公历月
	solarDate  int //公历日
	lunarYear  int //农历年
	lunarMonth int //农历月
	lunarDate  int //农历日
	leap       int //闰月
}

//公历日是否可转换到农历。参数是公历年月日
func (p *SolarLunar) canConvert2Lunar(y, m, d int) bool {
	if y < 1900 || y > 2100 {
		//不支持
		return false
	}

	if y == 1900 && m == 1 {
		//不支持
		return false
	}

	if 1 <= m && m <= 12 {

	} else {
		return false
	}

	//golang的Date内部不会作日期验证，所以这里要验证一下该日期是否有效
	if time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC).Day() != d {
		return false
	}

	return true
}

//农历日是否可转换到农历。参数是农历年月日，leap是否闰月
func (p *SolarLunar) canConvert2Solar(y, m, d, leap int) bool {

	if y < 1900 || y > 2100 {
		//不支持
		return false
	}

	if 1 <= m && m <= 12 {

	} else {
		return false
	}
	return true
}

//设置公历日期
func (p *SolarLunar) SetSolarDate(y, m, d int) {
	p.solarYear = y
	p.solarMonth = m
	p.solarDate = d
}

//设置农历日期
func (p *SolarLunar) SetLunarDate(y, m, d, leap int) {
	p.lunarYear = y
	p.lunarMonth = m
	p.lunarDate = d
	p.leap = leap
}

//获取闰月天数
func daysOfYearLeap(y int) int {
	yearItem := int32(y / 32)
	if yearItem >= int32(len(LUNAR_LEAP_DAY)) {
		return 0
	}
	yearItem = int32(LUNAR_LEAP_DAY[yearItem])
	yearBit := y % 32

	for i := 1; i < yearBit; i++ {
		yearItem <<= 1
	}
	if yearItem < 0 {
		return 30
	}
	return 29
}

//转换为农历
func (p *SolarLunar) ToLunar() bool {

	if !p.canConvert2Lunar(p.GetSolarDate()) {
		return false
	}

	t := time.Date(p.solarYear, time.Month(p.solarMonth), p.solarDate, 0, 0, 0, 0, time.UTC)
	offsetDay := int((t.Unix() - BASE_LUNAR_DATE.Unix()) / 86400)

	var y int
	var short16 int16
	var leapMonth int
	var leapMday int

	for y = 0; y < len(LUNAR_DATA); y++ {
		short16 = int16(LUNAR_DATA[y])
		leapMonth = int(short16 & 0xF)
		for j := 0; j < 12; j++ {
			if leapMonth > 0 && leapMonth == j {
				leapMday = daysOfYearLeap(y)
				offsetDay -= leapMday
				if offsetDay <= 0 {
					offsetDay += leapMday
					p.lunarYear = 1900 + y
					p.lunarMonth = j
					p.lunarDate = offsetDay + 1
					p.leap = 1
					return true
				}
			}

			if short16 < 0 {
				offsetDay -= 30
			} else {
				offsetDay -= 29
			}

			if offsetDay <= 0 {
				if short16 < 0 {
					offsetDay += 30
				} else {
					offsetDay += 29
				}

				p.lunarYear = 1900 + y
				p.lunarMonth = j + 1
				p.lunarDate = offsetDay + 1

				return true
			}
			short16 <<= 1
		}
	}

	return false
}

//转换为公历
func (p *SolarLunar) ToSolar() bool {
	if !p.canConvert2Solar(p.GetLunarDate()) {
		return false
	}

	var offsetDay int
	var short16 int16
	var leapMonth int

	//计算前面年的天数的总和
	var y int
	for y = 0; y < p.lunarYear-1900; y++ {
		short16 = int16(LUNAR_DATA[y])
		leapMonth = int(short16 & 0xF)
		for j := 0; j < 12; j++ {
			if short16 < 0 {
				offsetDay += 30
			} else {
				offsetDay += 29
			}
			short16 <<= 1
		}

		if leapMonth > 0 {
			offsetDay += daysOfYearLeap(y)
		}
	}

	//计算这一年的天数的总和
	short16 = int16(LUNAR_DATA[p.lunarYear-1900])
	leapMonth = int(short16 & 0xF)
	var j int
	for j = 1; j < p.lunarMonth; j++ {
		if short16 < 0 {
			offsetDay += 30
		} else {
			offsetDay += 29
		}

		if leapMonth > 0 && leapMonth == j {
			offsetDay += daysOfYearLeap(y)
		}
		short16 <<= 1
	}

	if p.leap == 1 && leapMonth > 0 && leapMonth == j {
		if short16 < 0 {
			offsetDay += 30
		} else {
			offsetDay += 29
		}
	} else {
		p.leap = 0
	}
	offsetDay += p.lunarDate - 1

	t := time.Unix(BASE_LUNAR_DATE.Unix()+int64(offsetDay*86400), 0).UTC()

	var m time.Month
	p.solarYear, m, p.solarDate = t.Date()
	p.solarMonth = int(m)

	return true
}

//获取农历年、月、日、闰
func (p *SolarLunar) GetLunarDate() (int, int, int, int) {
	return p.lunarYear, p.lunarMonth, p.lunarDate, p.leap
}

//获取公历年、月、日
func (p *SolarLunar) GetSolarDate() (int, int, int) {
	return p.solarYear, p.solarMonth, p.solarDate
}

//公历日期转换为农历日
func (p *SolarLunar) Solar2Lunar(y, m, d int) (int, int, int, int) {
	p.SetSolarDate(y, m, d)
	if !p.ToLunar() {
		return -1, 0, 0, 0
	}
	return p.GetLunarDate()
}

//农历日期转换为公历日
func (p *SolarLunar) Lunar2Solar(y, m, d, leap int) (int, int, int) {
	p.SetLunarDate(y, m, d, leap)
	if !p.ToSolar() {
		return -1, 0, 0
	}
	return p.GetSolarDate()
}

//重置属性。适合Single模式
func (p *SolarLunar) Reset() {
	p.solarYear = 0
	p.solarMonth = 0
	p.solarDate = 0
	p.lunarYear = 0
	p.lunarMonth = 0
	p.lunarDate = 0
}
