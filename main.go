package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var dbpath string

func init() {
	flag.StringVar(&dbpath, "db", os.ExpandEnv("$PWD"), "path to where sqlite db file")
}

func main() {
	flag.Parse()

	blockdb, accountdb, witnessdb := initDB()
	defer func() {
		blockdb.Close()
		accountdb.Close()
		witnessdb.Close()
	}()

	gin.DisableConsoleColor()

	server := gin.Default()
	server.Static("/statics", "./statics")
	server.Delims("{{!#", "#!}}")
	server.LoadHTMLGlob("templates/*")

	server.GET("/", func(ctx *gin.Context) {
		top10 := []Account{}
		accountdb.Select("address, balance").Order("balance desc").Limit(10).Find(&top10)
		top10Data := gin.H{}
		for _, item := range top10 {
			top10Data[item.Address] = item.Balance
		}

		top10bytes, err := json.Marshal(top10Data)
		if err != nil {

		}

		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"top10": string(top10bytes),
		})
	})

	server.GET("/getindexdata", func(ctx *gin.Context) {
		accountNum := 0
		accountdb.Table(Account{}.TableName()).Count(&accountNum)

		nodeNum := 0
		witnessdb.Table(Nodes{}.TableName()).Count(&nodeNum)

		b := &Block{}
		blockdb.Select("number").Order("number desc").First(b)

		ctx.JSON(200, gin.H{
			"accounts":   accountNum,
			"nodes":      nodeNum,
			"lastheight": b.Number,
		})
	})

	server.GET("/blocklist", func(ctx *gin.Context) {
		size, err := strconv.ParseInt(ctx.DefaultQuery("size", "25"), 10, 64)
		if err != nil {
			size = 25
		}
		page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			page = 1
		}

		var blocks []Block
		blockdb.Select("number, witnessaddr").Order("number desc").Offset(((page - 1) * size)).Limit(size).Find(&blocks)

		var blocklist = make([]gin.H, 0)
		for _, blk := range blocks {
			cnt := 0
			blockdb.Table(Tx{}.TableName()).Where("blocknumber = ?", blk.Number).Count(&cnt)
			blocklist = append(blocklist,
				gin.H{"num": blk.Number, "wit": blk.Witnessaddr, "tx_num": cnt})
		}

		ctx.JSON(200, gin.H{
			"blocklist": blocklist,
		})
	})

	server.GET("/block", func(ctx *gin.Context) {
		num, exists := ctx.GetQuery("num")
		if !exists {
			ctx.Redirect(204, "/")
			return
		}
		number, err := strconv.ParseInt(num, 10, 64)
		if err != nil {
			ctx.Redirect(204, "/")
			return
		}

		blk := Block{}
		txs := make([]Tx, 0)
		blockdb.Find(&blk, number)
		blockdb.Table(Tx{}.TableName()).Where("blocknumber = ?", number).Find(&txs)

		t := time.Unix(blk.Timestamp/1000, 0)
		ctx.HTML(200, "block.html", gin.H{
			"num":         blk.Number,
			"trieroot":    blk.Trieroot,
			"time":        t.String(),
			"parent_hash": blk.Parenthash,
			"witness":     blk.Witnessaddr,
			"num_txs":     len(txs),

			"txs": txs,
		})
	})

	server.GET("/blocks", func(ctx *gin.Context) {
		cnt := 0
		blockdb.Table("block").Count(&cnt)
		pages := cnt / 25
		ctx.HTML(200, "blocklist.html", gin.H{
			"pages": pages,
		})
	})

	server.GET("/nodes", func(ctx *gin.Context) {
		cnt := 0
		witnessdb.Table(Nodes{}.TableName()).Count(&cnt)
		pages := cnt / 100
		ctx.HTML(200, "nodes.html", gin.H{
			"pages": pages,
		})
	})

	server.GET("/nodelist", func(ctx *gin.Context) {
		size, err := strconv.ParseInt(ctx.DefaultQuery("size", "25"), 10, 64)
		if err != nil {
			size = 25
		}
		page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			page = 1
		}

		var nodes []Nodes
		witnessdb.Offset(((page - 1) * size)).Limit(size).Find(&nodes)

		var nodelist = make([]gin.H, 0)
		for _, nd := range nodes {
			nodelist = append(nodelist,
				gin.H{"host": nd.Host, "port": nd.Port})
		}

		ctx.JSON(200, gin.H{
			"list": nodelist,
		})
	})

	server.GET("/account", func(ctx *gin.Context) {
		addr, exists := ctx.GetQuery("address")
		if !exists {
			ctx.Redirect(204, "/")
			return
		}

		account := &Account{}
		accountdb.Where("address = ?", addr).First(account)

		aType := ""
		switch account.Type {
		case 1:
			aType = "AssetIssue"
		case 2:
			aType = "Contract"
		default:
			aType = "Normal"
		}

		txs := []Tx{}
		blockdb.Table(Tx{}.TableName()).Where("owneraddr = ?", addr).Find(&txs)

		txlist := make([]gin.H, 0, len(txs))
		for _, t := range txs {
			txlist = append(txlist, gin.H{"toaddr": t.Toaddr, "contracts": t.Contracts})
		}

		ctx.HTML(200, "account.html", gin.H{
			"address": addr,
			"balance": account.Balance,
			"type":    aType,
			"name":    account.AccountName,
			"asset":   account.Asset,
			"votes":   account.Votes,
			"num_txs": len(txs),
			"txs":     txlist,
		})
	})

	server.GET("/accounts", func(ctx *gin.Context) {
		cnt := 0
		accountdb.Table(Account{}.TableName()).Count(&cnt)
		pages := cnt / 25
		ctx.HTML(200, "accountlist.html", gin.H{
			"pages": pages,
		})
	})

	server.GET("/accountlist", func(ctx *gin.Context) {
		size, err := strconv.ParseInt(ctx.DefaultQuery("size", "25"), 10, 64)
		if err != nil {
			size = 25
		}
		page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
		if err != nil {
			page = 1
		}

		var accounts []Account
		accountdb.Order("balance desc").Offset(((page - 1) * size)).Limit(size).Find(&accounts)

		var accountlist = make([]gin.H, 0)
		for _, act := range accounts {
			item := gin.H{
				"id":          act.ID,
				"accountname": act.AccountName,
				"type":        act.Type,
				"address":     act.Address,
				"balance":     act.Balance,
				"votes":       act.Votes,
				"asset":       act.Asset,
				"latest":      act.LatestOperationTime,
			}

			switch act.Type {
			case 1:
				item["typename"] = "AssetIssue"
			case 2:
				item["typename"] = "Contract"
			default:
				item["typename"] = "Normal"
			}
			accountlist = append(accountlist, item)
		}

		ctx.JSON(200, gin.H{
			"accountlist": accountlist,
		})
	})

	server.GET("/search", func(ctx *gin.Context) {
		search := ctx.DefaultQuery("search", "")
		if search == "" {
			ctx.Redirect(204, "/")
		}

		var (
			block   = &Block{}
			account = &Account{}
			tx      = &Tx{}
		)

		if blkNum, err := strconv.ParseInt(search, 10, 64); err == nil {
			blockdb.First(block, blkNum)
			if block != nil {
				ctx.Redirect(302, "/block?num="+search)
			} else {
				ctx.Redirect(404, "/404")
			}
			return
		}

		accountdb.Table(account.TableName()).Where("address = ?", search).First(account)
		if account != nil {
			ctx.Redirect(302, "/account?address="+search)
			return
		}

		blockdb.Where("hash = ?", search).First(tx)
		if tx != nil {

			return
		}
	})

	server.Run()
}
