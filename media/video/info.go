package video

import (
	"regexp"
	"strconv"
	"strings"
)

type hook struct {
	matcher *regexp.Regexp
	handler func([][]byte)
}

func (p *FfmpegProcess) readInfo() error {
	var err error
	var line []byte

	hooks := []hook{
		hook{
			matcher: regexp.MustCompile(`Duration: (\d+:\d+:\d+\.\d+),`),
			handler: func(data [][]byte) {
				p.Info.Duration = string(data[1])
			},
		},
		hook{
			matcher: regexp.MustCompile(`Output #0, (\S+), to `),
			handler: func(data [][]byte) {
				p.Info.Format = "video/" + string(data[1])
			},
		},
		hook{
			matcher: regexp.MustCompile(`Stream #\d+:.*: Video:.*, (\d+x\d+).*, (\d+.?\d*) fps`),
			handler: func(data [][]byte) {
				resolution := strings.Split(string(data[1]), ":")
				p.Info.Width, err = strconv.ParseInt(resolution[0], 10, 64)
				if err != nil {
					return
				}
				p.Info.Height, err = strconv.ParseInt(resolution[1], 10, 64)
				if err != nil {
					return
				}
				p.Info.FPS, err = strconv.ParseFloat(string(data[2]), 64)
			},
		},
	}
	i := 0

	for i < len(hooks) {
		line, err = p.reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		matches := hooks[i].matcher.FindSubmatch(line)
		if len(matches) != 0 {
			hooks[i].handler(matches)
			if err != nil {
				return err
			}
			i++
		}
	}

	return nil
}
