package main

import (
	"fmt"
	"git.lcc.lib/core"
	//"time"
)

func main() {

	ini := core.NewIni("config/web.ini")
	fmt.Println(ini.GetItem("SERVER", "Host"))

	ini.ReLoad()
	fmt.Println(ini.GetItem("REDIS", "Host"))

	fmt.Println(ini.GetItem("LOGER", "Mask"))
	fmt.Println(ini.GetItem("LOGER", "Dir"))

	//time.Sleep(time.Second * 30)
}
