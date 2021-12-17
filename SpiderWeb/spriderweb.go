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
		//err不是空且没有到文件尾部
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}

		res += string(buf[:n])
	}
	return
}

func handler(start int, end int) {
	fmt.Printf("正在爬取第%d到页%d\n", start, end)

	//循环爬取每页数据
	for i := start; i <= end; i++ {
		url := "https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn=" + strconv.Itoa((i-1)*50)
		//有了url网页接下来访问网页对应的内容
		res, err := HttpGet(url, i)
		if err != nil {
			fmt.Println("http get error ....", err)
			continue
		}
		//fmt.Println("res=", res)
		//将读到网页文件储存
		file, err1 := os.Create("第" + strconv.Itoa(i) + "页" + ".html")
		if err1 != nil {
			fmt.Println("os create error ....", err)
			continue
		}

		file.WriteString(res)
		//这不能用defer关闭 因为你的文件要创建关闭多次
		file.Close() //保存一个 关闭一个
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
