package main

//百度贴吧爬取

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HttpGet(url string, i int) (res string, err error) {
	resp, err1 := http.Get(url)
	//相当于把错误传出去给外部判断
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	//循环读取网页数据 把数据传输给调用者 读到buf里面
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("爬取第" + strconv.Itoa(i) + "页成功！")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		res += string(buf[:n])
	}
	return
}

//爬取单个页面的函数
func SpriderPage(i int, page chan int) {
	url := "https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn=" + strconv.Itoa((i-1)*50)
	//有了url网页接下来访问网页对应的内容
	res, err := HttpGet(url, i)
	if err != nil {
		fmt.Println("http get error ....", err)
		return
	}
	//fmt.Println("res=", res)
	//将读到网页文件储存
	file, err1 := os.Create("第" + strconv.Itoa(i) + "页" + ".html")
	if err1 != nil {
		fmt.Println("os create error ....", err)
		return
	}
	file.WriteString(res)
	//这不能用defer关闭 因为你的文件要创建关闭多次
	file.Close() //保存一个 关闭一个
	//这里没读取完成一次 向page里传输  告诉我主go程 已经爬取完成
	page <- i
}

func handler(start int, end int) {
	fmt.Printf("正在爬取第%d到页%d\n", start, end)

	//定义一个chan不要让主go程先死亡
	page := make(chan int)

	//循环爬取每页数据
	for i := start; i <= end; i++ {
		//利用go程去承载每一页网页爬取 注意这里不能让主go程死亡 不然程序结束了
		go SpriderPage(i, page)
		//把page放这里就是严格的一个一个去读取  你非page读走才能创建go程 所以比放在外面循环慢
	}

	for i := start; i <= end; i++ {
		<-page
	}
}

func main() {
	//指定爬取起始和终止页面
	var start, end int

	fmt.Print("请输入起始页面(>=1)：")
	fmt.Scan(&start)
	fmt.Print("请输入终止页面(>=start)：")
	fmt.Scan(&end)

	handler(start, end)
}
