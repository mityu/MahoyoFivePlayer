package main

import (
	"io"
    "fmt"
	"os"

	"github.com/hajimehoshi/go-mp3"
)

func copyBytes(o []byte, d io.Reader) error {
	l := len(o)
	i := 0

	fmt.Printf("Start to copy %d bytes...\n", l)
	for {
		n, err := d.Read(o[i :])
		i += n
		if err != nil {
			return err
		}
		if i >= l {
			fmt.Printf("Expected: %d bytes, Copied %d bytes\n", l, i)
			return nil
		}
	}
}

const (
	FramePerSec = 44100
	BytesPerFrame = 4
)

func transport(audio_file string) {
	f, err := os.Open(audio_file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

    mp3, err := mp3.NewDecoder(f)
    if err != nil {
        fmt.Println(err)
        return
    }


	pre := make([]byte, FramePerSec / 100 * BytesPerFrame * (13200 + 45))
	loop := make([]byte, FramePerSec / 100 * BytesPerFrame * (10300 + 40))

	copyBytes(pre, mp3)
	copyBytes(loop, mp3)

	pre_file, err := os.Create("FivePre.dmp3")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pre_file.Close()
	pre_file.Write(pre)

	loop_file, err := os.Create("FiveLoop.dmp3")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer loop_file.Close()
	loop_file.Write(loop)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("%s: Just one argument must be given.", os.Args[0])
		return
	}
	transport(os.Args[1])
}
