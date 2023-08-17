package serializer

import (
	"encoding/base64"
	"encoding/binary"
	"git.multiverse.io/eventkit/kit/common/bytebuf"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/common/protocol"
	"git.multiverse.io/eventkit/kit/constant"
)

// LengthOfErrorCode default the fix length of error code
var LengthOfErrorCode = 10

var byteZero = byte(0)
var byteOne = byte(1)
var byteTwo = byte(2)

var zeroErrorLength = []byte{0, 0, 0, 0}

var byteZeroArray = []byte{byteZero}
var byteOneArray = []byte{byteOne}
var byteTwoArray = []byte{byteTwo}

func bool2bytes(b bool) []byte {
	if b {
		return byteOneArray
	}

	return byteZeroArray
}

func byte2bool(b byte) bool {
	if byteZero == b {
		return false
	}

	return true
}

func map2Bytes(m map[string]string) []byte {
	build := bytebuf.NewPointer()
	lenBytes := []byte{0, 0}
	for k, v := range m {
		keyBytes := []byte(k)
		valueBytes := []byte(v)
		lenKey := len(keyBytes)
		lenKeyValue := lenKey + len(valueBytes) + 1

		lenBytes[1] = byte(lenKeyValue)
		lenBytes[0] = byte(lenKeyValue >> 8)

		build.Write(lenBytes)         // 2 bytes:total length
		build.WriteByte(byte(lenKey)) // 1 byte:key length
		build.Write(keyBytes)         // bytes of key
		build.Write(valueBytes)       // bytes of value
	}
	len := build.Len()

	totalLenBytes := []byte{0, 0, 0, 0}
	totalLenBytes[3] = byte(len)
	totalLenBytes[2] = byte(len >> 8)
	totalLenBytes[1] = byte(len >> 16)
	totalLenBytes[0] = byte(len >> 24)

	result := append(totalLenBytes, build.Bytes()...)
	return result
}

func bytes2Map(inBytes []byte, totalLength uint32) map[string]string {
	var retMap map[string]string
	var tmpLen uint32

	retMap = make(map[string]string)
	for tmpLen = 0; tmpLen < totalLength; {
		keyLengthStartPos := tmpLen + 2
		keyLengthEndPos := tmpLen + 3
		pairLen := uint32(binary.BigEndian.Uint16([]byte(inBytes[tmpLen:keyLengthStartPos])))
		keyLength := uint32([]byte(inBytes[keyLengthStartPos:keyLengthEndPos])[0])
		keyEndPos := keyLengthEndPos + keyLength
		tmpLen = keyLengthStartPos + pairLen

		key := string(inBytes[keyLengthEndPos:keyEndPos])
		value := string(inBytes[keyEndPos:tmpLen])

		retMap[key] = value
	}

	return retMap
}

// String2ID parses the input string into uint64
func String2ID(inputStr string) (uint64, error) {
	inputBytes, _ := base64.StdEncoding.DecodeString(inputStr)

	if len(inputBytes) < 8 {
		err := errors.Errorf(constant.SystemInternalError, "inputstr invalid:[%s]", string(inputBytes))
		return 0, err
	}

	return binary.BigEndian.Uint64(inputBytes), nil
}

// ID2String serializes uint64 into bytes and base64 return the final result
func ID2String(id uint64) string {
	// id
	idBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	idBytes[7] = byte(id)
	idBytes[6] = byte(id >> 8)
	idBytes[5] = byte(id >> 16)
	idBytes[4] = byte(id >> 24)
	idBytes[3] = byte(id >> 32)
	idBytes[2] = byte(id >> 40)
	idBytes[1] = byte(id >> 48)
	idBytes[0] = byte(id >> 56)

	return base64.StdEncoding.EncodeToString(idBytes)
}

// Error2StringCompatible encodes the error into string
func Error2StringCompatible(err error, enableErrorCode bool) string {
	return base64.StdEncoding.EncodeToString(Error2BytesCompatible(err, enableErrorCode))
}

// String2Error decodes string into error
func String2Error(inputStr string) error {
	inputBytes, err := base64.StdEncoding.DecodeString(inputStr)
	if nil != err {
		return err
	}

	return Bytes2Error(inputBytes)
}

