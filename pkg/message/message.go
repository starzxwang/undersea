package message

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net"
	error_v2 "undersea/pkg/err"
	"undersea/pkg/log"
)

const (
	MesTypePeerText       = "PEER_TEXT" // 单聊消息
	MesTypeReplyPeerText  = "REPLY_PEER_TEXT" // 单聊ack消息
	MesTypeGroupText      = "GROUP_TEXT" // 普通群聊
	MesTypeReplyGroupText = "REPLY_GROUP_TEXT" // 普通群聊ack消息
	MesTypeSuperGroupText = "SUPER_GROUP_TEXT" // 超级群聊
	MesTypeReplySuperGroupText = "REPLY_SUPER_GROUP_TEXT" // 超级群聊ack消息
	MesTypeReplyExit      = "REPLY_EXIT" // 注销
	MesTypeHeartBeat      = "HEART_BEAT" // im的心跳
	MesTypePickNodeIp = "PICK_IP" // 客户端向im_balance获取ip
	MesTypeReplyPickNodeIp = "REPLY_PICK_IP"
)

type Message struct {
	Id         string `json:"msg_id,omitempty"`
	Type       string `json:"type"`
	Data       string `json:"data"`
	Length     int    `json:"length"`
}

type PickNodeIpReplyMessage struct {
	Ip string `json:"ip"`
}

func ConvertBytes2Message(ctx context.Context, data []byte) (mes *Message, err error) {
	err = json.Unmarshal(data, &mes)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "json.Unmarshal err")
		return
	}

	if mes.Length != len(mes.Data) {
		err = error_v2.PrintError(ctx, nil, "json.Unmarshal err")
		return
	}

	return
}

// 组装消息
func formatMessage(ctx context.Context, msgData interface{}, msgType string, msgId string) (mes *Message, err error) {
	dataBytes, err := json.Marshal(msgData)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "json.Marshal err")
		return
	}

	mes = &Message{
		Id: msgId,
		Type: msgType,
		Data: string(dataBytes),
	}

	mes.Length = len(mes.Data)
	return
}

// 发送tcp消息
func SendTcpMessage(ctx context.Context, conn net.Conn, msgData interface{}, msgType, msgId string) (err error) {
	mes, err := formatMessage(ctx, msgData, msgType, msgId)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "formatMessage err")
		return
	}

	mesBytes, err := json.Marshal(mes)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "json.Marshal err")
		return
	}

	_, err = conn.Write(mesBytes)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "conn.Write err")
		return
	}

	return
}

//发送websocket消息
func SendWebSocketMessage(ctx context.Context, conn *websocket.Conn, msgData interface{}, msgType, msgId string) (err error) {
	mes, err := formatMessage(ctx, msgData, msgType, msgId)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "formatMessage err")
		return
	}


	mesBytes, err := json.Marshal(mes)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "json.Marshal err")
		return
	}

	log.I(ctx).Msgf("服务端返回客户端消息=%v", mes)

	err = sendWebSocketByteMessage(ctx, conn, mesBytes)
	if err != nil {
		err = error_v2.PrintError(ctx, err, "sendByteMessage err")
		return
	}
	return
}

func sendWebSocketByteMessage(ctx context.Context, conn *websocket.Conn, mesBytes []byte) (err error) {
	if err = conn.WriteMessage(websocket.TextMessage, mesBytes); err != nil {
		//报错了
		err = error_v2.PrintError(ctx, err, "handlePeerTextMessage:conn.WriteMessage()")
		return
	}

	return
}
