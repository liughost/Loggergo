package Loggergo

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)


var (
	logName      = "ghost.log"
	logPath      = "../../log"
	_self        *GuttvLogger
	PrintConsole = true //是否显示到屏幕
)

type GuttvLogger struct {
	path     string
	fileName string
	logFile  string
	fp       *os.File
	locker   *sync.Mutex
}

func SetLog(name, path string, log_console bool) {
	logName = name
	logPath = path
	PrintConsole = log_console
	Instance(logName, logPath).Start()
}
func PrintLog(text string, args ...interface{}) {
	Instance(logName, logPath).Println(text, args)
}

func Instance(args ...string) *GuttvLogger {
	if _self == nil {
		//fmt.Println("create logger.")
		_self = new(GuttvLogger)

		_self.Start()
	}
	return _self
}

func (this *GuttvLogger) Start() {
	if this.locker == nil {
		this.locker = new(sync.Mutex)
	}
	this.path = logPath
	this.fileName = logName
	this.locker.Lock()
	this.CreateFile()
	this.locker.Unlock()
}

/**
创建文件，并绑定到logger
**/
func (this *GuttvLogger) CreateFile() {
	this.logFile = this.path + string(os.PathSeparator) + this.fileName
	this.HistoryFile()
	var err error
	//创建、打开文件
	this.fp, err = os.OpenFile(this.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	//defer this.fp.Close() //函数结束即立即关闭
	if err != nil {
		_, file, line, _ := runtime.Caller(4)
		_, file = path.Split(file)
		fmt.Println(fmt.Sprintf("%s(%d)\topening file: %s, error:%v", file, line, this.logFile, err))
	}
	//this.fp.WriteString("slkdjfksdfksf")
}

/**
修改文件名，发现同名文件改成追加
**/
func Rename(srcName, dstName string) {
	_, err := os.Stat(dstName)
	//fmt.Println("Stat ==========", err)
	if err != nil {
		//没有同名文件
		os.Rename(srcName, dstName)
		return
	}
	//合并文件
	fs, err := os.OpenFile(srcName, os.O_RDONLY, 0666)
	defer fs.Close()
	if err != nil {
		//打开失败，不合并
		return
	}
	//fmt.Println("open src", srcName)
	bReader := bufio.NewReader(fs)
	fd, err := os.OpenFile(dstName, os.O_RDWR|os.O_APPEND, 0666)
	defer fd.Close()
	if err != nil {
		//打开失败，不合并
		return
	}
	//fmt.Println("open dst", dstName)
	//bWriter := bufio.NewWriter(fd)
	//合并内容
	buffer := make([]byte, 1024)
	for {
		count, readErr := bReader.Read(buffer)
		//fmt.Println("read err", readErr)
		if readErr != nil {
			break
		}
		//fmt.Println("write", string(buffer[:count]))
		fd.Write(buffer[:count])
		//fmt.Println(wc, err)
	}
	//删除原文件
	os.Remove(srcName)
}

/**
检测历史文件，不在同一日则改名
返回：
true =不存在
false = 存在
**/
func (this *GuttvLogger) HistoryFile() (bool, error) {
	day, err := this.FileInfo()
	if err != nil {
		_, file, line, ok := runtime.Caller(4)
		if ok {
			_, file = path.Split(file)
		}
		fmt.Println(fmt.Sprintf("%s(%d)\tfile info error: %v", file, line, err))

		//文件打开失败
		return false, err
	} else {
		now := time.Now().Format("2006-01-02")
		//fmt.Println(now, day)
		//比较日期
		if now != day {
			//关闭文件
			if this.fp != nil {
				this.fp.Close()
				this.fp = nil
			}
			//改名称
			dstName := this.logFile + "_" + day
			Rename(this.logFile, dstName)
			fmt.Println("file:", this.logFile, ",change to ", dstName)
			return true, nil
		} else {
			return false, nil
		}
	}
}

func (this *GuttvLogger) FileInfo() (string, error) {
	f, err := os.Stat(this.logFile)
	if err != nil {
		return "", err
	}
	return f.ModTime().Format("2006-01-02"), nil
}

func (this *GuttvLogger) Println(text string, args ...interface{}) {
	for _, v := range args {
		str := fmt.Sprintf("%v", v)
		if len(str) > 2 {
			rs := []rune(str)
			//fmt.Println("aaaaa", str)
			str = string(rs[1 : len(rs)-1]) //带有汉字，需要使用rs的长度
			if len(str) > 0 {
				text += "\t" + str //fmt.Sprintf("%v", v)
			}
		}
	}
	//获取上一级的上一级
	_, file, line, ok := runtime.Caller(2)
	if ok {
		_, file = path.Split(file)
		text = fmt.Sprintf("%s(%d)\t%s", file, line, text)
	}

	this.Writer(text)
}

func (this *GuttvLogger) Writer(text string) {
	//启动多线程锁
	this.locker.Lock()
	//检查文件并改名
	change, _ := this.HistoryFile()
	if change || this.fp == nil {
		//创建新的日志文件
		this.CreateFile()
	}
	str := fmt.Sprintf("%s\t%s\r\n", time.Now().Format("2006-01-02 15:04:05"), text)
	if PrintConsole {
		fmt.Println(str)
	}
	//写入文件
	this.fp.WriteString(str)
	//解除锁
	this.locker.Unlock()
}
