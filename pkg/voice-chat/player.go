package chat

import (
	"github.com/mesilliac/pulse-simple"
	"github.com/pkg/errors"
	voicechat "github.com/vectorhacker/voice-chat/pb"

	proto "github.com/gogo/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
)

func Play(in *zmq.Socket, name string) chan error {

	errChan := make(chan error)

	go func() {

		ss := pulse.SampleSpec{pulse.SAMPLE_FLOAT32LE, sampleRate, 1}
		playback, err := pulse.Playback("pulse-simple test", "playback test", &ss)

		if err != nil {
			errChan <- err
			return
		}
		defer playback.Free()
		defer playback.Drain()

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

			if sample.Speaker == name {
				continue
			}

			playback.Write(sample.Sample)
		}
	}()

	return errChan
}
