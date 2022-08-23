package raycast

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"io"
	"log"
	"os"
)

type SoundPlayer struct {
	audioContext *audio.Context
	players      map[string]*audio.Player
}

const sampleRate = 44100

func NewSoundPlayer() *SoundPlayer {
	audioContext := audio.NewContext(sampleRate)
	return &SoundPlayer{
		audioContext: audioContext,
		players:      map[string]*audio.Player{},
	}
}

func (r *SoundPlayer) LoadSound(name string) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	var s audioStream
	b, err := os.ReadFile("res/sound/" + name + ".mp3")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	s, err = mp3.DecodeWithSampleRate(sampleRate, bytes.NewReader(b))
	if err != nil {
		log.Fatalln(err)
	}
	p, err := r.audioContext.NewPlayer(s)
	if err != nil {
		log.Fatalln(err)
	}
	r.players[name] = p
}

func (r SoundPlayer) PlaySound(name string) {
	p, ok := r.players[name]
	if !ok {
		fmt.Println("failed to get player for sound: ", name)
		return
	}
	p.Rewind() // we need to rewind the tape
	p.Play()
}
