package chat

import (
	proto "github.com/gogo/protobuf/proto"
	"github.com/mesilliac/pulse-simple"
	zmq "github.com/pebbe/zmq4"
	voicechat "github.com/vectorhacker/voice-chat/pb"
)

const sampleRate = 11025

func Record(out *zmq.Socket, name string) chan error {
	errChan := make(chan error)

	go func() {
		ss := pulse.SampleSpec{pulse.SAMPLE_FLOAT32LE, sampleRate, 1}

		capture, err := pulse.Capture("audio", "audio", &ss)
		if err != nil {
			errChan <- err
			return
		}
		defer capture.Free()

		for {
			buffer := make([]byte, sampleRate*4/25)

			_, err := capture.Read(buffer)
			if err != nil {
				errChan <- err
				return
			}

			voiceSample := &voicechat.VoiceSample{Sample: buffer, Speaker: name}
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
