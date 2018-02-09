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
  prn := Loggergo.PrintLog
  prn("日志输出测试：",123)
}
```
