package core

import (
	"regexp"
)

func IsEmail(vdata string) bool {
	email := regexp.MustCompile(`^[\w\.\_]{2,}@([\w\-]+\.){1,}\.[a-z]$`)
	return email.MatchString(vdata)
}

func IsPhone(vdata string) bool {
	phone := regexp.MustCompile(`^(\+86)?1[3-9][0-9]{9}$`)
	return phone.MatchString(vdata)
}

func IsNumeric(vdata string) bool {
	numeric := regexp.MustCompile(`^[0-9\.]{1,}$`)
	return numeric.MatchString(vdata)
}

func IsAlphaNumeric(vdata string) bool {
	alpha := regexp.MustCompile(`^[a-zA-Z0-9]{1,}$`)
	return alpha.MatchString(vdata)
}

func IsIp(vdata string) bool {
	ip := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
	return ip.MatchString(vdata)
}

func IsDate(vdata string) bool {
	date := regexp.MustCompile(`^[0-9]{4}\-[0-9]{2}\-[0-9]{2}$`)
	return date.MatchString(vdata)
}

func IsIdCard(vdata string) bool {
	idcard := regexp.MustCompile(`(^\d{15}$)|(^\d{17}[\d|x|X]$)`)
	return idcard.MatchString(vdata)
}

func IsUserName(vdata string) bool {
	user := regexp.MustCompile(`^[\w\@\-\.]{3,}$`)
	return user.MatchString(vdata)
}

func IsPasswd(vdata string) bool {
	//要求六位以上且还有英文和字母
	isnull := regexp.MustCompile(`^[^\s]{6,}$`)
	isalpha := regexp.MustCompile(`[a-zA-Z]`)
	isnumeric := regexp.MustCompile(`[0-9]`)
	return isnull.MatchString(vdata) && isnumeric.MatchString(vdata) && isalpha.MatchString(vdata)
}
