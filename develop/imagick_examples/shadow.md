package main

import "github.com/gographics/imagick/imagick"

func main() {
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	mw.ReadImage("/home/chen/61.png")
	width := mw.GetImageWidth()
	hight := mw.GetImageHeight()

	shadow := mw.Clone()

	p := imagick.NewPixelWand()
	defer p.Destroy()
	p.SetColor("black")
	
	shadow.SetImageBackgroundColor(p)
	shadow.ShadowImage(80, 3, 6, 6) // 2,6倍数且x=y较好
	shadow.CompositeImage(mw, imagick.COMPOSITE_OP_OVER, 0, 0)
	shadow.CropImage(width, hight, 0, 0)
	shadow.WriteImage("/home/chen/shdow.png")
}

