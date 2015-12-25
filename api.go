package main

import (
	"github.com/gin-gonic/gin"

	"database/sql"
	"fmt"
	"github.com/jaytaylor/html2text"
	_ "github.com/netroby/mysql"
	"log"
	"net/http"
	"strconv"
)

type api struct {
}

type apiBlogList struct {
	Aid   string `form:"aid" json:"aid"  binding:"required"`
	Title string `form:"title" json:"title"  binding:"required"`
}

func (a *api) index(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		log.Fatal(err)
	}
	page -= 1
	if page < 0 {
		page = 0
	}

	prev_page := page
	if prev_page < 1 {
		prev_page = 1
	}

	rpp := 20
	offset := page * rpp
	CKey := fmt.Sprintf("api-home-page-%d-rpp-%d", page, rpp)
	var blogListSlice []apiBlogList
	val, ok := Cache.Get(CKey)
	if val != nil && ok == true {
		blogListSlice = val.([]apiBlogList)
	} else {
		rows, err := DB.Query("Select aid, title from top_article where publish_status = 1 order by aid desc limit ? offset ? ", &rpp, &offset)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var (
			aid   sql.NullString
			title sql.NullString
		)
		blogListSlice = make([]apiBlogList, 0) //Must be zero slice
		var aBlog apiBlogList
		for rows.Next() {
			err := rows.Scan(&aid, &title)
			if err != nil {
				log.Fatal(err)
			}
			aBlog.Aid = aid.String
			aBlog.Title = title.String
			blogListSlice = append(blogListSlice, aBlog)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		Cache.Add(CKey, blogListSlice)
	}
	c.JSON(http.StatusOK, blogListSlice)
}

type apiBlogItem struct {
	Aid     string `form:"aid" json:"aid"  binding:"required"`
	Title   string `form:"title" json:"title"  binding:"required"`
	Content string `form:"content" json:"content"  binding:"required"`
}

func (a *api) view(c *gin.Context) {
	aid, err := strconv.Atoi(c.Param("id"))
	fmt.Println(aid)
	if err != nil {
		log.Fatal(err)
	}
	CKey := fmt.Sprintf("api-view-aid-%d", aid)
	var b apiBlogItem
	val, ok := Cache.Get(CKey)
	fmt.Println(val)
	if val != nil && ok == true {
		b = val.(apiBlogItem)
	} else {
		rows, err := DB.Query("Select aid, title, content from top_article where aid =  ? limit 1 ", &aid)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var (
			aid     sql.NullString
			title   sql.NullString
			content sql.NullString
		)
		for rows.Next() {
			err := rows.Scan(&aid, &title, &content)
			if err != nil {
				fmt.Println(err)
			}
			b.Aid = aid.String
			b.Title = title.String
			b.Content = content.String
		}
		fmt.Println(b)
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		Cache.Add(CKey, b)
	}
	fmt.Println(b)
	text, err := html2text.FromString(b.Content)
	if err != nil {
		fmt.Println("Error when convert html to text")
	} else {
		b.Content = text
	}
	c.JSON(http.StatusOK, b)
}
