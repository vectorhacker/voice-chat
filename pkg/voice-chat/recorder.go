package chat

import (
	"time"

	"github.com/MarkKremer/microphone"
	"github.com/faiface/beep/wav"
	proto "github.com/gogo/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	"github.com/pkg/errors"
	voicechat "github.com/vectorhacker/voice-chat/pb"
)

const sampleRate = 44100

func Record(out *zmq.Socket) chan error {
	errChan := make(chan error)

	go func() {
		err := microphone.Init()
		if err != nil {
			errChan <- err
			return
		}
		defer microphone.Terminate()

		streamer, format, err := microphone.OpenDefaultStream(sampleRate)
		if err != nil {
			errChan <- errors.Wrap(err, "unable to create stream")
		}
		defer streamer.Close()

		for {
			streamer.Start()

			<-time.After(1 * time.Second)

			streamer.Stop()

			f := &writer{}

			wav.Encode(f, streamer, format)

			sample := f.Bytes()
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
