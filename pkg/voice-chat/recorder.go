package chat

import (
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/faiface/beep"
	proto "github.com/gogo/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	voicechat "github.com/vectorhacker/voice-chat/pb"
)

const sampleRate = 44100

func Record(out *zmq.Socket) chan error {
	errChan := make(chan error)

	go func() {

		for {

			// select {
			// 	case <-time.After(17 * time.Millisecond):
			// 		rawSample, err := ioutil.TempFile("", "sound.wav")

			// 		if err != nil {
			// 			errChan <- err
			// 			return
			// 		}
			// 		err = wav.Encode(rawSample, stream, format)
			// 		if err != nil {
			// 			errChan <- err
			// 			return
			// 		}

			// 		sample, err := ioutil.ReadAll(rawSample)

			// 		voiceSample := &voicechat.VoiceSample{Sample: sample}

			// 		msg, err := proto.Marshal(voiceSample)
			// 		if err != nil {
			// 			errChan <- err
			// 			return
			// 		}

			// 		_, err = out.Send(string(msg), 0)
			// 		if err != nil {
			// 			errChan <- err
			// 			return
			// 		}
			// 	}

			rawSample, err := os.Open("./Lame_Drivers_-_01_-_Frozen_Egg.mp3")
			sample, err := ioutil.ReadAll(rawSample)

			voiceSample := &voicechat.VoiceSample{Sample: sample}

			msg, err := proto.Marshal(voiceSample)
			if err != nil {
				errChan <- err
				return
			}

			_, err = out.Send(string(msg), 0)
			if err != nil {
				errChan <- err
				return
			}

		}
	}()

	return errChan
}

func Noise() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			samples[i][0] = rand.Float64()*2 - 1
			samples[i][1] = rand.Float64()*2 - 1
		}
		return len(samples), true
	})
}
