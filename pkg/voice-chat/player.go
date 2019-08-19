package chat

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

	"github.com/pkg/errors"
	voicechat "github.com/vectorhacker/voice-chat/pb"

	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	proto "github.com/gogo/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

func Play(in *zmq.Socket) chan error {

	errChan := make(chan error)

	go func() {
		for {
			rawSample, err := in.Recv(0)
			if err != nil {
				err = errors.WithStack(err)
				errChan <- errors.Wrap(err, "cannot receive")
				return
			}
			if rawSample == "" {
				continue
			}

			sample := &voicechat.VoiceSample{}
			err = proto.Unmarshal([]byte(rawSample), sample)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot unmarshal")
				return
			}

			buf := bytes.NewBuffer(sample.Sample)

			sound := ioutil.NopCloser(buf)
			streamer, format, err := wav.Decode(sound)
			if err != nil {
				errChan <- errors.Wrap(err, "cannot decode")
				return
			}
			defer streamer.Close()

			log.Println("playing...")

			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

			speaker.Play(streamer)
			<-time.After(17 * time.Millisecond)
		}
	}()

	return errChan
}
