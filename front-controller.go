package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/netroby/mysql"
	"html/template"
	"net/http"
	"strconv"
)

type FrontController struct {
}

func (fc *FrontController) HomeCtr(c *gin.Context) {
	content := c.DefaultQuery("content", "")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		fmt.Println(err)
	}
	page -= 1
	if page < 0 {
		page = 0
	}

	prev_page := page
	if prev_page < 1 {
		prev_page = 1
	}
	next_page := page + 2

	rpp := 20
	offset := page * rpp
	CKey := fmt.Sprintf("%s-home-page-%d-rpp-%d", GetMinutes(), page, rpp)
	var blogList string
	val, ok := Cache.Get(CKey)
	if val != nil && ok == true {
		fmt.Println("Ok, we found cache, Cache Len: ", Cache.Len())
		blogList = val.(string)
	} else {
		rows, err := DB.Query("Select aid, content from article where publish_status = 1 order by aid desc limit ? offset ? ", &rpp, &offset)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		var (
			aid     int
			content sql.NullString
		)
		for rows.Next() {
			err := rows.Scan(&aid, &content)
			if err != nil {
				fmt.Println(err)
			}
			blogList += fmt.Sprintf(
				"<div>%s<hr /></div>",
				content.String,
			)
		}
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
		go func(CKey string, blogList string) {
			Cache.Add(CKey, blogList)
		}(CKey, blogList)
	}
	session := sessions.Default(c)
	username := session.Get("username")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"content":          content,
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
		"bloglist":         template.HTML(blogList),
		"username":         username,
		"prev_page":        prev_page,
		"next_page":        next_page,
	})
}
