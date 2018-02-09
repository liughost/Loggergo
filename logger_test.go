package Loggergo

import (
	"testing"
)

func Test_create(t *testing.T) {
	SetLog("test.log", "../../logs", true)
	PrintLog("abc")
	//g.Println("hello log", "sss", "www", 1, 2)
	// g = Instance()
	// g.Println("hello log", 333333)
	t.Log("Fine")
}
