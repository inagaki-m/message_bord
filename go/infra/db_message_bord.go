package infra

import (
	"context"
	"log"
	"messageBord/go/model"
	"messageBord/go/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var _ repository.MessageBoardRepository = new(DbMessageBoardRepository)

type DbMessageBoardRepository struct {
	DB *sqlx.DB
}

func (m *DbMessageBoardRepository) RegisterMessageInfo(messageInfo *model.MessageInfo, ctx context.Context) error {
	tx, err := m.DB.BeginTxx(ctx, nil)
	if err != nil {
		log.Fatalln("[ERROR] BeginTxx: ", err)
		return err
	}

	// データベースにクエリを送信する
	query := `INSERT INTO messages(
			name,
			message,
			createTime
		)
		VALUES(
			:name,
			:message,
			NOW()
		)`
	if _, err = tx.NamedExecContext(ctx, query, messageInfo); err != nil {
		log.Fatalln("[ERROR] NamedExecContext: ", err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *DbMessageBoardRepository) GetMessageList(ctx context.Context) ([]*model.MessageInfo, error) {
	rows, err := m.DB.QueryxContext(ctx, "SELECT name, message, createTime FROM messages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messageInfoList []*model.MessageInfo

	for rows.Next() {
		var messageInfo model.MessageInfo

		err := rows.Scan(
			&messageInfo.Name,
			&messageInfo.Message,
			&messageInfo.CreateTime,
		)
		if err != nil {
			return nil, err
		}
		messageInfoList = append(messageInfoList, &messageInfo)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return messageInfoList, nil
	}
}
