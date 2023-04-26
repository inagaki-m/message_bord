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
	DB  *sqlx.DB
	Ctx context.Context
}

func (m *DbMessageBoardRepository) RegisterMessageInfo(messageInfo *model.MessageInfo) error {
	tx, err := m.DB.BeginTxx(m.Ctx, nil)
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
	if _, err = tx.NamedExecContext(m.Ctx, query, messageInfo); err != nil {
		log.Fatalln("[ERROR] NamedExecContext: ", err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (m *DbMessageBoardRepository) GetMessageList() ([]*model.MessageInfo, error) {
	rows, err := m.DB.Queryx("SELECT name, message, createTime FROM messages")
	if err != nil {
		panic(err.Error())
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
			log.Fatalln("[ERROR] rows.Scan: ", err)
			panic(err.Error())
		}
		messageInfoList = append(messageInfoList, &messageInfo)
	}
	return messageInfoList, nil
}
