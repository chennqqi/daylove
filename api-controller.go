package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/groupcache/lru"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"path/filepath"
)

type BlogItemFull struct {
	Aid         string
	Content     string
	Images      []string
	PublishTime string
}

// APILoginForm is the login form for Admin
type APILoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type APIController struct {
	Token string
}

func Sha512RandomString() string {
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha512.New()
	h.Write([]byte(s))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

func genNewFileName(n string) string {
	return Sha512RandomString() + filepath.Ext(n)
}
func (ac *APIController) FileUpload(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "" || token != ac.Token {
		c.JSON(http.StatusForbidden, gin.H{"msg": "token not valid"})
		return
	}
	s := awsSession.New(&aws.Config{
		Region: aws.String(Config.ObjectStorage.Aws_region),
		Credentials: credentials.NewStaticCredentials(
			Config.ObjectStorage.Aws_access_key_id,
			Config.ObjectStorage.Aws_secret_access_key,
			"",
		),
	})
	s3o := s3.New(s)

	file, fileHeader, err := c.Request.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg":"uploading error"})
		return
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg":"uploading error"})
		return
	}
	prefix := time.Now().In(loc).Format("2006/01/02")
	body, err := ioutil.ReadAll(file)
	storageTargetPath := fmt.Sprintf("%s/%s", prefix, genNewFileName(fileHeader.Filename))
	params := &s3.PutObjectInput{
		Bucket:      aws.String(Config.ObjectStorage.Aws_bucket),
		Key:         aws.String(storageTargetPath),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(fileHeader.Header.Get("content-type")),
	}
	_, err = s3o.PutObject(params)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg":"Can not upload file"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": Config.ObjectStorage.Cdn_url + storageTargetPath})
	(&umsg{"Upload success", "/admin/files"}).ShowMessage(c)
}
func (ac *APIController) ListCtr(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "" || token != ac.Token {
		c.JSON(http.StatusForbidden, gin.H{"msg": "token not valid"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		fmt.Println(err)
	}
	page -= 1
	if page < 0 {
		page = 0
	}

	rpp := 20
	offset := page * rpp
	CKey := fmt.Sprintf("%s-home-page-%d-rpp-%d", GetMinutes(), page, rpp)
	var blogList []BlogItemFull
	val, ok := Cache.Get(CKey)
	if val != nil && ok == true {
		fmt.Println("Ok, we found cache, Cache Len: ", Cache.Len())
		blogList = val.([]BlogItemFull)
	} else {
		rows, err := DB.Query("Select aid, content, images, publish_time from article where publish_status = 1 order by aid desc limit ? offset ? ", &rpp, &offset)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		var Images string
		for rows.Next() {
			blog := BlogItemFull{}

			err := rows.Scan(&blog.Aid, &blog.Content, &Images, &blog.PublishTime)
			if err != nil {
				fmt.Println(err)
			}
			err = json.Unmarshal([]byte(Images), &blog.Images)
			if err != nil {
				fmt.Println(err)
			}
			if blog.Images == nil {
				blog.Images = make([]string, 0)
			}
			blogList = append(blogList, blog)
		}
		err = rows.Err()
		if err != nil {
			fmt.Println(err)
		}
		go func(CKey string, blogList []BlogItemFull) {
			Cache.Add(CKey, blogList)
		}(CKey, blogList)
	}
	c.JSON(http.StatusOK, gin.H{"data": blogList})
}

func (ac *APIController) SaveBlogEditCtr(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "" || token != ac.Token {
		c.JSON(http.StatusForbidden, gin.H{"msg": "token not valid"})
		return
	}
	var BI EditBlogItem
	c.BindJSON(&BI)
	if BI.Aid == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "can not find the blog to edit"})
		return
	}
	if BI.Content == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Content should not be empty"})
		return
	}
	_, err := DB.Exec("update article set title=?, content=? where aid = ? limit 1", BI.Content, BI.Aid)
	if err == nil {
		Cache = lru.New(CacheSize)
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Error"})
	}

}
func (ac *APIController) SaveBlogAddCtr(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "" || token != ac.Token {
		c.JSON(http.StatusForbidden, gin.H{"msg": "token not valid"})
		return
	}
	var BI BlogItem
	c.BindJSON(&BI)
	if BI.Content == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Content can not empty"})
		return
	}
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Time zone setting error"})
		return
	}
	_, err = DB.Exec(
		"insert into article (content, publish_time, publish_status) values (?, ?, 1)",
		BI.Content, time.Now().In(loc).Format("2006-01-02 15:04:05"))
	if err == nil {
		Cache = lru.New(CacheSize)
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Save failed"})
	}

}

func (ac *APIController) LoginCtr(c *gin.Context) {
	login := struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}{}
	c.BindJSON(&login)
	if login.Username == Config.Admin_user && login.Password == Config.Admin_password {
		ac.Token = Sha512RandomString()
		c.JSON(http.StatusOK, gin.H{"msg": "login success", "token": ac.Token})
	} else {
		ac.Token = ""
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": "Login failed"})
	}
}

func (ac *APIController) LogoutCtr(c *gin.Context) {
	token := c.DefaultQuery("token", "")
	if token == "" || token != ac.Token {
		c.JSON(http.StatusForbidden, gin.H{"msg": "token not valid"})
		return
	}
	ac.Token = ""
	c.JSON(http.StatusOK, gin.H{"msg": "logout success"})
}
