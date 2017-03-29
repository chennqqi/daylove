package main

import (
	"database/sql"
	//	"github.com/fvbock/endless"
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
	ac := new(AdminController)
	admin := r.Group("/admin")
	{
		admin.GET("/login", ac.LoginCtr)
		admin.POST("/login-process", ac.LoginProcessCtr)
		admin.GET("/logout", ac.LogoutCtr)
		admin.GET("/addblog", ac.AddBlogCtr)
		admin.POST("/save-blog-add", ac.SaveBlogAddCtr)
		admin.GET("/deleteblog/:id", ac.DeleteBlogCtr)
		admin.POST("/save-blog-edit", ac.SaveBlogEditCtr)
		admin.GET("/editblog/:id", ac.EditBlogCtr)
	}

	apiCtrl := new(APIController)
	api := r.Group("/api")
	{
		api.HEAD("/", apiCtrl.HomeCtr)
		api.POST("/list", apiCtrl.ListCtr)
		api.POST("/login", apiCtrl.LoginCtr)
		api.POST("/logout", apiCtrl.LogoutCtr)
		api.POST("/file-upload", apiCtrl.FileUpload)
		api.POST("/save-blog-add", apiCtrl.SaveBlogAddCtr)
		api.POST("/save-blog-edit", apiCtrl.SaveBlogEditCtr)
	}
	r.Run(":8080")
	//endless.ListenAndServe(":8080", r)
}
