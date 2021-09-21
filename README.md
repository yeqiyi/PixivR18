# pixivR18

仿造项目<https://github.com/secriy/pixivrank>爬取p站r18每日排行榜，下载下来的图像保存在`./image/`下

## 食用方法

在项目根目录下添加`.cookies`文件，里面放的是你的p站cookie。

### 两个可选参数
| parameter | value | 
| :-: | :-: | 
|  -p | http代理地址，如`http://127.0.0.1:xxxx` |
|  -n | 下载涩图数量，默认为0 |

运行命令
```
go build -o pixivR18 ./cmd/ 
```
开冲
```
./pixivR18 -p http://127.0.0.1:xxxx -n x
```
