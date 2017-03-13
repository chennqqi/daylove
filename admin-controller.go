package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/groupcache/lru"
)

// AdminLoginForm is the login form for Admin
type AdminLoginForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// BlogItem is the blog item
type BlogItem struct {
	Content string `form:"content" json:"content" binding:"required"`
	Images  string `form:"images" json:"images" `
}
type EditBlogItem struct {
	Aid     string `form:"aid" json:"aid" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

type AdminController struct {
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
		rows, err := DB.Query("select aid, content, publish_time, publish_statusfrom top_article where aid = ?", &id)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		var ()
		for rows.Next() {
			err := rows.Scan(&blog.aid, &blog.content, &blog.publish_time, &blog.publish_status)
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
		"content":          template.HTML(blog.content.String),
		"publish_time":     blog.publish_time.String,
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
	if BI.Content == "" {
		(&umsg{"Content can not empty", "/"}).ShowMessage(c)
		return
	}
	_, err := DB.Exec("update article set title=?, content=? where aid = ? limit 1", BI.Content, BI.Aid)
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
	if BI.Content == "" {
		(&umsg{"Content can not empty", "/"}).ShowMessage(c)
		return
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = DB.Exec(
		"insert into article (content, images, publish_time, publish_status) values (?, '[]', ?, 1)",
		BI.Content, time.Now().In(loc).Format("2006-01-02 15:04:05"))
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
		session.Set("username", Config.Admin_user)
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
