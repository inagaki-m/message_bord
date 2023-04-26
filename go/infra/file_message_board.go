package infra

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"messageBord/go/model"
	"messageBord/go/repository"
	"os"
	"time"
)

var _ repository.MessageBoardRepository = new(FileMessageBoardRepository)

type FileMessageBoardRepository struct {
	FilePath string
}

func (m *FileMessageBoardRepository) RegisterMessageInfo(messageInfo *model.MessageInfo) error {
	inputData, err := GetInputMessageList(m.FilePath)
	if err != nil {
		return err
	}
	messageInfo.CreateTime = time.Now()
	inputData = append(inputData, *messageInfo)

	// JSONに変換
	jsonData, err := json.Marshal(inputData)
	if err != nil {
		fmt.Println("JSONへ変換失敗")
		return err
	}

	if err := ioutil.WriteFile(m.FilePath, jsonData, os.ModePerm); err != nil {
		fmt.Println("ファイル書き込み失敗")
		return err
	}
	return nil
}

func (m *FileMessageBoardRepository) GetMessageList() ([]model.MessageInfo, error) {
	messageInfo := make([]model.MessageInfo, 0)
	bytes, err := ioutil.ReadFile(m.FilePath)
	if err != nil {
		fmt.Println("ファイル内容取得 失敗", err)
		return nil, err
	}

	if err := json.Unmarshal(bytes, &messageInfo); err != nil {
		fmt.Println("ファイル構造体に変換 失敗", err)
		return nil, err
	}
	return messageInfo, nil
}

func GetInputMessageList(filePath string) ([]model.MessageInfo, error) {
	messageInfo := make([]model.MessageInfo, 0)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("ファイル内容取得 失敗", err)
		return nil, err
	}

	if err := json.Unmarshal(bytes, &messageInfo); err != nil {
		fmt.Println("ファイル構造体に変換 失敗", err)
		return nil, err
	}
	return messageInfo, nil
}
