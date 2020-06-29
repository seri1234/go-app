package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//Renderer interface
type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

//MySQLの接続設定
func sqlConnect() (database *gorm.DB, err error) {
	DBMS := "mysql"
	USER := "user"
	PASS := "mysql"
	PROTOCOL := "tcp(db:3306)"
	DBNAME := "mydb"

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?parseTime=true"
	return gorm.Open(DBMS, CONNECT)
}

func main() {
	//DB接続
	db, err := sqlConnect()
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("DB接続成功")
	}

	defer db.Close()

	//portの指定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Echo World!!")
	})

	e.GET("/index", func(c echo.Context) error {
		// Postsテーブルデータの取得
		postdata := []Posts{}
		db.Find(&postdata)

		//viewへのデータ受け渡し
		return c.Render(http.StatusOK, "index", postdata)
	})

	e.POST("/post", func(c echo.Context) error {
		//フォームデータの取得
		postdata := Posts{
			Title:     c.FormValue("post"),
			Content:   "posttest",
			UpdatedAt: getDate(),
			CreatedAt: getDate(),
		}
		//Postsテーブルデータの追加
		db.Create(&postdata)

		return c.Render(http.StatusOK, "hello", postdata)
	})

	e.Logger.Fatal(e.Start(":" + port))
}

//Postsテーブル
type Posts struct {
	ID        int
	Title     string
	Content   string
	UpdatedAt string
	CreatedAt string
}

//日付データの取得
func getDate() string {
	const layout = "2006-01-02 15:04:05"
	now := time.Now()
	return now.Format(layout)
}
