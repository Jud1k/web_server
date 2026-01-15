package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

const bufferSize = 8

type parseState int 

const (
	stateInit parseState = iota
	stateDone
)

type Request struct{
	RequestLine RequestLine
	State parseState
}

type RequestLine struct{
	HttpVersion string
	RequestTarget string
	Method string
}

func (r *Request)parse(data []byte)(int,error){
	switch r.State {
	case stateInit:
		n,rl,err:=parseRequestLine(data)
		if err!=nil{
			return 0,err
		}
		if n==0{
			return 0,nil
		}
		r.RequestLine=*rl
		r.State=stateDone
		return n,nil
	case stateDone:
		return 0,nil
	default: 
		return 0,errors.New("error: unknown state")
	}
}

func RequestFromReader(reader io.Reader)(*Request, error){
	readToIndex := 0
	buffer := make([]byte,bufferSize)
	r := &Request{State: stateInit}
	for r.State!=stateDone{
		if readToIndex==len(buffer){
			doubleBuf(&buffer)
		}
		numBytesRead,err := reader.Read(buffer[readToIndex:])
		if err!=nil{
			if errors.Is(err,io.EOF){
				r.State=stateDone
				break
			}
			return nil,err
		}
		readToIndex+=numBytesRead
		// if readToIndex>=len(buffer){
		// 	doubleBuf(&buffer)
		// }
		numBytesParsed,err:=r.parse(buffer[:readToIndex])
		if err!=nil{
			return nil,err
		}
		if numBytesParsed==0{
			continue
		}
		copy(buffer,buffer[numBytesParsed:readToIndex])
		readToIndex-=numBytesParsed
	}
	return r,nil
}

func doubleBuf(buffer *[]byte){
	newSlice := make([]byte,len(*buffer)*2)
	copy(newSlice,*buffer)
	*buffer=newSlice
}

func parseRequestLine(data []byte)(int,*RequestLine, error){
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0,nil, nil 
	}
	line := data[:idx]
	consumed := idx+2
	partsLine := strings.Split(string(line)," ")
	if len(partsLine)!=3{
		return 0,nil,errors.New("error: Not valid format")
	}
	method := partsLine[0]
	if !isValidMethod(method){
		return 0,nil,errors.New("error: Invalid HTTP method")
	}	
	//TODO: Add validation to target
	target := partsLine[1]
	if !strings.HasPrefix(partsLine[2], "HTTP/") {
		return 0,nil, errors.New("error: invalid HTTP version format")
	}

	version := strings.TrimPrefix(partsLine[2], "HTTP/")
	if version != "1.1" {
		return 0,nil, errors.New("error: unsupported HTTP version")
	}	
	req := RequestLine{
			Method: method,
			RequestTarget: target,
			HttpVersion: version,
		}
	return consumed,&req,nil

}

func isValidMethod(str string)bool{
	switch str{
	case "GET","POST","PUT","PATCH","DELETE","OPTIONS","HEAD","TRACE","CONNECT": return true
	default: return false
	}
}