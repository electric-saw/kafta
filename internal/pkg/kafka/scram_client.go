package kafka

import (
	"github.com/xdg/scram"
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) error {
	var err error
	x.Client, err = x.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (string, error) {
	response, err := x.ClientConversation.Step(challenge)
	return response, err
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
