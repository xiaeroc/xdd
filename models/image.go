package models

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
)

// 创建宽 500 高 300 的图片
const (
	w = 500
	h = 300
)

func main() {

	//返回一个矩形
	rectangle := image.Rect(0, 0, w, h)
	rgba := image.NewRGBA(rectangle)
	// 设置一条白色的斜线
	for y := 0; y < h; y++ {
		rgba.Set(y, y, color.RGBA{255, 255, 255, 1})
	}
	// 文字水印
	fontbyte, err := ioutil.ReadFile("./simhei.ttf")
	if err != nil {
		fmt.Println("ioutil.ReadFile error : ", err)
		return
	}
	font, err := freetype.ParseFont(fontbyte)
	if err != nil {
		fmt.Println("freetype.ParseFont error : ", err)
		return
	}
	// 创建一个新的上下文
	context := freetype.NewContext()
	context.SetDPI(70)                                             // 设置屏幕分辨率，单位为每英寸点数。
	context.SetClip(rgba.Bounds())                                 //设置用于绘制的剪辑矩形。
	context.SetDst(rgba)                                           //设置绘制操作的目标图像。
	context.SetFont(font)                                          //设置用于绘制文本的字体。
	context.SetFontSize(16)                                        //以点为单位设置字体大小(如“12点字体”)。
	context.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 1})) //设置用于绘制操作的源图像
	pt := freetype.Pt(10, 260)                                     //从一个以像素度量的坐标对转换为一个固定的点
	context.DrawString("文字水印", pt)
	// 图片水印
	img, _ := os.Open("./img.jpg")
	defer img.Close()
	img1, _ := jpeg.Decode(img) //读取一个JPEG图像并将其作为image.Image返回
	offset := image.Pt(300, 10)
	draw.Draw(rgba, img1.Bounds().Add(offset), img1, image.ZP, draw.Over)

	//创建图片
	file, err := os.Create("./text.jpg")
	if err != nil {
		fmt.Println("os.Open error : ", err)
		return
	}
	// 将图像写入file
	//&jpeg.Options{100} 取值范围[1,100]，越大图像编码质量越高
	jpeg.Encode(file, rgba, &jpeg.Options{100})
	defer file.Close()
}
