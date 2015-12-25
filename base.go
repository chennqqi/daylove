package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/naoina/toml"
	_ "github.com/netroby/mysql"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type ShowMessage interface {
	ShowMessage(c *gin.Context)
}
type msg struct {
	msg string
}
type umsg struct {
	msg string
	url string
}

type VBlogItem struct {
	aid            int
	title          sql.NullString
	content        sql.NullString
	publish_time   sql.NullString
	publish_status sql.NullInt64
	views          int
}

/*
	Show message with template
*/
func (m *msg) ShowMessage(c *gin.Context) {
	c.HTML(http.StatusOK, "message.html", gin.H{
		"message": template.HTML(m.msg),
	})
}

func (m *umsg) ShowMessage(c *gin.Context) {

	c.HTML(http.StatusOK, "message.html", gin.H{
		"message": template.HTML(m.msg),
		"url":     m.url,
	})
}

func GetMinutes() string {
	return time.Now().Format("200601021504")
}

func GetDB(config *appConfig) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=30s&charset=utf8mb4",
		config.Db_user, config.Db_password, config.Db_host, config.Db_port, config.Db_name))
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)
	return db
}

type appConfig struct {
	Db_host          string
	Db_port          int
	Db_name          string
	Db_user          string
	Db_password      string
	Admin_user       string
	Admin_password   string
	Site_name        string
	Site_description string
}

func GetConfig() *appConfig {
	f, err := os.Open("config.toml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var config appConfig
	if err := toml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	return &config
}
