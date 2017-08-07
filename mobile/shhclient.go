// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Contains a wrapper for the Ethereum whisper client.

package geth

import (
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv5"
)

type ShhClient struct {
	client *shhclient.Client
}

func NewShhClient(rawurl string) (clieint *ShhClient, _ error) {
	rawClient, err := shhclient.Dial(rawurl)
	return &ShhClient{rawClient}, err
}

func (sc *ShhClient) Version(ctx *Context) (int64, error) {
	ver, err := sc.client.Version(ctx.context)
	return int64(ver), err
}

func (sc *ShhClient) Info(ctx *Context) (*Info, error) {
	rawInfo, err := sc.client.Info(ctx.context)
	return &Info{rawInfo}, err
}

func (sc *ShhClient) SetMaxMessageSize(ctx *Context, size int64) error {
	return sc.client.SetMaxMessageSize(ctx.context, uint32(size))
}

func (sc *ShhClient) MarkTrustedPeer(ctx *Context, enode string) error {
	return sc.client.MarkTrustedPeer(ctx.context, enode)
}

func (sc *ShhClient) NewKeyPair(ctx *Context) (string, error) {
	return sc.client.NewKeyPair(ctx.context)
}

func (sc *ShhClient) AddPrivateKey(ctx *Context, key []byte) (string, error) {
	return sc.client.AddPrivateKey(ctx.context, key)
}

func (sc *ShhClient) DeleteKeyPair(ctx *Context, id string) (string, error) {
	return sc.client.DeleteKeyPair(ctx.context, id)
}

func (sc *ShhClient) HasKeyPair(ctx *Context, id string) (bool, error) {
	return sc.client.HasKeyPair(ctx.context, id)
}

func (sc *ShhClient) PublicKey(ctx *Context, id string) ([]byte, error) {
	return sc.client.PublicKey(ctx.context, id)
}

func (sc *ShhClient) PrivateKey(ctx *Context, id string) ([]byte, error) {
	return sc.client.PrivateKey(ctx.context, id)
}

func (sc *ShhClient) NewSymmetricKey(ctx *Context) (string, error) {
	return sc.client.NewSymmetricKey(ctx.context)
}

func (sc *ShhClient) AddSymmetricKey(ctx *Context, key []byte) (string, error) {
	return sc.client.AddSymmetricKey(ctx.context, key)
}

func (sc *ShhClient) GenerateSymmetricKeyFromPassword(ctx *Context, passwd []byte) (string, error) {
	return sc.client.GenerateSymmetricKeyFromPassword(ctx.context, passwd)
}

func (sc *ShhClient) HasSymmetricKey(ctx *Context, id string) (bool, error) {
	return sc.client.HasSymmetricKey(ctx.context, id)
}

func (sc *ShhClient) GetSymmetricKey(ctx *Context, id string) ([]byte, error) {
	return sc.client.GetSymmetricKey(ctx.context, id)
}

func (sc *ShhClient) DeleteSymmetricKey(ctx *Context, id string) error {
	return sc.client.DeleteSymmetricKey(ctx.context, id)
}

func (sc *ShhClient) Post(ctx *Context, message *NewMessage) error {
	return sc.client.Post(ctx.context, message.newMessage)
}

type NewMessageHandler interface {
	OnNewMessage(message *Message)
	OnError(failure string)
}

func (sc *ShhClient) SubscribeMessages(ctx *Context, criteria *Criteria, handler NewMessageHandler, buffer int) (*Subscription, error) {
	ch := make(chan *whisper.Message, buffer)
	rawSub, err := sc.client.SubscribeMessages(ctx.context, criteria.criteria, ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case message := <-ch:
				handler.OnNewMessage(&Message{message})
			case err := <-rawSub.Err():
				handler.OnError(err.Error())
				return
			}
		}
	}()
	return &Subscription{rawSub}, nil
}

func (sc *ShhClient) NewMessageFilter(ctx *Context, criteria *Criteria) (string, error) {
	return sc.client.NewMessageFilter(ctx.context, criteria.criteria)
}

func (sc *ShhClient) DeleteMessageFilter(ctx *Context, id string) error {
	return sc.client.DeleteMessageFilter(ctx.context, id)
}

func (sc *ShhClient) FilterMessages(ctx *Context, id string) ([]*Message, error) {
	messages, err := sc.client.FilterMessages(ctx.context, id)
	if err != nil {
		return nil, err
	}

	res := make([]*Message, len(messages))
	for i, message := range messages {
		res[i] = &Message{message}
	}

	return res, nil
}
