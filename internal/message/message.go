package message

type Message struct {
	Header           Header
	QuestionsSection QuestionsSection
	AnswersSection   AnswersSection
}

func (m Message) Bytes() []byte {
	res := []byte{}
	res = append(res, m.Header.Bytes()...)
	res = append(res, m.QuestionsSection.Bytes()...)
	res = append(res, m.AnswersSection.Bytes()...)
	return res
}

func Parse(datagram []byte) Message {
	header := parseHeader(datagram[0:12])
	questionsSection, pointer := parseQuestionsSection(header.QDCOUNT, 12, datagram)
	answersSection := parseAnswersSection(header.ANCOUNT, pointer, datagram)

	return Message{
		Header:           header,
		QuestionsSection: questionsSection,
		AnswersSection:   answersSection,
	}
}
