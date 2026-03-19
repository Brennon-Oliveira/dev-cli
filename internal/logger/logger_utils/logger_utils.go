package loggerutils

import (
	"math"
	"strings"
)

var RegularBlackColor = "\x1b[30m"
var RegularRedColor = "\x1b[31m"
var RegularGreenColor = "\x1b[32m"
var RegularYellowColor = "\x1b[33m"
var RegularBlueColor = "\x1b[34m"
var RegularPurpleColor = "\x1b[35m"
var RegularCyanColor = "\x1b[36m"
var RegularWhiteColor = "\x1b[37m"

var BoldBlackColor = "\x1b[1;30m"
var BoldRedColor = "\x1b[1;31m"
var BoldGreenColor = "\x1b[1;32m"
var BoldYellowColor = "\x1b[1;33m"
var BoldBlueColor = "\x1b[1;34m"
var BoldPurpleColor = "\x1b[1;35m"
var BoldCyanColor = "\x1b[1;36m"
var BoldWhiteColor = "\x1b[1;37m"

var UnderlineBlackColor = "\x1b[4;30m"
var UnderlineRedColor = "\x1b[4;31m"
var UnderlineGreenColor = "\x1b[4;32m"
var UnderlineYellowColor = "\x1b[4;33m"
var UnderlineBlueColor = "\x1b[4;34m"
var UnderlinePurpleColor = "\x1b[4;35m"
var UnderlineCyanColor = "\x1b[4;36m"
var UnderlineWhiteColor = "\x1b[4;37m"

var BackgroundBlackColor = "\x1b[40m"
var BackgroundRedColor = "\x1b[41m"
var BackgroundGreenColor = "\x1b[42m"
var BackgroundYellowColor = "\x1b[43m"
var BackgroundBlueColor = "\x1b[44m"
var BackgroundPurpleColor = "\x1b[45m"
var BackgroundCyanColor = "\x1b[46m"
var BackgroundWhiteColor = "\x1b[47m"

var HighIntensityBlackColor = "\x1b[90m"
var HighIntensityRedColor = "\x1b[91m"
var HighIntensityGreenColor = "\x1b[92m"
var HighIntensityYellowColor = "\x1b[93m"
var HighIntensityBlueColor = "\x1b[94m"
var HighIntensityPurpleColor = "\x1b[95m"
var HighIntensityCyanColor = "\x1b[96m"
var HighIntensityWhiteColor = "\x1b[97m"

var BoldHighIntensityBlackColor = "\x1b[1;90m"
var BoldHighIntensityRedColor = "\x1b[1;91m"
var BoldHighIntensityGreenColor = "\x1b[1;92m"
var BoldHighIntensityYellowColor = "\x1b[1;93m"
var BoldHighIntensityBlueColor = "\x1b[1;94m"
var BoldHighIntensityPurpleColor = "\x1b[1;95m"
var BoldHighIntensityCyanColor = "\x1b[1;96m"
var BoldHighIntensityWhiteColor = "\x1b[1;97m"

var HighIntensityBackgroundBlackColor = "\x1b[100m"
var HighIntensityBackgroundRedColor = "\x1b[101m"
var HighIntensityBackgroundGreenColor = "\x1b[102m"
var HighIntensityBackgroundYellowColor = "\x1b[103m"
var HighIntensityBackgroundBlueColor = "\x1b[104m"
var HighIntensityBackgroundPurpleColor = "\x1b[105m"
var HighIntensityBackgroundCyanColor = "\x1b[106m"
var HighIntensityBackgroundWhiteColor = "\x1b[107m"

var ResetColor = "\x1b[0m"

var BoldStyle = "\x1b[1m"
var ItalicStyle = "\x1b[3m"
var UnderlineStyle = "\x1b[4m"
var StrikethroughStyle = "\x1b[9m"

var BlockDivider = "="
var SimplesDivider = "-"

func HeaderDivider(title string) string {
	totalSize := 40
	titleSize := len(title)

	rest := totalSize - titleSize

	var before int

	if rest%2 == 0 {
		middle := rest / 2
		before = middle
	} else {
		middle := math.Floor(float64(rest) / 2)
		before = int(middle)
	}

	var str strings.Builder

	str.Grow(totalSize)

	var startAfter = before + titleSize

	for i := range totalSize {
		if i < startAfter && i >= before {
			str.WriteByte(title[i-before])
			continue
		}
		str.WriteString(BlockDivider)
	}

	return str.String()
}
