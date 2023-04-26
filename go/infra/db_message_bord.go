package infra

import (
	"context"
	"log"
	"messageBord/go/model"
	"messageBord/go/repository"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var _ repository.MessageBoardRepository = new(DbMessageBoardRepository)

type DbMessageBoardRepository struct {
	DB *sqlx.DB
}

func (m *DbMessageBoardRepository) RegisterMessageInfo(messageInfo *model.MessageInfo) error {
	ctx := context.Background()

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
		return err
	}
	tx.Commit()
	return nil
}

func (m *DbMessageBoardRepository) GetMessageList() ([]model.MessageInfo, error) {
	db, err := sqlx.Open("mysql", "root:myrootpassword@tcp(127.0.0.1:3306)/test_1")
	if err != nil {
		log.Fatalln("[ERROR]sqlx.Open: ", err)
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Queryx("SELECT name, message, createTime FROM messages")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var messageInfoList []model.MessageInfo

	const mysqlDatetimeLayout = "2006-01-02 15:04:05"

	for rows.Next() {
		var messageInfo model.MessageInfo
		var createTimeStr string

		err := rows.Scan(
			&messageInfo.Name,
			&messageInfo.Message,
			&createTimeStr,
		)
		if err != nil {
			log.Fatalln("[ERROR] rows.Scan: ", err)
			panic(err.Error())

		}
		// 時間文字列をパースして、MessageInfoのCreateTimeに代入する
		createTime, err := time.Parse(mysqlDatetimeLayout, createTimeStr)
		if err != nil {
			log.Fatalln("[ERROR] time.Parse: ", err)
			panic(err.Error())
		}
		messageInfo.CreateTime = createTime

		messageInfoList = append(messageInfoList, messageInfo)
	}
	return messageInfoList, nil
}
