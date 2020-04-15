package media_test

import (
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/kevinwylder/sbvision"

	"github.com/kevinwylder/sbvision/media"
)

func TestRanger(t *testing.T) {
	file := make([]byte, 1024*1024*10)
	for i := 0; i < len(file); i++ {
		file[i] = byte(i % 256)
	}

	location, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(location)
	}()

	assets, err := media.NewAssetDirectory(location)
	if err != nil {
		t.Fatal(err)
	}

	fake := sbvision.Video{
		Type:     sbvision.YoutubeVideo,
		ShareURL: "thisisatest",
	}
	assets.VideoPath(&fake)

	ioutil.WriteFile(assets.VideoPath(&fake), file, 0777)

	ranger := assets.GetVideo(&fake)

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go testRanger(ranger, t, &wg)
	}
	wg.Wait()

}

func testRanger(ranger *media.Ranger, t *testing.T, wg *sync.WaitGroup) {
	start := rand.Int63() % (1024 * 1024 * 10)
	var end int64
	if rand.Int()%2 == 0 {
		end = start + rand.Int63()%(1024*1024)
	} else {
		end = 0
	}
	data, err := ranger.GetRange(start, end)
	if err != nil {
		t.Fail()
		t.Log(err)
	}
	for i, byte := range data {
		if (i+int(start))%256 != int(byte) {
			t.Fail()
			t.Log("Wrong byte")
		}
	}
	wg.Done()
}
