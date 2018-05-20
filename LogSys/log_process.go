package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
	"regexp"
)

type Reader interface {
	Read(rc chan []byte)
}
type Writer interface {
	Write(wc chan *LogMessage)
}

type ReadFromFile struct {
	LogFilePath string // 日志文件路径
}

func (r *ReadFromFile) Read(rc chan []byte)  {
	file, err := os.Open(r.LogFilePath)
	handleErr(err)

	_, seekErr := file.Seek(0, 2)
	handleErr(seekErr)

	buffReader := bufio.NewReader(file)

	for {
		data, err := buffReader.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500*time.Millisecond)
			continue
		}else if err != nil {
			handleErr(err)
		}

		rc <- data[:len(data)-1]  // 去掉末尾换行符
	}

}

type WriteToDB struct {
	dbConfig string  // 数据库配置信息
}

func (w *WriteToDB) Write(wc chan *LogMessage)  {
	for line := range wc {
		fmt.Println(line)
	}
}

type LogProcess struct {
	rc chan []byte  // 读文件通道
	wc chan *LogMessage  // 写文件通道
	reader Reader
	writer Writer
}

func (l *LogProcess) Process()  {
	// 解析模块

	for line := range l.rc {
		//使用正则提取出需要的信息
		r := regexp.MustCompile(`([a-zA-Z])([^<])`)
		substr := r.FindStringSubmatch(string(line))
		if len(substr)!=10 {
			panic(fmt.Sprintf("regexp ERR!"))
			continue
		}

		//从substr中获得需要信息 塞入LogMessage中
		lm := &LogMessage{}

		l.wc <- lm
	}

}

// 提取日志中的关键信息
type LogMessage struct {
	name string
	num int
}

func handleErr(err error)  {
	if err != nil {
		panic(fmt.Sprintf("Error! %s", err.Error()))
	}
}

func main() {

	r := &ReadFromFile{
		"D:\\GOWORKSPACE\\src\\LogSys\\cln.log",
	}
	w := &WriteToDB{
		"root&root@123",
	}
	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan *LogMessage),
		reader: r,
		writer: w,
	}

	go lp.reader.Read(lp.rc)
	go lp.Process()
	go lp.writer.Write(lp.wc)

	time.Sleep(30*time.Second)
}
