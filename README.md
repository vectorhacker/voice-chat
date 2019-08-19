Voice Chat
====

Voice Chat is a voice chat client and server that aims to be as simple as possible at hosting voice conversations.

## Goals:
- Be CLI based
- Be as simple as possible to use
- Provide a CLI based GUI (much like htop or vtop and the like)
- Use as little bandwidth as needed to transmit sound

## How to build
Make sure that the zeromq and portaudio libraries are installed

```shell
sudo apt install libzmq3-dev portaudio19-dev
```

Then go get the binary

```
go get github.com/vectorhacker/voice-chat/cmd/voice-chat # for the client cli
go get github.com/vectorhacker/voice-chat/cmd/server # to host your own server
```
