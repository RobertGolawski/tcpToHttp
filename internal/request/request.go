package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	// "log"
	"strings"
	"tcpToHttp/internal/headers"
)

const (
	crlf       = "\r\n"
	bufferSize = 8
)

const (
	requestStateInitialised int = iota
	requestStateParsingHeaders
	requestStateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	State       int
}

type RequestLine struct {
	HTTPVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	ret := Request{State: requestStateInitialised, Headers: headers.NewHeaders()}

	buf := make([]byte, bufferSize)
	readToIndex := 0

	for ret.State != requestStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if ret.State != requestStateDone {
					return nil, errors.New("incomplete request")
				}
				break
			}
			return nil, err
		}

		readToIndex += n

		i, err := ret.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if i > 0 {
			copy(buf, buf[i:])
			readToIndex -= i
		}

	}
	return &ret, nil
}

func (r *Request) parse(line []byte) (int, error) {
	totalBytesParsed := 0

	for r.State != requestStateDone {
		n, err := r.parseSingle(line[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		totalBytesParsed += n
		if n == 0 {
			break
		}
	}

	return totalBytesParsed, nil

}

func (r *Request) parseSingle(line []byte) (int, error) {
	switch r.State {
	case requestStateInitialised:
		reqLine, n, err := parseRequestLine(line)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *reqLine
		r.State = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(line)
		if err != nil {
			return 0, err
		}

		if done {
			r.State = requestStateDone
		}

		return n, nil

	case requestStateDone:
		return 0, fmt.Errorf("error: trying to parse data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func parseRequestLine(line []byte) (*RequestLine, int, error) {
	idx := bytes.Index(line, []byte(crlf))
	if idx == -1 {

		return nil, 0, nil
	}

	strLine := string(line[:idx])

	reqParse := strings.Split(strLine, " ")
	if len(reqParse) != 3 {
		return nil, 0, errors.New("incorrect number of request parts")
	}

	rLine := &RequestLine{}

	method, err := verifyMethod(reqParse[0])
	if err != nil {
		return nil, 0, fmt.Errorf("error verifying method: %w for string: %s", err, reqParse[0])
	}

	rLine.Method = method

	target, err := verifyTarget(reqParse[1])
	if err != nil {
		return nil, 0, fmt.Errorf("error verifying target: %w for string: %s", err, reqParse[1])
	}

	rLine.RequestTarget = target

	version, err := verifyVersion(reqParse[2])
	if err != nil {
		return nil, 0, fmt.Errorf("error verifying version: %w for string: %s", err, reqParse[2])
	}

	rLine.HTTPVersion = version

	return rLine, idx + len(crlf), nil
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
	//placeholder function
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
