package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/netroby/feeds"
	_ "github.com/netroby/mysql"
	"net/http"
	"strconv"
	"time"
)

// RSS controll group
type RSS struct {
}

func (rss *RSS) Alter(c *gin.Context) {
	c.Redirect(301, "/rss")
}

// Out Render and output RSS
func (rss *RSS) Out(c *gin.Context) {
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

	rpp := 20
	offset := page * rpp
	rows, err := DB.Query("Select aid, title, content, publish_time from top_article where publish_status = 1 order by aid desc limit ? offset ? ", &rpp, &offset)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var (
		aid          sql.NullString
		title        sql.NullString
		content      sql.NullString
		publish_time sql.NullString
	)

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "HardCoder",
		Link:        &feeds.Link{Href: "https://www.netroby.com"},
		Description: "Opensource , linux, golang",
		Created:     now,
	}
	feed.Items = make([]*feeds.Item, 0)
	for rows.Next() {
		err := rows.Scan(&aid, &title, &content, &publish_time)
		if err != nil {
			fmt.Println(err)
			break
		}
		itemTime, err := time.Parse("2006-01-02 15:04:05", publish_time.String)
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       title.String,
			Link:        &feeds.Link{Href: fmt.Sprintf("https://www.netroby.com/view/%s", aid.String)},
			Description: content.String,
			Created:     itemTime,
		})
	}
	c.XML(http.StatusOK, (&feeds.Rss{feed}).FeedXml())
}
