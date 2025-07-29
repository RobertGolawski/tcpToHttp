package request

import (
	"errors"
	"io"
	"log"
	"strings"
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

type Request struct {
	RequestLine RequestLine
	State       int
}

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	ret := Request{State: 0}

	buf := make([]byte, bufferSize)
	readToIndex := 0

	for ret.State == 0 {
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			return nil, err
		}

		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		readToIndex += n

		i, err := ret.parse(buf)
		if err != nil {
			return nil, err
		}
		if i != 0 {
			buf = buf[n+1:]
		}

	}

	return &ret, nil
}

func (r *Request) parse(data []byte) (int, error) {
	rl, n, err := parseRequestLine(data)
	if n != 0 && err == nil {
		r.RequestLine = rl
		r.State = 1
	}

	return n, err
}

func parseRequestLine(line []byte) (RequestLine, int, error) {
	strLine := string(line)
	strArr := strings.Split(strLine, "\r\n")

	if !strings.Contains(strLine, crlf) {
		return RequestLine{}, 0, nil
	}

	reqParse := strings.Split(strArr[0], " ")

	if len(reqParse) != 3 {
		return RequestLine{}, 0, errors.New("incorrect number of request parts")
	}

	rLine := RequestLine{}

	method, err := verifyMethod(reqParse[0])
	if err != nil {
		log.Printf("error verifying method: %v for string: %s", err, reqParse[0])
		return RequestLine{}, 0, err
	}

	rLine.Method = method

	target, err := verifyTarget(reqParse[1])
	if err != nil {
		log.Printf("error verifying target: %v for string: %s", err, reqParse[1])
		return RequestLine{}, 0, err
	}

	rLine.RequestTarget = target

	version, err := verifyVersion(reqParse[2])
	if err != nil {
		log.Printf("error verifying version: %v for string: %s", err, reqParse[2])
		return RequestLine{}, 0, err
	}

	rLine.HTTPVersion = version

	return rLine, len(strArr[0]) + 2, nil
}

func verifyMethod(m string) (string, error) {
	// switch m {
	// case "GET", "POST", "PUT", "DELETE":
	// 	return m, nil
	// default:
	// 	return "", errors.New("invalid HTTP method")
	// }

	for _, c := range m {
		if c < 'A' || c > 'Z' {
			return "", errors.New("method contains non-uppercase character")
		}
	}

	return m, nil
}

func verifyTarget(t string) (string, error) {
	return t, nil
}

func verifyVersion(v string) (string, error) {
	vArr := strings.Split(v, "/")

	if len(vArr) != 2 {
		return "", errors.New("incorrect version formatting")
	}

	if vArr[0] != "HTTP" {
		return "", errors.New("incorrect protocol")
	}

	ver := vArr[1]
	switch ver {
	case "1.1":
		return ver, nil
	default:
		return "", errors.New("unsupported version")
	}
}
