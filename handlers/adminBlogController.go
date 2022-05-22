package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"project/tech-blog-go/models"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type BlogPayload struct {
	ID          int                  `json:"id"`
	UserID      int                  `json:"user_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	EyeCatch    string               `json:"eye_catch"`
	Body        string               `json:"body"`
	State       int                  `json:"state"`
	PublishAt   string               `json:"publish_at"`
	Category    []int                `json:"categories"`
	Tag         []int                `json:"tags"`
	Content     []models.BlogContent `json:"contents"`
}

type BlogCategoryPayload struct {
	CategoryID int `json:"category_id"`
}

type BlogEyeCatchPayload struct {
	Buf *bytes.Buffer
}

// all
func (app *Application) getAllBlogs(w http.ResponseWriter, r *http.Request) {
	ctg, err := app.Models.DB.BlogGetAll()
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, ctg, "blogs")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

// one
func (app *Application) getOneBlog(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.Logger.Print(errors.New("invalid id parameter"))
		app.ErrorJSON(w, err)
		return
	}
	b, err := app.Models.DB.GetBlog(id)
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, b, "blog")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

// create or update
func (app *Application) editBlog(w http.ResponseWriter, r *http.Request) {
	var bp BlogPayload
	err := json.NewDecoder(r.Body).Decode(&bp)
	if err != nil {
		log.Println(err)
		app.ErrorJSON(w, err)
		return
	}

	var bm models.Blog

	if bp.ID != 0 {
		id := bp.ID
		b, _ := app.Models.DB.GetBlog(id)
		bm = *b
		bm.UpdatedAt = time.Now()
	}

	//　カテゴリ
	for i := 0; i < len(bp.Category); i++ {
		var b models.BlogCategory
		b.CategoryID = bp.Category[i]
		bm.Category = append(bm.Category, b)
	}

	// タグ
	for i := 0; i < len(bp.Tag); i++ {
		var t models.BlogTag
		t.TagID = bp.Tag[i]
		bm.Tag = append(bm.Tag, t)
	}

	// 目次
	for i := 0; i < len(bp.Content); i++ {
		var c models.BlogContent
		c.Name = bp.Content[i].Name
		c.Anchor = bp.Content[i].Anchor
		c.Position = bp.Content[i].Position
		bm.Content = append(bm.Content, c)
	}

	bm.ID = bp.ID
	bm.UserID = bp.UserID
	bm.Title = bp.Title
	bm.Description = bp.Description
	bm.EyeCatch = bp.EyeCatch
	bm.Body = bp.Body
	bm.State = bp.State
	bm.PublishAt = bp.PublishAt
	bm.UpdatedAt = time.Now()

	if bp.ID == 0 {
		err = app.Models.DB.BlogCreate(bm)
		if err != nil {
			app.ErrorJSON(w, err)
			return
		}
	}
	// else {
	// 	err = app.Models.DB.BlogUpdate(bm)
	// 	if err != nil {
	// 		app.ErrorJSON(w, err)
	// 		return
	// 	}
	// }
	ok := JsonResp{
		OK: true,
	}

	err = app.WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}

func (app *Application) eyeCatchUpload(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load(fmt.Sprintf("./%s.env", os.Getenv("GO_ENV")))
	if err != nil {
		log.Fatalln(err)
	}
	ak := os.Getenv("AWS_ACCESS_KEY")         // .envのAWS_ACCESS_KEY
	ask := os.Getenv("AWS_SECRET_ACCESS_KEY") // .envのAWS_SECRET_ACCESS_KEYx
	ab := os.Getenv("AWS_BUCKET")             // .envのAWS_SECRET_ACCESS_KEYx
	creds := credentials.NewStaticCredentials(ak, ask, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-northeast-1"),
	})
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	file, reader, err := r.FormFile("image")
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	defer file.Close()

	// img := &BlogEyeCatchPayload{
	// 	Buf: &bytes.Buffer{},
	// }

	// , err = img.Buf.ReadFrom(file)
	// if err != nil {
	// 	app.ErrorJSON(w, err)
	// 	return
	// }

	// img, t, err := image.Decode()
	t := time.Now().Format("20060102030405")
	fileAry := strings.Split(reader.Filename, ".")
	pos := len(fileAry) - 1
	fileAry = append(fileAry[:pos+1], fileAry[pos:]...)
	fileAry[pos] = t
	newFile := strings.Join(fileAry, ".")

	uploader := s3manager.NewUploader(sess)
	uploadData, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(ab),
		Key:    aws.String("eyecatch/" + newFile),
		Body:   file,
	})
	if err != nil {
		app.ErrorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, uploadData.Location, "image")

	if err != nil {
		app.ErrorJSON(w, err)
		return
	}
}
