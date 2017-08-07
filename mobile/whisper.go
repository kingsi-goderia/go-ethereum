package geth

import whisper "github.com/ethereum/go-ethereum/whisper/whisperv5"

type Info struct {
	info whisper.Info
}

func (i *Info) GetMemory() int64         { return int64(i.info.Memory) }
func (i *Info) GetMessages() int64       { return int64(i.info.Messages) }
func (i *Info) GetMinPow() float64       { return i.info.MinPow }
func (i *Info) GetMaxMessageSize() int64 { return int64(i.info.MaxMessageSize) }

type TopicType struct {
	topicType whisper.TopicType
}

func NewTopicType() *TopicType {
	return new(TopicType)
}

func (t *TopicType) SetTopic(topic []byte) { t.topicType = whisper.BytesToTopic(topic) }

type NewMessage struct {
	newMessage whisper.NewMessage
}

func NewNewMessage() *NewMessage {
	return new(NewMessage)
}

func (nm *NewMessage) GetSymKeyID() string   { return nm.newMessage.SymKeyID }
func (nm *NewMessage) GetPublicKey() []byte  { return nm.newMessage.PublicKey }
func (nm *NewMessage) GetSig() string        { return nm.newMessage.Sig }
func (nm *NewMessage) GetTTL() int64         { return int64(nm.newMessage.TTL) }
func (nm *NewMessage) GetTopic() TopicType   { return TopicType{topicType: nm.newMessage.Topic} }
func (nm *NewMessage) GetPayload() []byte    { return nm.newMessage.Payload }
func (nm *NewMessage) GetPadding() []byte    { return nm.newMessage.Padding }
func (nm *NewMessage) GetPowTime() int64     { return int64(nm.newMessage.PowTime) }
func (nm *NewMessage) GetPowTarget() float64 { return nm.newMessage.PowTarget }
func (nm *NewMessage) GetTargetPeer() string { return nm.newMessage.TargetPeer }

func (nm *NewMessage) SetSymKeyID(symKeyID string)     { nm.newMessage.SymKeyID = symKeyID }
func (nm *NewMessage) SetPublicKey(publicKey []byte)   { nm.newMessage.PublicKey = publicKey }
func (nm *NewMessage) SetSig(sig string)               { nm.newMessage.Sig = sig }
func (nm *NewMessage) SetTTL(ttl int64)                { nm.newMessage.TTL = uint32(ttl) }
func (nm *NewMessage) SetTopic(topic TopicType)        { nm.newMessage.Topic = topic.topicType }
func (nm *NewMessage) SetPayload(payload []byte)       { nm.newMessage.Payload = payload }
func (nm *NewMessage) SetPadding(padding []byte)       { nm.newMessage.Padding = padding }
func (nm *NewMessage) SetPowTime(powTime int64)        { nm.newMessage.PowTime = uint32(powTime) }
func (nm *NewMessage) SetPowTarget(powTarget float64)  { nm.newMessage.PowTarget = powTarget }
func (nm *NewMessage) SetTargetPeer(targetPeer string) { nm.newMessage.TargetPeer = targetPeer }

type Message struct {
	message *whisper.Message
}

func (m *Message) GetSig() []byte      { return m.message.Sig }
func (m *Message) GetTTL() int64       { return int64(m.message.TTL) }
func (m *Message) GetTimestamp() int64 { return int64(m.message.Timestamp) }
func (m *Message) GetTopic() TopicType { return TopicType{topicType: m.message.Topic} }
func (m *Message) GetPayload() []byte  { return m.message.Payload }
func (m *Message) GetPadding() []byte  { return m.message.Padding }
func (m *Message) GetPoW() float64     { return m.message.PoW }
func (m *Message) GetHash() []byte     { return m.message.Hash }
func (m *Message) GetDst() []byte      { return m.message.Dst }

type Messages struct {
	messages []*Message
}

type Criteria struct {
	criteria whisper.Criteria
}

func NewCriteria() *Criteria {
	return new(Criteria)
}

func (c *Criteria) GetSymKeyID() string     { return c.criteria.SymKeyID }
func (c *Criteria) GetPrivateKeyID() string { return c.criteria.PrivateKeyID }
func (c *Criteria) GetSig() []byte          { return c.criteria.Sig }
func (c *Criteria) GetMinPow() float64      { return c.criteria.MinPow }
func (c *Criteria) GetTopics() []TopicType {
	topics := make([]TopicType, len(c.criteria.Topics))
	for i, topic := range c.criteria.Topics {
		topics[i] = TopicType{topicType: topic}
	}
	return topics
}
func (c *Criteria) GetAllowP2P() bool { return c.criteria.AllowP2P }

func (c *Criteria) SetSymKeyID(symKeyID string)         { c.criteria.SymKeyID = symKeyID }
func (c *Criteria) SetPrivateKeyID(privateKeyID string) { c.criteria.PrivateKeyID = privateKeyID }
func (c *Criteria) SetSig(sig []byte)                   { c.criteria.Sig = sig }
func (c *Criteria) SetMinPow(minPow float64)            { c.criteria.MinPow = minPow }
func (c *Criteria) SetTopics(topics []TopicType) {
	topicTypes := make([]whisper.TopicType, len(topics))
	for i, topic := range topics {
		topicTypes[i] = topic.topicType
	}
	c.criteria.Topics = topicTypes
}
func (c *Criteria) SetAllowP2P(allowP2P bool) { c.criteria.AllowP2P = allowP2P }
