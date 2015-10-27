package main

import(
	"github.com/gographics/imagick/imagick"
)

/*
convert -size 70x70 canvas:none -fill black -draw "circle 35,35 10,30" black_circle.png
convert -size 70x70 canvas:none -draw "circle 35,35 35,24" -negate -channel A -gaussian-blur 0x9 white_highlight.png
composite -compose atop -geometry -13-17 white_highlight.png black_circle.png black_ball.png
*/

func main(){
	//画布
	bg := imagick.NewMagickWand()
	p_bg := imagick.NewPixelWand()
	p_bg.SetColor("none")
	bg.NewImage(70, 70, p_bg)
	bg.SetImageFormat("png")
	//bg.WriteImage("/home/chen/bg.png")
	
	//画圆
	draw_c := imagick.NewDrawingWand()
	p2 := imagick.NewPixelWand()
	p2.SetColor("#000")
	draw_c.SetFillColor(p2)
	draw_c.Circle(35,35,60,30)
	bg.DrawImage(draw_c)
	//bg.WriteImage("/home/chen/bg1.png")
	
	//高光(highlight)覆盖层
	hl:=imagick.NewMagickWand()
	p_hl := imagick.NewPixelWand()
	p_hl.SetColor("none")
	hl.NewImage(70, 70, p_hl)
	hl.SetImageFormat("png")
	
	//高光区域
	draw_hl:= imagick.NewDrawingWand()
	draw_hl.Circle(35,35,35,24)
	hl.DrawImage(draw_hl)
	//hl.WriteImage("/home/chen/bg2.png")
	
	//反色
	hl.NegateImage(true)
	//hl.WriteImage("/home/chen/bg3.png")
	
	//hl.GaussianBlurImage(0,9)
	hl.GaussianBlurImageChannel(imagick.CHANNEL_ALPHA,0,9)
	//hl.WriteImage("/home/chen/bg4.png")   

	bg.CompositeImage(hl,imagick.COMPOSITE_OP_ATOP,-13,-17)
		
	bg.WriteImage("/home/chen/bg5.png")
}
