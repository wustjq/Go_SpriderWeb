package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func Http_Get(url string, i int) (res string, err error) {
	resp, err1 := http.Get(url)
	//将错误传出交给外界处理
	if err1 != nil {
		err = err1
		return
	}

	defer resp.Body.Close()
	//将得到数据读到buf里面
	buf := make([]byte, 4096)
	//循环读取网页数据
	for {
		//表示剩余还有多少没读
		n, err2 := resp.Body.Read(buf)
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		// if err2 != nil {
		// 	fmt.Println("err2=", err2)//err=eof
		// }
		if n == 0 {
			fmt.Println("恭喜您，爬取第" + strconv.Itoa(i) + "页成功")
			break
		}
		res += string(buf[:n])
	}
	return
}

func DoUrl(i int, page chan int) {
	url := "https://bbs.hupu.com/46683153-" + strconv.Itoa(i) + ".html"
	//调用该函数进行处理url读写 将数据读给调用者
	res, err := Http_Get(url, i)
	if err != nil {
		fmt.Println("get error .....", err)
		return
	}
	file, err1 := os.Create("第" + strconv.Itoa(i) + "个" + ".html")
	if err1 != nil {
		fmt.Println("create error .....", err)
		return
	}
	file.WriteString(res)
	file.Close()

	page <- i
}

func handler(start, end int) {
	//爬虎扑
	//https://bbs.hupu.com/46683153.html //加个-1
	//https://bbs.hupu.com/46683153-2.html
	//https://bbs.hupu.com/46683153-3.html
	page := make(chan int)

	for i := start; i <= end; i++ {
		go DoUrl(i, page)
	}

	for i := start; i <= end; i++ {
		<-page
	}

}

func main() {
	var start, end int
	fmt.Print("请输入你要爬取的开始页码(>=1):")
	fmt.Scan(&start)

	fmt.Print("请输入你要爬取的终止页码(>=1):")
	fmt.Scan(&end)

	handler(start, end)
}
