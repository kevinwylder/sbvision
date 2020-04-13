package video

import (
	"strconv"
	"strings"
)

func generateFfmpegFilter(bits, width, height int) string {
	var filter strings.Builder
	pow := 1
	for i := 0; i < bits; i++ {
		for j := 0; j < 2; j++ {
			filter.WriteString("drawbox=enable='eq(mod(floor(n/")
			filter.WriteString(strconv.Itoa(pow))
			filter.WriteString("),2),")
			filter.WriteString(strconv.Itoa(j))
			filter.WriteString(")':x=")
			filter.WriteString(strconv.Itoa(i * width))
			filter.WriteString(":y=0:w=")
			filter.WriteString(strconv.Itoa(width))
			filter.WriteString(":h=")
			filter.WriteString(strconv.Itoa(height))
			filter.WriteString(":color=")
			if j == 0 {
				filter.WriteString("black")
			} else {
				filter.WriteString("white")
			}
			if j == 0 || i != bits-1 {
				filter.WriteRune(',')
			}
		}
		pow *= 2
	}
	return filter.String()
}
