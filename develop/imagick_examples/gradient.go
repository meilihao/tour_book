package main

import (
	"github.com/gographics/imagick/imagick"
)

func main() {
	//画布
	bg := imagick.NewMagickWand()
	bg.SetSize(400, 200)
	//bg.ReadImage("gradient:#A02B2B-#126B27") //渐变
    bg.ReadImage("gradient:gray70-gray30")
	
	p := imagick.NewPixelWand()
	p.SetColor("none")
	//bg.RotateImage(p, -60) //旋转

	bg.SetImageFormat("png")
	bg.WriteImage("/home/chen/bg.png")
}

