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
	"strings"
)

type FrontController struct {
}

func (fc *FrontController) AboutCtr(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{		
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
	})
}
func (fc *FrontController) PingCtr(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
func (fc *FrontController) HomeCtr(c *gin.Context) {
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
		rows, err := DB.Query("Select aid, title from top_article where publish_status = 1 order by aid desc limit ? offset ? ", &rpp, &offset)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		var (
			aid   int
			title sql.NullString
		)
		for rows.Next() {
			err := rows.Scan(&aid, &title)
			if err != nil {
				fmt.Println(err)
			}
			blogList += fmt.Sprintf(
				"<li><a href=\"/view/%d\">%s</a></li>",
				aid,
				title.String,
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
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
		"bloglist":         template.HTML(blogList),
		"username":         username,
		"prev_page":        prev_page,
		"next_page":        next_page,
	})
}

func (fc *FrontController) SearchCtr(c *gin.Context) {
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
	keyword := c.DefaultQuery("keyword", "")
	fmt.Println(keyword)
	if len(keyword) <= 0 {
		(&msg{"Keyword can not empty"}).ShowMessage(c)
		return
	}
	orig_keyword := keyword
	keyword = strings.Replace(keyword, " ", "%", -1)

	var blogList string
	rpp := 20
	offset := page * rpp
	rows, err := DB.Query(
		"Select aid, title from top_article where publish_status = 1 and (title like ? or content like ?) order by aid desc limit ? offset ? ",
		"%"+keyword+"%", "%"+keyword+"%", &rpp, &offset)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var (
		aid   int
		title sql.NullString
	)
	for rows.Next() {
		err := rows.Scan(&aid, &title)
		if err != nil {
			fmt.Println(err)
		}
		blogList += fmt.Sprintf(
			"<li><a href=\"/view/%d\">%s</a></li>",
			aid,
			title.String,
		)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}
	session := sessions.Default(c)
	username := session.Get("username")

	c.HTML(http.StatusOK, "search.html", gin.H{
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
		"bloglist":         template.HTML(blogList),
		"keyword":          orig_keyword,
		"username":         username,
		"prev_page":        prev_page,
		"next_page":        next_page,
	})
}

func (fc *FrontController) ViewAltCtr(c *gin.Context) {
	id := c.DefaultQuery("id", "0")
	c.Redirect(301, fmt.Sprintf("/view/%s", id))
}

func (fc *FrontController) CountViewCtr(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println("Can not get id")
		return
	}
	_, err = DB.Exec("update top_article set views=views+1 where aid = ? limit 1", &id)
	if err != nil {
		fmt.Println(err)
	}
	c.String(http.StatusOK, "[1]")
}

func (fc *FrontController) ViewCtr(c *gin.Context) {
	id := c.Param("id")
	var blog VBlogItem
	CKey := fmt.Sprintf("%s-blogitem-%d", GetMinutes(), id)
	val, ok := Cache.Get(CKey)
	if val != nil && ok == true {
		blog = val.(VBlogItem)
	} else {
		rows, err := DB.Query("select aid, title, content, publish_time, publish_status, views from top_article where aid = ? limit 1", &id)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		var ()
		for rows.Next() {
			err := rows.Scan(&blog.aid, &blog.title, &blog.content, &blog.publish_time, &blog.publish_status, &blog.views)
			if err != nil {
				fmt.Println(err)
			}
		}
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
		go func(CKey string, blog VBlogItem) {
			Cache.Add(CKey, blog)
		}(CKey, blog)
	}
	session := sessions.Default(c)
	username := session.Get("username")
	c.HTML(http.StatusOK, "view.html", gin.H{
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
		"aid":              blog.aid,
		"title":            blog.title.String,
		"content":          template.HTML(blog.content.String),
		"publish_time":     blog.publish_time.String,
		"views":            blog.views,
		"username":         username,
	})

}
