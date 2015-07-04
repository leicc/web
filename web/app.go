package main

import (
	"fmt"
	"git.lcc.lib/core"
	"git.lcc.lib/verify"
)

func main() {
	log := core.NewLoger("./cache/log", 2)
	fmt.Println(log)
	core.Log(1, "1111", "aaaaaaa", "cccccccccccc")
	core.Log(2, "2222", "aaaaaaa", "cccccccccccc")
	core.Log(3, "2222", "aaaaaaa", "cccccccccccc")

	core.Logf(1, "%s-%s-%s", "333333333", "aaaaaaa", "cccccccccccc")
	core.Logf(4, "%s-%s-%s", "333333333", "aaaaaaa", "cccccccccccc")

	fmt.Println(verify.Email("lchenchun@sina.com.12"))

	fmt.Println(verify.Phone("135140768061"))

	fmt.Println(verify.Numeric("13aa516"))

	fmt.Println(verify.Alpha("13516"))

	fmt.Println(verify.Ip("1a1.0.11.0"))
}
