package hardware

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type SoundCard struct {
	samples         [4]rl.Sound
	prevSoundParams Tone
	currSound       rl.Sound
	isPlaying       bool
}

type Tone struct {
	Wave   Nybble
	Octave Nybble
	Pitch  Nybble
	Volume Nybble
}

const (
	OCTAVE_2 = 0
	OCTAVE_3 = 1
	OCTAVE_4 = 2
	OCTAVE_5 = 3
)

const (
	WAVE_SQUARE   = 0
	WAVE_TRIANGLE = 1
	WAVE_SAWTOOTH = 2
	WAVE_NOISE    = 3
)

const SEMITONE_INTERVAL = 1.059463

func NewSoundCard(masterVol float32, sampleFiles [4]string) *SoundCard {
	rl.InitAudioDevice()
	rl.SetMasterVolume(masterVol)
	return &SoundCard{
		samples: [4]rl.Sound{
			WAVE_SQUARE:   rl.LoadSound(sampleFiles[WAVE_SQUARE]),
			WAVE_TRIANGLE: rl.LoadSound(sampleFiles[WAVE_TRIANGLE]),
			WAVE_SAWTOOTH: rl.LoadSound(sampleFiles[WAVE_SAWTOOTH]),
			WAVE_NOISE:    rl.LoadSound(sampleFiles[WAVE_NOISE]),
		},
		isPlaying: false,
	}
}

// Plays a beep with a given pitch, octave, waveform and volume
func (sc *SoundCard) PlayTone(t Tone) {
	t.Volume %= 16
	if sc.prevSoundParams != t {
		if t.Volume > 0 {
			rl.SetSoundVolume(sc.samples[t.Wave], (1/float32(0xF))*float32(t.Volume))
			pitchMul := float32(math.Pow(SEMITONE_INTERVAL, float64(t.Pitch))) * float32(math.Pow(2, float64(t.Octave)))
			rl.SetSoundPitch(sc.samples[t.Wave], pitchMul)
			rl.StopSound(sc.currSound)
			sc.currSound = sc.samples[t.Wave]
			sc.isPlaying = true
		} else {
			sc.isPlaying = false
		}

		sc.prevSoundParams = t
	}
}

func (sc *SoundCard) StopAll() {
	sc.isPlaying = false
	for _, s := range sc.samples {
		rl.StopSound(s)
	}
}

func (sc *SoundCard) Tick(vm *Machine) {
	if sc.isPlaying {
		if !rl.IsSoundPlaying(sc.currSound) {
			rl.PlaySound(sc.currSound)
		}
	} else {
		rl.StopSound(sc.currSound)
	}
}

func (sc *SoundCard) GetListener(vm *Machine) ([]byte, RAMListener) {
	return []byte{
			(PERIPHERAL_PAGE << 4) | FPG_SND_OPT,
			(PERIPHERAL_PAGE << 4) | FPG_SND_PTC,
			(PERIPHERAL_PAGE << 4) | FPG_SND_VOL,
		},
		func(val Nybble) {

			opt := vm.RAM[PERIPHERAL_PAGE][FPG_SND_OPT]
			ptc := vm.RAM[PERIPHERAL_PAGE][FPG_SND_PTC]
			vol := vm.RAM[PERIPHERAL_PAGE][FPG_SND_VOL]
			tone := Tone{
				Wave:   (opt >> 2) % 4,
				Octave: opt % 4,
				Pitch:  ptc,
				Volume: vol,
			}

			sc.PlayTone(tone)

		}
}

func (sc *SoundCard) Reset() {
	sc.StopAll()
}

func (sc *SoundCard) Terminate() {
	for _, s := range sc.samples {
		rl.UnloadSound(s)
	}
	rl.CloseAudioDevice()
}
