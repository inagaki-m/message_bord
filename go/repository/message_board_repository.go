package repository

import (
	"context"
	"messageBord/go/model"
)

type MessageBoardRepository interface {
	RegisterMessageInfo(messageInfo *model.MessageInfo, ctx context.Context) error
	GetMessageList(ctx context.Context) ([]*model.MessageInfo, error)
}
