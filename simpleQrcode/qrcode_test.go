package simpleQrcode

import (
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
	"os"
	"testing"
)

func TestQrcode(t *testing.T) {
	//定义二维码内容
	qrcode, _ := qr.Encode("https://www.cnblogs.com/hongyeci", qr.M, qr.Auto)
	//设置二维码的宽高
	qrcode, _ = barcode.Scale(qrcode, 256, 256)
	//创建qrcode.png文件
	file, _ := os.Create("qrcode2.png")
	//关闭文件流--defer-延迟关闭资源
	defer file.Close()
	//将二维码写入文件中
	if err := png.Encode(file, qrcode); err != nil {
		fmt.Println("生成二维码失败")
	} else {
		fmt.Println("生成二维码成功")
	}
}

//func TestQrcode(t *testing.T) {
//
//}
