package MyRegex

import "regexp"

var intNumber, _ = regexp.Compile("^[+-]{0,1}[0-9]+$")

func TestNumber(s string) bool {
	return number.MatchString(s)
}

var number, _ = regexp.Compile("^[+-]{0,1}[0-9.]+$")

func TestIntNumber(s string) bool {
	return intNumber.MatchString(s)
}

var scienceNumber, _ = regexp.Compile("^[0-9.]+e[0-9.]+$")

func TestScienceNumber(s string) bool {
	return scienceNumber.MatchString(s)
}

var numberAndNotSence, _ = regexp.Compile("^[0-9.e, \\-]+$")

func TestNumberAndNotSence(s string) bool {
	return numberAndNotSence.MatchString(s)
}

var noWord, _ = regexp.Compile("^[0-9.e,'\"`:%$#@!&*^+\\/;? \\-]+$")

func TestNoWord(s string) bool {
	return noWord.MatchString(s)
}
