package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"messageBord/go/infra"
	"messageBord/go/model"
	"messageBord/go/repository"
	"os"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
)

func main() {
	///////////////////////////////////////////////
	// main()の内容
	db, err := sqlx.Open("mysql", "root:myrootpassword@tcp(127.0.0.1:3306)/test_1?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Fatalln("[ERROR]sqlx.Open: ", err)
		panic(err.Error())
	}
	defer db.Close()

	ctx := context.Background()

	var messageRepo repository.MessageBoardRepository
	// messageRepo = &infra.DbMessageBoardRepository{DB: db} // DB
	messageRepo = &infra.FileMessageBoardRepository{FilePath: "./messageList.json"} // File

	if len(os.Args) > 1 {
		///////////////////////////////////////////////
		// 以下usecaseの内容
		messageList, err := useRepositoryGetMessageList(messageRepo, ctx)
		if err != nil {
			log.Fatalln("[ERROR]sqlx.Open: ", err)
		}
		// 一旦Printfで出力
		for _, m := range messageList {
			fmt.Printf("[Name]%s, [Message]%s, [CreateTime]%s \n", m.Name, m.Message, m.CreateTime)
		}

	} else {
		userName, err := input("名前")
		if err {
			return
		}
		message, err := input("メッセージ")
		if err {
			return
		}

		///////////////////////////////////////////////
		// 以下usecaseの内容
		messageInfo := &model.MessageInfo{
			Name:    userName,
			Message: message,
		}
		useRepositoryRegisterMessageInfo(messageRepo, messageInfo, ctx)
	}
}

// [memo]
// messageRepoをrepository.MessageBoardRepositoryにすることによりinfraのどのメソッド使っても良いようにする
// *infra.DbMessageBoardRepositoryにした場合、使いたいinfraが変わった場合、呼び出し元も同じにする必要がある
func useRepositoryRegisterMessageInfo(messageRepo repository.MessageBoardRepository, messageInfo *model.MessageInfo, ctx context.Context) {
	repositoryErr := messageRepo.RegisterMessageInfo(messageInfo, ctx)
	if repositoryErr != nil {
		return
	}
}

func useRepositoryGetMessageList(messageRepo repository.MessageBoardRepository, ctx context.Context) ([]*model.MessageInfo, error) {
	messageList, repositoryErr := messageRepo.GetMessageList(ctx)
	if repositoryErr != nil {
		return nil, repositoryErr
	}
	return messageList, nil
}

func input(item string) (string, bool) {
	fmt.Printf("%sを入力してください\n", item)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Text() == "" {
		fmt.Println("１文字以上を入力してください")
		return "", true
	}

	if utf8.RuneCountInString(scanner.Text()) >= 100 {
		fmt.Println("utfLen: ", utf8.RuneCountInString(scanner.Text()))
		fmt.Println("100文字以内で入力してください")
		return "", true
	}
	return scanner.Text(), false
}
