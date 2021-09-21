package pixivR18

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"sync"
)

type PixivSlot struct {
	Client *http.Client
	Cookie string
	proxy string
	Numbers int //所要下载图片的数量
}

type illust struct {
	Contents []struct {
		IllustID int `json:"illust_id"`
	} `json:"contents"`
}

//获取当日排行榜ID列表
func (p *PixivSlot) RankIDList()([]int,error){
	url:="https://www.pixiv.net/ranking.php?mode=daily_r18&content=illust&p=1&format=json"
	req,err:=http.NewRequest("GET",url,nil)
	if err!=nil{
		return nil,err
	}
	req.Header.Add("accept","text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("user-agent","Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15")
	req.Header.Add("cookie",p.Cookie)
	resp,err:=p.Client.Do(req)
	if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return nil,err
	}
	illust:=&illust{}
	if err:=json.Unmarshal(body,illust);err!=nil{
		return nil,err
	}
	ids :=make([]int,len(illust.Contents))
	for k,v :=range illust.Contents{
		ids[k]=v.IllustID
	}
	return ids,err
}

//通过RankIDList返回的id列表获取图片url
func (p *PixivSlot) ImageURL(id int) string{
	url:=fmt.Sprintf("https://www.pixiv.net/artworks/%d",id)
	req,err:=http.NewRequest("GET",url,nil)
	if err!=nil{
		return ""
	}
	resp,err:=p.Client.Do(req)
	if err!=nil{
		return ""
	}
	defer resp.Body.Close()
	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return ""
	}
	reg:=regexp.MustCompile(`(?:"original":")(.*?)(?:"\})`)
	str:=reg.FindStringSubmatch(string(body))

	if len(str)<2{
		return ""
	}
	fmt.Println("str =",str[1])
	return str[1]
}

//通过图片url下载图片到指定路径
func (p *PixivSlot) Download(url string,id int) error{
	pathStr:="./image"

	if err:=os.MkdirAll(pathStr,0777);err!=nil{
		return err
	}
	req,err:=http.NewRequest("GET",url,nil)
	if(err!=nil){
		return err
	}
	req.Header.Add("referer",fmt.Sprint("https://www.pixiv.net/artworks/%d",id))
	resp,err:=p.Client.Do(req)//发起http请求并获取响应包
	if err!=nil{
		return err
	}
	defer resp.Body.Close()//防止内存泄漏


	file,err:=os.Create(pathStr+"/"+path.Base(req.URL.Path))
	if err!=nil{
		return err
	}
	wt:=bufio.NewWriter(file)
	defer file.Close()
	_,err=io.Copy(file,resp.Body)
	if err !=nil{
		return err
	}
	fmt.Println(file.Name()+" download complete")
	wt.Flush()
	return nil
}

func (p *PixivSlot) readCookies() error{
	file,err:=os.Open("./.cookies")
	if err!=nil{
		return err
	}
	defer file.Close()
	text,err:=ioutil.ReadAll(file)
	if err!=nil{
		return err
	}
	p.Cookie=string(text)
	return nil
}

func (p *PixivSlot) Task(){
	if err:=p.readCookies();err!=nil{
		fmt.Println("Failed to read cookies: "+err.Error())
		return
	}

	IdLIst,err:=p.RankIDList()
	if err!=nil{
		fmt.Println("Failed to get ranklist:"+err.Error())
		return
	}
	wg:=&sync.WaitGroup{}

	if p.Numbers>0&&p.Numbers<=len(IdLIst){
		IdLIst=IdLIst[:p.Numbers]
	}

	for k,v :=range IdLIst{
		wg.Add(1)
		go func(k,v int){
			defer wg.Done()
			url:=p.ImageURL(v)
			if url!=""{
				err:=p.Download(url,v)
				if err!=nil{
					fmt.Println(err)
				}
			}
		}(k,v)
	}
	wg.Wait()
}

//创建task实例
func (p *PixivSlot) RunTask(proxy string,numbers int) {
	if proxy == "" {
		p.Client = NewClient()
	} else {
		p.Client = NewClientWithPorxy(proxy)
	}
	p.Numbers = numbers
	p.Task()
}


