package main

import (
	// "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type (
	Block struct {
		Number      int64 `gorm:"primary_key"`
		Witnessid   int64
		Timestamp   int64
		Parenthash  string `gorm:"type:varchar(64)"`
		Trieroot    string `gorm:"type:varchar(64)"`
		Witnessaddr string `gorm:"type:varchar(64)"`
		witnesssig  string `gorm:"type:varchar(128)"`
	}

	Account struct {
		ID                  int64  `gorm:"column:Id;type:bigint;primary_key;auto_increment"`
		AccountName         string `gorm:"column:account_name;type:text;not null"`
		Type                int32  `gorm:"not null"`
		Address             string `gorm:"type:varchar(64);unique;not null"`
		Balance             int64  `gorm:"not null"`
		Votes               string `gorm:"type:text"`
		Asset               string `gorm:"type:text"`
		LatestOperationTime int64  `gorm:"column:latest_operation_time;not null"`
	}

	Tx struct {
		ID             int64 `gorm:"primary_key"`
		Blocknumber    int64
		Refblocknumber int64
		Owneraddr      string `gorm:"type:varchar(64);index:txfrom"`
		Toaddr         string `gorm:"type:varchar(64);index:txto"`
		Expiration     int64
		Timestamp      int64
		Refblockhash   string `gorm:"type:varchar(64)"`
		Scripts        string `gorm:"type:text"`
		Data           string `gorm:"type:text"`
		Contracts      string `gorm:"type:text"`
		Sigs           string `gorm:"type:text"`
	}

	Witness struct {
		ID             int64  `gorm:"column:Id;type:int;primary_key;auto_increment"`
		Address        string `gorm:"type:varchar(64);unique;not null"`
		VoteCount      int64  `gorm:"column:voteCount;type:int;not null"`
		PubKey         string `gorm:"column:pubKey;type:varchar(128)"`
		Url            string `gorm:"type:text"`
		TotalProduced  int64  `gorm:"column:totalProduced;type:int(10)"`
		TotalMissed    int64  `gorm:"column:totalMissed;type:int(10)"`
		LatestBlockNum int64  `gorm:"column:latestBlockNum;type:int(10)"`
		LatestSlotNum  int64  `gorm:"column:latestSlotNum;type:int(10)"`
		IsJobs         bool   `gorm:"column:isJobs"`
	}

	Nodes struct {
		ID   int64 `gorm:"column:Id;primary_key;auto_increment"`
		Host string
		Port int32
	}
)

func (b Block) TableName() string {
	return "block"
}

func (t Tx) TableName() string {
	return "tx"
}

func (a Account) TableName() string {
	return "tron_accounts"
}

func (w Witness) TableName() string {
	return "tron_witness"
}

func (n Nodes) TableName() string {
	return "tron_nodes"
}
