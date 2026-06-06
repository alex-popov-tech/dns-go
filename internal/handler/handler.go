package handler

import (
	"net"

	"github.com/codecrafters-io/dns-server-starter-go/internal/config"
	"github.com/codecrafters-io/dns-server-starter-go/internal/message"
)

func Handler(c config.Config, datagram []byte) []byte {
	input := message.Parse(datagram)

	if c.Resolver.IsValid() {
		return proxiedResponse(c, input)
	} else {
		return ownResponse(input)
	}
}

func ownResponse(input message.Message) []byte {
	output := message.Message{
		Header:           message.Header{},
		QuestionsSection: input.QuestionsSection,
		AnswersSection:   message.AnswersSection{},
	}

	output.Header.ID = input.Header.ID
	var rcode uint16
	if input.Header.Opcode() == 0 {
		rcode = 0
	} else {
		rcode = 4
	}
	output.Header.SetFlags(1, input.Header.Opcode(), 0, 0, input.Header.Rd(), 0, 0, rcode)
	output.Header.QDCOUNT = input.Header.QDCOUNT
	output.Header.ANCOUNT = input.Header.QDCOUNT
	output.Header.NSCOUNT = 0
	output.Header.ARCOUNT = 0

	for _, question := range input.QuestionsSection.Questions {
		output.AnswersSection.Answers = append(output.AnswersSection.Answers, message.Answer{
			Name:   question.Name,
			Type:   1,
			Class:  1,
			TTL:    60,
			Length: 4,
			Data:   "8.8.8.8",
		})
	}

	return output.Bytes()
}

func proxiedResponse(c config.Config, clientInput message.Message) []byte {
	addr := net.UDPAddrFromAddrPort(c.Resolver)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	answersSection := message.AnswersSection{}
	buf := make([]byte, 512)
	for i := range len(clientInput.QuestionsSection.Questions) {
		msg := message.Message{
			Header: clientInput.Header,
			QuestionsSection: message.QuestionsSection{
				Questions: []message.Question{clientInput.QuestionsSection.Questions[i]},
			},
		}
		msg.Header.QDCOUNT = 1

		_, _ = conn.Write(msg.Bytes())
		size, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		output := message.Parse(buf[:size])
		answersSection.Answers = append(answersSection.Answers, output.AnswersSection.Answers...)
	}

	res := message.Message{
		Header:           clientInput.Header,
		QuestionsSection: clientInput.QuestionsSection,
		AnswersSection:   answersSection,
	}
	res.Header.ID = clientInput.Header.ID
	var rcode uint16
	if clientInput.Header.Opcode() == 0 {
		rcode = 0
	} else {
		rcode = 4
	}
	res.Header.SetFlags(1, clientInput.Header.Opcode(), 0, 0, clientInput.Header.Rd(), 0, 0, rcode)
	res.Header.QDCOUNT = clientInput.Header.QDCOUNT
	res.Header.ANCOUNT = clientInput.Header.QDCOUNT
	res.Header.NSCOUNT = 0
	res.Header.ARCOUNT = 0

	return res.Bytes()
}
