package message

import (
	"encoding/binary"
	"strings"
)

func parseQuestionsSection(count uint16, pointer int, data []byte) (QuestionsSection, int) {
	qs := QuestionsSection{}
	for i := 0; i < int(count); i++ {
		q, read := parseQuestion(pointer, data)
		qs.Questions = append(qs.Questions, q)
		pointer = read
	}
	return qs, pointer
}

func parseQuestion(read int, data []byte) (Question, int) {
	q := Question{}

	q.Name, read = parseQuestionName(read, data)

	q.Type = binary.BigEndian.Uint16(data[read:])
	read = read + 2
	q.Class = binary.BigEndian.Uint16(data[read:])
	read = read + 2

	return q, read
}

func parseQuestionName(pointer int, data []byte) (string, int) {
	labels := []string{}
	for {
		first := data[pointer]
		pointer++
		// if null byte - end
		if first == 0 {
			break
		}
		// if redirect - read labels from there
		if (first & 0b11000000) == 0b11000000 {
			second := data[pointer]
			pointer++
			redirectAddress := binary.BigEndian.Uint16([]byte{first & 0b00111111, second})
			question, _ := parseQuestionName(int(redirectAddress), data)
			labels = append(labels, question)
			break
		}
		// else read from here
		label := data[pointer : pointer+int(first)]
		pointer = pointer + int(first)
		labels = append(labels, string(label))
	}

	return strings.Join(labels, "."), pointer
}

type QuestionsSection struct {
	Questions []Question
}

func (qs QuestionsSection) Bytes() []byte {
	res := []byte{}
	for _, q := range qs.Questions {
		res = append(res, q.Bytes()...)
	}
	return res
}

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

func (q Question) Bytes() []byte {
	labelsBytes := []byte{}
	for label := range strings.SplitSeq(q.Name, ".") {
		labelsBytes = append(labelsBytes, byte(len(label)))
		labelsBytes = append(labelsBytes, []byte(label)...)
	}
	labelsBytes = append(labelsBytes, byte(0))
	// QNAME + QTYPE (2 bytes) + QCLASS (2 bytes)
	res := make([]byte, len(labelsBytes)+4)
	copy(res, labelsBytes)
	binary.BigEndian.PutUint16(res[len(labelsBytes):], q.Type)
	binary.BigEndian.PutUint16(res[len(labelsBytes)+2:], q.Class)
	return res
}
