package verify

import (
	"regexp"
)

var (
	email   = regexp.MustCompile(`^[\w\.\_]{2,}@([\w\-]+\.){1,}\.[a-z]$`)
	phone   = regexp.MustCompile(`^(\+86)?1[3-9][0-9]{9}$`)
	numeric = regexp.MustCompile(`^[0-9]{1,}$`)
	alpha   = regexp.MustCompile(`^[a-zA-Z0-9]{1,}$`)
	ip      = regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
)

func Email(vdata string) bool {
	return email.MatchString(vdata)
}

func Phone(vdata string) bool {
	return phone.MatchString(vdata)
}

func Numeric(vdata string) bool {
	return numeric.MatchString(vdata)
}

func Alpha(vdata string) bool {
	return alpha.MatchString(vdata)
}

func Ip(vdata string) bool {
	return ip.MatchString(vdata)
}
