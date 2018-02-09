// 可视化go并发
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ajstarks/svgo"
)

var width = 800
var height = 400
var startTime = time.Now().UnixNano()

func drawPoint(osvg *svg.SVG, pnt int, process int) {
	sec := time.Now().UnixNano()
	diff := (int64(sec) - int64(startTime)) / 10000000
	pointLocation, pointLocationV := int(diff), 0

	color := "#000000"
	switch {
	case process == 1:
		pointLocationV = 60
		color = "#cc6666"
	default:
		pointLocationV = 180
		color = "#66cc66"
	}
	if pointLocation > width {
		return
	}
	//fmt.Println(pointLocation, pointLocationV)
	osvg.Rect(pointLocation, pointLocationV, 3, 5, "fill:"+color+";stroke:none;")
	time.Sleep(150 * time.Millisecond)
}

func visualize(rw http.ResponseWriter, req *http.Request) {
	startTime = time.Now().UnixNano()
	fmt.Println("Request to /visualize")
	rw.Header().Set("Content-Type", "image/svg+xml")

	buf := bytes.NewBuffer(nil)
	outputSVG := svg.New(buf)
	// outputSVG := svg.New(rw) // 直接写rw可能会导致svg内容错乱
	outputSVG.Start(width, height)

	outputSVG.Rect(10, 10, 780, 100, "fill:#eeeeee;stroke:none")
	outputSVG.Text(20, 30, "Process 1 Timeline", "textanchor:start;font-size:12px;fill:#333333")
	outputSVG.Rect(10, 130, 780, 100, "fill:#eeeeee;stroke:none")
	outputSVG.Text(20, 150, "Process 2 Timeline", "textanchor:start;font-size:12px;fill:#333333")

	for i := 0; i < 801; i++ {
		timeText := strconv.FormatInt(int64(i), 10)
		if i%100 == 0 {
			outputSVG.Text(i, 380, timeText, "text-anchor:middle;fontsize:10px;fill:#000000")
		} else if i%4 == 0 {
			outputSVG.Circle(i, 377, 1, "fill:#cccccc;stroke:none")
		}

		if i%10 == 0 {
			outputSVG.Rect(i, 0, 1, 400, "fill:#dddddd")
		}
		if i%50 == 0 {
			outputSVG.Rect(i, 0, 1, 400, "fill:#cccccc")
		}
	}
	for i := 0; i < 100; i++ {
		go drawPoint(outputSVG, i, 1)
		drawPoint(outputSVG, i, 2)
	}
	outputSVG.Text(650, 360, "Run without goroutines", "textanchor:start;font-size:12px;fill:#333333")
	outputSVG.End()

	fmt.Println(buf.String())
	buf.WriteTo(rw)
}

func main() {
	http.Handle("/visualize", http.HandlerFunc(visualize))
	err := http.ListenAndServe(":1900", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
