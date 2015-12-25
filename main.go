package main

import (
	"database/sql"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/groupcache/lru"
)

var (
	Config    *appConfig
	DB        *sql.DB
	Cache     *lru.Cache
	CacheSize int = 8192
)

func main() {

	Config = GetConfig()
	DB = GetDB(Config)
	Cache = lru.New(CacheSize)

	r := gin.Default()
	r.Static("/assets", "assets")
	store := sessions.NewCookieStore([]byte("gssecret"))
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("templates/*.html")

	fc := new(FrontController)
	r.GET("/", fc.HomeCtr)
	r.GET("/about", fc.AboutCtr)
	r.GET("/view/:id", fc.ViewCtr)
	r.GET("/view.php", fc.ViewAltCtr)
	r.GET("/ping", fc.PingCtr)
	r.GET("/search", fc.SearchCtr)
	r.GET("/countview/:id", fc.CountViewCtr)

	ac := new(AdminController)
	admin := r.Group("/admin")
	{
		admin.GET("/", ac.ListBlogCtr)
		admin.GET("/login", ac.LoginCtr)
		admin.POST("/login-process", ac.LoginProcessCtr)
		admin.GET("/logout", ac.LogoutCtr)
		admin.GET("/addblog", ac.AddBlogCtr)
		admin.POST("/save-blog-add", ac.SaveBlogAddCtr)
		admin.GET("/listblog", ac.ListBlogCtr)
		admin.GET("/deleteblog/:id", ac.DeleteBlogCtr)
		admin.POST("/save-blog-edit", ac.SaveBlogEditCtr)
		admin.GET("/editblog/:id", ac.EditBlogCtr)
	}

	a := new(api)
	api := r.Group("/api")
	{
		api.GET("/", a.index)
		api.GET("view/:id", a.view)
	}
	rss := new(RSS)
	r.GET("/rss.php", rss.Alter)
	r.GET("/rss", rss.Out)
	endless.ListenAndServe(":8080", r)
}
