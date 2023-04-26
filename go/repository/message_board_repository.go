package repository

import (
	"messageBord/go/model"
)

type MessageBoardRepository interface {
	RegisterMessageInfo(messageInfo *model.MessageInfo) error
	GetMessageList() ([]model.MessageInfo, error)
}