// Bytes2Error decodes bytes into error
func Bytes2Error(inputBytes []byte) error {
	if len(inputBytes) < 1 {
		err := errors.Errorf(constant.SystemInternalError, "inputstr invalid:[%s]", string(inputBytes))
		return err
	}
	errFlag := inputBytes[0]
	if byteOne == errFlag {
		errLength := binary.BigEndian.Uint32(inputBytes[1:5])
		errorMessage := string(inputBytes[5 : 5+errLength])

		return errors.Errorf(constant.SystemInternalError, errorMessage)
	} else if byteTwo == errFlag {
		errLength := binary.BigEndian.Uint32(inputBytes[1:5])
		errorCode := string(inputBytes[5 : 5+LengthOfErrorCode])
		errorMessage := string(inputBytes[5+LengthOfErrorCode : 5+errLength])

		return errors.Errorf(errorCode, errorMessage)
	} else {
		return nil
	}
}

// ProtoMsg2String encodes protocol.ProtoMessage or error into string
func ProtoMsg2String(msg *protocol.ProtoMessage, err error) string {
	return base64.StdEncoding.EncodeToString(ProtoMsg2BytesCompatible(msg, err, true))
}

// ProtoMsg2StringCompatible encodes protocol.ProtoMessage or error into string
func ProtoMsg2StringCompatible(msg *protocol.ProtoMessage, err error, enableErrorCode bool) string {
	return base64.StdEncoding.EncodeToString(ProtoMsg2BytesCompatible(msg, err, enableErrorCode))
}

// String2protoMsg decodes the input string into protocol.ProtoMessage
func String2protoMsg(inputStr string) (protocol.ProtoMessage, error) {
	msg := &protocol.ProtoMessage{}
	inputBytes, err := base64.StdEncoding.DecodeString(inputStr)
	if nil != err {
		return *msg, err
	}
	return Bytes2protoMsg(inputBytes)
}

var lpadSpaces = [...]string{
	"",                     // 0
	" ",                    // 1
	"  ",                   // 2
	"   ",                  // 3
	"    ",                 // 4
	"     ",                // 5
	"      ",               // 6
	"       ",              // 7
	"        ",             // 8
	"         ",            // 9
	"          ",           // 10
	"           ",          // 11
	"            ",         // 12
	"             ",        // 13
	"              ",       // 14
	"               ",      // 15
	"                ",     // 16
	"                 ",    // 17
	"                  ",   // 18
	"                   ",  // 19
	"                    ", // 20
}

func fixLen(targetLength int, src string) string {
	srcLen := len(src)
	if srcLen > targetLength {
		return src[0:targetLength]
	}

	return lpadSpaces[targetLength-srcLen] + src
}

// Error2BytesCompatible encodes the error into bytes
func Error2BytesCompatible(err error, enableErrorCode bool) []byte {
	build := bytebuf.NewPointer()
	if nil == err {
		build.Write(byteZeroArray)
		build.Write(zeroErrorLength)
	} else {
		if enableErrorCode {
			build.Write(byteTwoArray)
		} else {
			// for old protocol version can decode the error response
			build.Write(byteOneArray)
		}
		errorCode := errors.GetErrorCode(err)
		if "" == errorCode {
			errorCode = constant.SystemInternalError
		}
		errorCode = fixLen(LengthOfErrorCode, errorCode)
		errStr := errors.ErrorToString(err)
		errBytes := []byte(errorCode + errStr)
		lengthOfErr := len(errBytes)

		lengthOfErrBytes := []byte{0, 0, 0, 0}
		lengthOfErrBytes[3] = byte(lengthOfErr)
		lengthOfErrBytes[2] = byte(lengthOfErr >> 8)
		lengthOfErrBytes[1] = byte(lengthOfErr >> 16)
		lengthOfErrBytes[0] = byte(lengthOfErr >> 24)

		build.Write(lengthOfErrBytes)
		build.Write(errBytes)
	}
	return build.Bytes()
}

// ProtoMsg2Bytes encodes protocol.ProtoMessage or error to bytes
func ProtoMsg2Bytes(msg *protocol.ProtoMessage, err error) []byte {
	return ProtoMsg2BytesCompatible(msg, err, true)
}

