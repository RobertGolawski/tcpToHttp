package request

import (
	"errors"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {

	line, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("0 bytes in the message")
	}

	l, err := parseRequestLine(line)
	if err != nil {
		return nil, err
	}

	ret := Request{
		RequestLine: *l,
	}

	return &ret, nil
}

func parseRequestLine(line []byte) (*RequestLine, error) {
	strLine := string(line)
	strArr := strings.Split(strLine, "\r\n")

	reqParse := strings.Split(strArr[0], " ")

	if len(reqParse) < 3 {
		return nil, errors.New("Incorrect number of request parts")
	}

	rLine := RequestLine{}

	method, err := verifyMethod(reqParse[0])
	if err != nil {
		log.Printf("error verifying method: %v for string: %s", err, reqParse[0])
		return nil, err
	}

	rLine.Method = method

	target, err := verifyTarget(reqParse[1])
	if err != nil {
		log.Printf("error verifying target: %v for string: %s", err, reqParse[1])
		return nil, err
	}

	rLine.RequestTarget = target

	version, err := verifyVersion(reqParse[2])
	if err != nil {
		log.Printf("error verifying version: %v for string: %s", err, reqParse[2])
		return nil, err
	}

	rLine.HttpVersion = version

	return &rLine, nil
}

func verifyMethod(m string) (string, error) {
	switch m {
	case "GET":
		return m, nil
	case "POST":
		return m, nil
	case "PUT":
		return m, nil
	case "DELETE":
		return m, nil
	default:
		return "", errors.New("Invalid HTTP method")
	}
}

func verifyTarget(t string) (string, error) {
	return t, nil
}

func verifyVersion(v string) (string, error) {
	vArr := strings.Split(v, "/")

	if len(vArr) != 2 {
		return "", errors.New("Incorrect version formatting")
	}

	if vArr[0] != "HTTP" {
		return "", errors.New("Incorrect protocol")
	}

	ver := vArr[1]
	switch ver {
	case "1.1":
		return ver, nil
	default:
		return "", errors.New("Unsupported version")
	}
}
