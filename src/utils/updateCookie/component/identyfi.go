package component

import (
	"fmt"
	"os/exec"
)

func Identify(isInServer bool) (string, error) {
	if isInServer {
		cmd := exec.Command("python3", "./ocr/main.py")
		outPut, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("请检查是否是服务器环境！")
			return "", err
		}
		str := string(outPut)
		str = str[len(str)-5 : len(str)-1]
		return str, nil
	} else {
		cmd := exec.Command("python", "./ocr/main.py")
		outPut, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("请检查是否是测试环境以及虚拟环境是否激活，关键包：ddddocr\n激活指令：./ocr/.venv/Scripts/activate")
			return "", err
		}
		str := string(outPut)
		str = str[len(str)-6 : len(str)-2]
		return str, nil
	}
}
