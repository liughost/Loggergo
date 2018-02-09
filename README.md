# Loggergo
## 功能：日志生成
* 指定输出目录和文件名
* 每天一个文件
* 可以控制是否输出到标准输出
* 支持多参数，但要求第一个参数必须是字符串
## 导入方法
* go get github.com/liughost/Loggergo
## 例子：
```
import (
	"github.com/liughost/Loggergo"
)
func main(){
  //设置初始参数
  //输出文件名为test.log，输出目录为当前目录的logs目录下
  //最后一个参数true代表写日志文件的同时输出到标准输出（例如：屏幕），false表示只输出到日志文件
  Loggergo.SetLog("test.log", "./logs", true)
  //
  prn := Loggergo.PrintLog
  //写入日志文件
  prn("日志输出测试：",123)
}
```
