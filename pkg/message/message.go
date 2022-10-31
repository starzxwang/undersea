package message

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"undersea/pkg/api"
	"undersea/pkg/log"
)

const (
	MesTypeLogin               = "LOGIN"                  // 登录消息
	MesTypeReplyLogin          = "REPLY_LOGIN"            // 登录ack消息
	MesTypeExit                = "EXIT"                   // 注销消息
	MesTypePeerText            = "PEER_TEXT"              // 单聊消息
	MesTypeReplyPeerText       = "REPLY_PEER_TEXT"        // 单聊ack消息
	MesTypeGroupText           = "GROUP_TEXT"             // 普通群聊
	MesTypeReplyGroupText      = "REPLY_GROUP_TEXT"       // 普通群聊ack消息
	MesTypeSuperGroupText      = "SUPER_GROUP_TEXT"       // 超级群聊
	MesTypeReplySuperGroupText = "REPLY_SUPER_GROUP_TEXT" // 超级群聊ack消息
	MesTypeReplyExit           = "REPLY_EXIT"             // 注销
	MesTypeHeartBeat           = "HEART_BEAT"             // im的心跳
	MesTypePickIp              = "PICK_IP"                // 客户端向im_balance获取ip
	MesTypeReplyPickIp         = "REPLY_PICK_IP"
	MesTypeReplyException      = "REPLY_EXCEPTION" // 服务内部异常导致
)

type Message struct {
	Id   string `json:"msg_id,omitempty"`
	Type string `json:"type"`
	Data string `json:"data"`
	Len  int    `json:"len"`
}

type ReplyMessageData struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type HeartBeatMessage struct {
	ServiceName string `json:"service_name"`
	Ip          string `json:"ip"`
}

type PickIpMessage struct {
	ServiceName string `json:"service_name"`
	Uid         int    `json:"uid"`
}

type LoginMessage struct {
	Uid int `json:"uid"`
}

type PickIpReplyMessage struct {
	Ip  string `json:"ip"`
	Uid int    `json:"uid"`
}

func ConvertBytes2Message(ctx context.Context, data []byte) (mes *Message, err error) {
	err = json.Unmarshal(data, &mes)
	if err != nil {
		log.E(ctx, err).Msgf("json.Unmarshal err")
		return
	}

	if mes.Len != len(mes.Data) {
		err = errors.New("json.Unmarshal err")
		log.E(ctx, err).Msgf("json.Unmarshal err")
		return
	}

	return
}

// 组装消息
func formatMessage(ctx context.Context, msgData interface{}, msgType string, msgId string) (mes *Message, err error) {
	dataBytes, err := json.Marshal(msgData)
	if err != nil {
		log.E(ctx, err).Msgf("json.Marshal err")
		return
	}

	mes = &Message{
		Id:   msgId,
		Type: msgType,
		Data: string(dataBytes),
	}

	mes.Len = len(mes.Data)
	return
}

// 发送服务异常消息给前端
func SendExceptionWebsocketMessage(ctx context.Context, conn *websocket.Conn) (err error) {
	return SendWebSocketMessage(ctx, conn, &ReplyMessageData{
		Code:    api.CodeException,
		Message: "服务开小差了，请稍候重试",
	}, MesTypeReplyException, "")
}

// 发送websocket消息
func SendWebSocketMessage(ctx context.Context, conn *websocket.Conn, msgData *ReplyMessageData, msgType, msgId string) (err error) {
	mes, err := formatMessage(ctx, msgData, msgType, msgId)
	if err != nil {
		log.E(ctx, err).Msgf("formatMessage err")
		return
	}

	mesBytes, err := json.Marshal(mes)
	if err != nil {
		log.E(ctx, err).Msgf("json.Marshal err")
		return
	}

	log.I(ctx).Msgf("服务端返回客户端消息=%v", string(mesBytes))

	err = sendWebSocketByteMessage(ctx, conn, mesBytes)
	if err != nil {
		log.E(ctx, err).Msgf("sendByteMessage err")
		return
	}
	return
}

func sendWebSocketByteMessage(ctx context.Context, conn *websocket.Conn, mesBytes []byte) (err error) {
	if err = conn.WriteMessage(websocket.TextMessage, mesBytes); err != nil {
		//报错了
		log.E(ctx, err).Msgf("handlePeerTextMessage:conn.WriteMessage()")
		return
	}

	return
}
