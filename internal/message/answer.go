package message

import (
	"encoding/binary"
	"strconv"
	"strings"
)

func parseAnswersSection(count uint16, pointer int, data []byte) AnswersSection {
	qs := AnswersSection{}
	for i := 0; i < int(count); i++ {
		q, read := parseAnswer(pointer, data)
		qs.Answers = append(qs.Answers, q)
		pointer = read
	}
	return qs
}

func parseAnswer(read int, data []byte) (Answer, int) {
	q := Answer{}

	q.Name, read = parseAnswerName(read, data)

	q.Type = binary.BigEndian.Uint16(data[read:])
	read = read + 2
	q.Class = binary.BigEndian.Uint16(data[read:])
	read = read + 2
	q.TTL = binary.BigEndian.Uint32(data[read:])
	read = read + 4
	q.Length = binary.BigEndian.Uint16(data[read:])
	read = read + 2
	datastr := make([]string, int(q.Length))
	for i, b := range data[read : read+int(q.Length)] {
		datastr[i] = strconv.Itoa(int(b))
	}
	q.Data = strings.Join(datastr, ".")
	read = read + int(q.Length)

	return q, read
}

func parseAnswerName(pointer int, data []byte) (string, int) {
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
			question, _ := parseAnswerName(int(redirectAddress), data)
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

type AnswersSection struct {
	Answers []Answer
}

func (qs AnswersSection) Bytes() []byte {
	res := []byte{}
	for _, a := range qs.Answers {
		res = append(res, a.Bytes()...)
	}
	return res
}

type Answer struct {
	Name   string
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   string
}

func (q Answer) Bytes() []byte {
	labelsBytes := []byte{}
	for label := range strings.SplitSeq(q.Name, ".") {
		labelsBytes = append(labelsBytes, byte(len(label)))
		labelsBytes = append(labelsBytes, []byte(label)...)
	}
	labelsBytes = append(labelsBytes, byte(0))
	// ANAME + TYPE (2 bytes) + CLASS (2 bytes) + TTL (4 bytes) + RDLENGTH (2 bytes)
	res := make([]byte, len(labelsBytes)+10)
	copy(res, labelsBytes)
	binary.BigEndian.PutUint16(res[len(labelsBytes):], q.Type)
	binary.BigEndian.PutUint16(res[len(labelsBytes)+2:], q.Class)
	binary.BigEndian.PutUint32(res[len(labelsBytes)+4:], q.TTL)
	binary.BigEndian.PutUint16(res[len(labelsBytes)+8:], q.Length)

	for ns := range strings.SplitSeq(q.Data, ".") {
		n, _ := strconv.Atoi(ns)
		res = append(res, byte(n))
	}

	return res
}
