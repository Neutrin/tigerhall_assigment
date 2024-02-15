package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/nitin/tigerhall/core/inits"
	"github.com/nitin/tigerhall/core/internal/model"
	repositiories "github.com/nitin/tigerhall/core/internal/repositiories"
)

// TODO : make valid request messages
var validate *validator.Validate

func init() {
	validate = validator.New()
}

type TigerControllers struct {
	repo repositiories.TigerRepo
}

func NewTigerController(repo repositiories.TigerRepo) TigerControllers {
	return TigerControllers{repo: repo}
}

// func NewTigerController()

//	func NewTigerController() TigerControllers {
//		return TigerControllers{
//			validator: validator.New(),
//		}
//	}
//
// TODO add multi inheriting
type TigerSightingReq struct {
	TigerId           string `json:"tiger_id" validate:"number"`
	LastSeenTimeStamp string `json:"last_seen" validate:"datetime=02-01-2006 15:04:05"`
	Lat               string `json:"latitude" validate:"latitude"`
	Long              string `json:"longitude" validate:"longitude"`
}

type CreateTigerReq struct {
	Name              string `json:"name" validate:"alphanum"`
	LastSeenTimeStamp string `json:"last_seen" validate:"datetime=02-01-2006 15:04:05"`
	Lat               string `json:"latitude" validate:"latitude"`
	Long              string `json:"longitude" validate:"longitude"`
	DOB               string `json:"date_of_birth" validate:"datetime=02-01-2006"`
}

func (controller TigerControllers) AddTigerSighting(c *gin.Context) {
	var (
		sightings model.TigerSightings
		err       error
		u         *url.URL
	)
	if sightings, err = fetchSightData(c); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	f, fileUpload, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()})
		return
	}

	log.Println(" ******* file path name becomes = *******  ", fileUpload.Filename)
	defer f.Close()
	u, err = inits.UploadFile(f, fileUpload)
	//log.Printf(" here the values becomes = %s\n", c.PostForm("tiger_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": u.EscapedPath(),
	})
	_ = sightings
	// log.Printf("%+v\n", sightings)
	// c.JSON(200, gin.H{"success": "user logged in"})
}

func (controller TigerControllers) AddTiger(c *gin.Context) {
	var (
		tiger model.Tiger
		err   error
		u     *url.URL
	)
	if tiger, err = fetchTigerData(c); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	f, fileUpload, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()})
		return
	}
	defer f.Close()
	u, err = inits.UploadFile(f, fileUpload)
	//log.Printf(" here the values becomes = %s\n", c.PostForm("tiger_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}
	//Transaction will start here
	_, err = controller.repo.CreateTiger(tiger, u.EscapedPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}
	// log.Println(" tiuger id becom,es -= ", tigerId)
	// sightings := model.TigerSightings{
	// 	TigerId:           tigerId,
	// 	LastSeenTimeStamp: tiger.LastSeenTimeStamp,
	// 	Lat:               tiger.Lat,
	// 	Long:              tiger.Long,
	// 	ImagePath:         u.EscapedPath(),
	// }
	//controller.repo.CreateTigerSighting(sightings)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("tiger = %s added ", tiger.Name),
	})

}

func fetchSightData(c *gin.Context) (model.TigerSightings, error) {
	var (
		tigerSighting model.TigerSightings
		err           error
	)
	sightingReq := TigerSightingReq{
		TigerId:           c.PostForm("tiger_id"),
		LastSeenTimeStamp: c.PostForm("last_seen"),
		Lat:               c.PostForm("latitude"),
		Long:              c.PostForm("longitude"),
	}

	if err = validate.Struct(sightingReq); err == nil {
		tigerSighting.TigerId, _ = strconv.Atoi(sightingReq.TigerId)
		lastSeemTime, _ := time.Parse("02-01-2006 15:04:05", sightingReq.LastSeenTimeStamp)
		log.Println(" last seen time stamp is = ", lastSeemTime)

		tigerSighting.LastSeenTimeStamp = lastSeemTime.UnixMilli()
		log.Println(" and after convetion time stamp becomes = ", time.UnixMilli(tigerSighting.LastSeenTimeStamp))
		tigerSighting.Lat, _ = strconv.ParseFloat(sightingReq.Lat, 64)
		tigerSighting.Long, _ = strconv.ParseFloat(sightingReq.Long, 64)

	}
	return tigerSighting, err

}

func fetchTigerData(c *gin.Context) (model.Tiger, error) {
	var (
		tiger model.Tiger
		err   error
	)
	tigerReq := CreateTigerReq{
		LastSeenTimeStamp: c.PostForm("last_seen"),
		Lat:               c.PostForm("latitude"),
		Long:              c.PostForm("longitude"),
		Name:              c.PostForm("name"),
		DOB:               c.PostForm("date_of_birth"),
	}
	if err = validate.Struct(tigerReq); err == nil {
		lastSeenTime, _ := time.Parse("15-01-2006 15:04:05", tigerReq.LastSeenTimeStamp)
		log.Println(" last seen time stamp is = ", lastSeenTime)
		tiger.LastSeenTimeStamp = lastSeenTime.UnixMilli()
		log.Println(" and the last seen tiome stamp is = ", tiger.LastSeenTimeStamp)
		log.Println(" and after convetion time stamp becomes = ", time.UnixMilli(tiger.LastSeenTimeStamp).UTC())
		tiger.Lat, _ = strconv.ParseFloat(tigerReq.Lat, 64)
		tiger.Long, _ = strconv.ParseFloat(tigerReq.Long, 64)
		tiger.Name = tigerReq.Name
		tiger.DOB, _ = time.Parse(tigerReq.DOB, "02-01-2006")
	}

	return tiger, err
}