// ProtoMsg2BytesCompatible encodes protocol.ProtoMessage or build-in error to bytes
func ProtoMsg2BytesCompatible(msg *protocol.ProtoMessage, err error, enableErrorCode bool) []byte {
	//var build strings.Builder
	build := bytebuf.NewPointer()
	if nil == err {
		build.Write(byteZeroArray)
		build.Write(zeroErrorLength)
	} else {
		return Error2BytesCompatible(err, enableErrorCode)
	}

	// id
	idBytes := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	idBytes[7] = byte(msg.ID)
	idBytes[6] = byte(msg.ID >> 8)
	idBytes[5] = byte(msg.ID >> 16)
	idBytes[4] = byte(msg.ID >> 24)
	idBytes[3] = byte(msg.ID >> 32)
	idBytes[2] = byte(msg.ID >> 40)
	idBytes[1] = byte(msg.ID >> 48)
	idBytes[0] = byte(msg.ID >> 56)
	//binary.BigEndian.PutUint64(idBytes, msg.ID)
	build.Write(idBytes)

	// need reply:0/1
	build.Write(bool2bytes(msg.NeedReply))
	// need ack:0/1
	build.Write(bool2bytes(msg.NeedAck))

	// length of session name(2 bytes)
	sessionNameLengthByte := []byte{0, 0}
	sessionNameBytes := []byte(msg.SessionName)
	lengthOfSessionName := len(sessionNameBytes)
	sessionNameLengthByte[1] = byte(lengthOfSessionName)
	sessionNameLengthByte[0] = byte(lengthOfSessionName >> 8)

	build.Write(sessionNameLengthByte)
	// session name
	build.Write(sessionNameBytes)

	// topicAttribute
	build.Write(map2Bytes(msg.TopicAttribute))

	// appProps
	build.Write(map2Bytes(msg.AppProps))

	// payload
	lengthPayloadBytes := []byte{0, 0, 0, 0}
	payloadBytes := []byte(msg.Body)
	lengthOfPayload := len(payloadBytes)
	lengthPayloadBytes[3] = byte(lengthOfPayload)
	lengthPayloadBytes[2] = byte(lengthOfPayload >> 8)
	lengthPayloadBytes[1] = byte(lengthOfPayload >> 16)
	lengthPayloadBytes[0] = byte(lengthOfPayload >> 24)
	build.Write(lengthPayloadBytes)
	build.Write(payloadBytes)
	//log.Errorf("ProtoMsg2String:%s, msg:%v", build.String(), msg)

	return build.Bytes()
}

// Bytes2protoMsg decodes bytes to protocol.ProtoMessage
func Bytes2protoMsg(inputBytes []byte) (protocol.ProtoMessage, error) {
	msg := &protocol.ProtoMessage{}
	if len(inputBytes) < 10 {
		err := errors.Errorf(constant.SystemInternalError, "inputstr invalid:[%s], bytes[%d]", string(inputBytes), inputBytes)
		return *msg, err
	}

	var pos uint32
	errFlag := []byte(inputBytes[0:1])[0]
	if byteOne == errFlag || byteTwo == errFlag {
		return *msg, Bytes2Error(inputBytes)
	}

	pos = 5
	msgIDBytes := inputBytes[pos : pos+8]
	msg.ID = binary.BigEndian.Uint64(msgIDBytes)
	pos = pos + 8

	// need reply
	msg.NeedReply = byte2bool(inputBytes[pos])
	pos++

	// need ack
	msg.NeedAck = byte2bool(inputBytes[pos])
	pos++

	// length of session name
	lengthOfSessionName := uint32(binary.BigEndian.Uint16(inputBytes[pos : pos+2]))
	pos = pos + 2

	// session name
	msg.SessionName = string(inputBytes[pos : pos+lengthOfSessionName])
	pos = pos + lengthOfSessionName

	// topicAttribute
	totalTopicAttributeLength := binary.BigEndian.Uint32(inputBytes[pos : pos+4])
	pos = pos + 4

	if totalTopicAttributeLength > 0 {
		msg.TopicAttribute = bytes2Map(inputBytes[pos:pos+totalTopicAttributeLength], totalTopicAttributeLength)
		pos = pos + totalTopicAttributeLength
	} else {
		msg.TopicAttribute = make(map[string]string)
	}

	// appProps
	appPropsLength := binary.BigEndian.Uint32(inputBytes[pos : pos+4])
	pos = pos + 4

	if appPropsLength > 0 {
		msg.AppProps = bytes2Map(inputBytes[pos:pos+appPropsLength], appPropsLength)
		pos = pos + appPropsLength
	} else {
		msg.AppProps = make(map[string]string)
	}

	// payload
	msg.Body = string(inputBytes[pos+4:])

	return *msg, nil
}
