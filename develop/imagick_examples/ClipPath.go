
package main

import (
	"github.com/gographics/imagick/imagick"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	p := imagick.NewPixelWand()
	defer p.Destroy()
	p.SetColor("red")

	p1 := imagick.NewPixelWand()
	defer p1.Destroy()
	p1.SetColor("green")

	p2 := imagick.NewPixelWand()
	defer p2.Destroy()
	p2.SetColor("blue")

	draw := imagick.NewDrawingWand()
	draw.SetStrokeColor(p)
	draw.SetFillColor(p1)
	draw.SetStrokeOpacity(1)
	draw.SetStrokeWidth(2)

	clipPathName := "testClipPath"

	draw.PushClipPath(clipPathName)
	draw.Rectangle(0, 0, 200, 200) //预计截取的区域
	draw.PopClipPath()
	draw.SetClipPath(clipPathName)
	draw.Rectangle(100, 100, 300, 300) //预计draw区域

	src := imagick.NewMagickWand()
	defer src.Destroy()

	src.NewImage(400, 400, p2)

	src.DrawImage(draw) //实际draw区域=预计draw区域和预计截取的区域相交的区域

	src.WriteImage("/home/chen/demo.png")
}
