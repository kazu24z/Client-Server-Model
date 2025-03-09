package protocol

// TCRPHeader はTCRPヘッダーを表します。
type TCRPHeader struct {
	RoomNameSize         uint8
	Operation            uint8
	State                uint8
	OperationPayloadSize [29]byte
}

// TCRPMessage はTCRPメッセージを表します。
type TCRPMessage struct {
	Header TCRPHeader
	Body   []byte
}

// EncodeTCRPMessage はTCRPメッセージをバイト列にエンコードします。
func EncodeTCRPMessage(msg TCRPMessage) ([]byte, error) {
	// ここでエンコード処理を実装します。
	return nil, nil
}

// DecodeTCRPMessage はバイト列をTCRPメッセージにデコードします。
func DecodeTCRPMessage(data []byte) (TCRPMessage, error) {
	// ここでデコード処理を実装します。
	return TCRPMessage{}, nil
}
