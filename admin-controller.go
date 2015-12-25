package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang/groupcache/lru"
	_ "github.com/netroby/mysql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// AdminLoginForm is the login form for Admin
type AdminLoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

// BlogItem is the blog item
type BlogItem struct {
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
}
type EditBlogItem struct {
	Aid     string `form:"aid" binding:"required"`
	Title   string `form:"title" binding:"required"`
	Content string `form:"content" binding:"required"`
}

type AdminController struct {
}

// ListBlogCtr is list blogs for admin
func (ac *AdminController) ListBlogCtr(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	fmt.Println(username)
	if username == nil {
		(&umsg{"You have no permission", "/"}).ShowMessage(c)
		return
	}
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
	next_page := page + 2

	var blogList string
	rpp := 20
	offset := page * rpp
	log.Println(rpp)
	log.Println(offset)
	rows, err := DB.Query("Select aid, title, views from top_article where publish_status = 1 order by aid desc limit ? offset ? ", &rpp, &offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		aid   int
		title sql.NullString
		views int
	)
	for rows.Next() {
		err := rows.Scan(&aid, &title, &views)
		if err != nil {
			log.Fatal(err)
		}
		blogList += fmt.Sprintf(
			"<li><a href=\"/view/%d\">%s</a> Views(%d)    [<a href=\"/admin/editblog/%d\">Edit</a>] [<a href=\"/admin/deleteblog/%d\">Delete</a>]</li>",
			aid,
			title.String,
			views,
			aid,
			aid,
		)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	c.HTML(http.StatusOK, "admin.list.blog.html", gin.H{
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
		"bloglist":         template.HTML(blogList),
		"username":         username,
		"prev_page":        prev_page,
		"next_page":        next_page,
	})
}

func (ac *AdminController) EditBlogCtr(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		(&umsg{"You have no permission", "/"}).ShowMessage(c)
		return
	}
	id := c.Param("id")
	var blog VBlogItem
	CKey := fmt.Sprintf("blogitem-%s", id)
	val, ok := Cache.Get(CKey)
	if val != nil && ok == true {
		fmt.Println("Ok, we found cache, Cache Len: ", Cache.Len())
		blog = val.(VBlogItem)
	} else {
		rows, err := DB.Query("select aid, title, content, publish_time, publish_status, views  from top_article where aid = ?", &id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var ()
		for rows.Next() {
			err := rows.Scan(&blog.aid, &blog.title, &blog.content, &blog.publish_time, &blog.publish_status, &blog.views)
			if err != nil {
				log.Fatal(err)
			}
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		Cache.Add(CKey, blog)
	}
	c.HTML(http.StatusOK, "edit-blog.html", gin.H{
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
		"aid":              blog.aid,
		"title":            blog.title.String,
		"content":          template.HTML(blog.content.String),
		"publish_time":     blog.publish_time.String,
		"views":            blog.views,
	})
}

func (ac *AdminController) DeleteBlogCtr(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		(&umsg{"You have no permission", "/"}).ShowMessage(c)
		return
	}
	var BI EditBlogItem
	c.BindWith(&BI, binding.Form)
	if BI.Aid == "" {
		(&umsg{"Can not find the blog been delete", "/"}).ShowMessage(c)
		return
	}
	_, err := DB.Exec("delete from top_article where aid = ? limit 1", BI.Aid)
	if err == nil {
		Cache = lru.New(CacheSize)
		(&umsg{"Deleted Success", "/"}).ShowMessage(c)
	} else {
		(&umsg{"Failed to delete blog", "/"}).ShowMessage(c)
	}
}

func (ac *AdminController) AddBlogCtr(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		(&umsg{"You have no permission", "/"}).ShowMessage(c)
		return
	}
	c.HTML(http.StatusOK, "add-blog.html", gin.H{
		"site_name":        Config.Site_name,
		"site_description": Config.Site_description,
	})
}

func (ac *AdminController) SaveBlogEditCtr(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		(&umsg{"You have no permission", "/"}).ShowMessage(c)
		return
	}
	var BI EditBlogItem
	c.BindWith(&BI, binding.Form)
	if BI.Aid == "" {
		(&umsg{"Can not find the blog been edit", "/"}).ShowMessage(c)
		return
	}
	if BI.Title == "" {
		(&umsg{"Title can not empty", "/"}).ShowMessage(c)
		return
	}
	if BI.Content == "" {
		(&umsg{"Content can not empty", "/"}).ShowMessage(c)
		return
	}
	_, err := DB.Exec("update top_article set title=?, content=? where aid = ? limit 1", BI.Title, BI.Content, BI.Aid)
	if err == nil {
		Cache = lru.New(CacheSize)
		(&umsg{"Success", "/"}).ShowMessage(c)
	} else {
		(&umsg{"Failed to save blog", "/"}).ShowMessage(c)
	}

}
func (ac *AdminController) SaveBlogAddCtr(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		(&umsg{"You have no permission", "/"}).ShowMessage(c)
		return
	}
	var BI BlogItem
	c.BindWith(&BI, binding.Form)
	if BI.Title == "" {
		(&umsg{"Title can not empty", "/"}).ShowMessage(c)
		return
	}
	if BI.Content == "" {
		(&umsg{"Content can not empty", "/"}).ShowMessage(c)
		return
	}
	_, err := DB.Exec(
		"insert into top_article (title, content, publish_time, publish_status) values (?, ?, ?, 1)",
		BI.Title, BI.Content, time.Now().Format("2006-01-02 15:04:05"))
	if err == nil {
		Cache = lru.New(CacheSize)
		(&umsg{"Success", "/"}).ShowMessage(c)
	} else {
		(&umsg{"Failed to save blog", "/"}).ShowMessage(c)
	}

}

func (ac *AdminController) LoginCtr(c *gin.Context) {
	c.HTML(http.StatusOK, "admin-login.html", gin.H{})
}

func (ac *AdminController) LoginProcessCtr(c *gin.Context) {
	var form AdminLoginForm
	c.BindWith(&form, binding.Form)
	session := sessions.Default(c)
	if form.Username == Config.Admin_user && form.Password == Config.Admin_password {
		session.Set("username", "netroby")
		session.Save()
		c.Redirect(301, "/")
	} else {
		session.Delete("username")
		session.Save()
		(&umsg{"Login Failed. You have no permission", "/"}).ShowMessage(c)
	}
}

func (ac *AdminController) LogoutCtr(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("username")
	session.Save()
	c.Redirect(301, "/")
}
