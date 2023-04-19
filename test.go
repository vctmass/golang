package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Comment 结构体
type Comment struct {
	Email string `json:"email"`
}

// baseURL 存储 API 的基础 URL
const baseURL = "https://jsonplaceholder.typicode.com/posts/"

func main() {
	// 开始时间
	startTime := time.Now()
	fmt.Printf("开始时间: %v\n", startTime)

	// 创建等待组，每个 URL 对应一个 goroutine
	var wg sync.WaitGroup
	wg.Add(100)

	// 电子邮件地址
	var emails []string

	// 创建通道，传输邮件地址
	emailChan := make(chan string)

	// 协程异步处理请求
	for i := 1; i <= 100; i++ {
		go func(postID int) {
			// 处理完一个URL后，减1
			defer wg.Done()

			// 构造URL
			url := fmt.Sprintf("%s%d/comments", baseURL, postID)

			// 发送HTTP GET请求，获取响应
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("请求失败： %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()

			// 解码JSON
			var comments []Comment
			err = json.NewDecoder(resp.Body).Decode(&comments)
			if err != nil {
				fmt.Printf("数据解析失败 %s: %v\n", url, err)
				return
			}

			// 将每个评论中的电子邮件地址发送到通道中
			for _, comment := range comments {
				emailChan <- comment.Email
			}
		}(i) // 传入URL的ID，遍历1到100
	}

	// 处理通道中的电子邮件地址，将它们存储到切片中
	go func() {
		for email := range emailChan {
			// fmt.Printf("邮件： %v\n", email)
			emails = append(emails, email)

		}
	}()

	// 等待所有goroutine处理完毕，关闭通道
	wg.Wait()
	close(emailChan)

	// 打开文件，准备写入电子邮件地址
	file, err := os.Create("emails.txt")
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer file.Close()

	// 将电子邮件地址写入文件
	for _, email := range emails {
		_, err := file.WriteString(email + "\n")
		if err != nil {
			fmt.Printf("写入文件失败: %v\n", err)
			return
		}
	}

	// 记录结束时间
	endTime := time.Now()
	fmt.Printf("结束时间: %v\n", endTime)

	// 打印出总共的电子邮件地址数和程序运行时间
	fmt.Printf("emails总数: %d\n", len(emails))
	fmt.Printf("程序完成时间为： %v seconds.\n", endTime.Sub(startTime).Seconds())
}
