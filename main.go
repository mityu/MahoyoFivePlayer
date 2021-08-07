package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/hajimehoshi/oto"
)

type Command_t int

const (
	Exit = iota
	Pause
	Resume
	Status
	Null
)

func generateHelp(source [][]string) string {
	maxlen := 0
	for i := 0; i < len(source); i++ {
		l := len(source[i][0])
		if l > maxlen {
			maxlen = l
		}
	}

	var help string
	formatter := fmt.Sprintf("%%-%ds    %%s\n", maxlen)
	for i := 0; i < len(source); i++ {
		help += fmt.Sprintf(formatter, source[i][0], source[i][1])
	}

	return help
}

const (
	FramePerSec      = 44100 // SampleRate
	BytesPerFrame    = 4
	OneFragmentBytes = FramePerSec * BytesPerFrame / 4
)

var HelpText = generateHelp([][]string{
	{"e[xit]|q[uit]", "Quit this player"},
	{"p[ause]", "Pause music (valid only when playing)"},
	{"r[esume]", "Cancel pausing (valid only when pausing)"},
	{"s[tatus]", "Show status"},
	{"h[elp]", "Show this help"},
})

//go:embed FivePre.dmp3
var PreMusic []byte

//go:embed FiveLoop.dmp3
var LoopMusic []byte

type Music struct {
	Content []byte
	Index   uint32
}

func (m *Music) length() uint32 {
	return uint32(len(m.Content))
}

func (m *Music) playFragment(p *oto.Player) error {
	var to_write []byte
	if (m.Index + OneFragmentBytes) >= m.length() {
		to_write = m.Content[m.Index:]
		m.Index = m.length()
	} else {
		to_write = m.Content[m.Index : m.Index+OneFragmentBytes]
		m.Index += OneFragmentBytes
	}
	_, err := p.Write(to_write)
	if err != nil {
		return err
	}
	return nil
}

func (m *Music) EOF() bool {
	return m.Index >= m.length()
}

func play(ch <-chan Command_t, done chan<- bool, reterr chan<- error) {
	c, err := oto.NewContext(FramePerSec, 2, 2, OneFragmentBytes)
	if err != nil {
		reterr <- err
		return
	}
	p := c.NewPlayer()
	defer p.Close()

	pauseNow := false

	processCommand := func(command Command_t) {
		switch command {
		case Exit:
			reterr <- nil
		case Pause:
			if pauseNow {
				fmt.Println("Already paused")
			} else {
				fmt.Println("Paused")
			}
			pauseNow = true
		case Resume:
			if pauseNow {
				fmt.Println("Resumed")
			} else {
				fmt.Println("Already resumed")
			}
			pauseNow = false
		case Status:
			if pauseNow {
				fmt.Println("Paused")
			} else {
				fmt.Println("Playing")
			}
		default:
			// Do-nothing
		}
		done <- true
	}

	music := &Music{Content: PreMusic, Index: 0}
	for !music.EOF() {
		if !pauseNow {
			music.playFragment(p)
		}

		select {
		case command := <-ch:
			processCommand(command)
		default:
			// Do-nothing
		}
	}

	music = &Music{Content: LoopMusic, Index: 0}
	for {
		music.Index = 0
		for !music.EOF() {
			if !pauseNow {
				music.playFragment(p)
			}

			select {
			case command := <-ch:
				processCommand(command)
			default:
				// Do-nothing
			}
		}
	}
	reterr <- nil
}

func run() error {
	ch := make(chan Command_t)
	done := make(chan bool)
	err := make(chan error)
	go play(ch, done, err)

	for {
		var input string
		fmt.Print("> ")
		fmt.Scanln(&input)

		var command Command_t = Null
		input = strings.ToLower(strings.Trim(input, " "))
		if strings.HasPrefix("exit", input) || strings.HasPrefix("quit", input) {
			command = Exit
			return nil
		} else if strings.HasPrefix("pause", input) {
			command = Pause
		} else if strings.HasPrefix("resume", input) {
			command = Resume
		} else if strings.HasPrefix("status", input) {
			command = Status
		} else if strings.HasPrefix("help", input) {
			fmt.Println(HelpText)
		}

		ch <- command
		<-done // Wait for processing
		select {
		case e := <-err:
			return e
		default:
			// Do-nothing
		}
	}
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
