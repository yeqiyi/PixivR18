package main

import(
	"flag"
	"pixivR18"
)

var(
	proxy      = flag.String("p", "", "Proxy address, such as 'http://127.0.0.1:1080'")
	numbers    = flag.Int("n", 0, "The number of images to download.")
)

func main(){
	flag.Parse()
	p:=new(pixivR18.PixivSlot)
	p.RunTask(*proxy,*numbers)
}