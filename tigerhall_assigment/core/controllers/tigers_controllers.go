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
	"github.com/umahmood/haversine"

	"github.com/nitin/tigerhall/core/inits"
	"github.com/nitin/tigerhall/core/internal/config"
	"github.com/nitin/tigerhall/core/internal/model"
	repositiories "github.com/nitin/tigerhall/core/internal/repositiories"
	"github.com/nitin/tigerhall/core/models"
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

func (controller TigerControllers) AddTigerSighting(c *gin.Context) {
	var (
		sightings model.TigerSightings
		err       error
		u         *url.URL
		tiger     model.Tiger
	)
	if sightings, err = fetchSightData(c); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tiger, err = controller.repo.TigerById(sightings.TigerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()})
		return
	}

	err = vSightingRule(sightings, tiger)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()})
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

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}
	sightings.ImagePath = u.EscapedPath()
	_, err = controller.repo.CreateTigerSighting(sightings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "sighting added ",
	})
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	_, err = controller.repo.CreateTiger(tiger, u.EscapedPath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("tiger = %s added ", tiger.Name),
	})

}

func (controller TigerControllers) ListAllTigers(c *gin.Context) {
	var (
		pageNo, _ = strconv.Atoi(c.Query("page_no"))
		limit, _  = strconv.Atoi(c.Query("limit"))
		//TODO : Do something of this
		sortParam = "last_seen"
		responses []models.TigerResp
	)
	pagParams := repositiories.NewPagination(limit, pageNo, sortParam)
	result, err := controller.repo.ListAllTigers(pagParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}
	for _, curTiger := range result.Rows.([]*model.Tiger) {
		responses = append(responses, models.TigerResp{
			Name:     curTiger.Name,
			Lat:      fmt.Sprintf("%f", curTiger.Lat),
			Long:     fmt.Sprintf("%f", curTiger.Long),
			LastSeen: time.UnixMilli(curTiger.LastSeenTimeStamp).UTC().Format(config.DateTimeFormat),
		})

	}
	result.Rows = responses
	c.JSON(http.StatusOK, result)

}

func fetchSightData(c *gin.Context) (model.TigerSightings, error) {
	var (
		tigerSighting model.TigerSightings
		err           error
	)
	sightingReq := models.TigerSightingReq{
		TigerId:           c.PostForm("tiger_id"),
		LastSeenTimeStamp: c.PostForm("last_seen"),
		Lat:               c.PostForm("latitude"),
		Long:              c.PostForm("longitude"),
	}

	if err = validate.Struct(sightingReq); err == nil {
		tigerSighting.TigerId, _ = strconv.Atoi(sightingReq.TigerId)
		lastSeemTime, _ := time.Parse(config.DateTimeFormat, sightingReq.LastSeenTimeStamp)
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
	tigerReq := models.CreateTigerReq{
		LastSeenTimeStamp: c.PostForm("last_seen"),
		Lat:               c.PostForm("latitude"),
		Long:              c.PostForm("longitude"),
		Name:              c.PostForm("name"),
		DOB:               c.PostForm("date_of_birth"),
	}
	if err = validate.Struct(tigerReq); err == nil {
		lastSeenTime, _ := time.Parse(config.DateTimeFormat, tigerReq.LastSeenTimeStamp)
		log.Println(" last seen time stamp is = ", lastSeenTime)
		tiger.LastSeenTimeStamp = lastSeenTime.UnixMilli()
		log.Println(" and the last seen tiome stamp is = ", tiger.LastSeenTimeStamp)
		log.Println(" and after convetion time stamp becomes = ", time.UnixMilli(tiger.LastSeenTimeStamp).UTC())
		tiger.Lat, _ = strconv.ParseFloat(tigerReq.Lat, 64)
		tiger.Long, _ = strconv.ParseFloat(tigerReq.Long, 64)
		tiger.Name = tigerReq.Name
		tiger.DOB, _ = time.Parse(tigerReq.DOB, config.DateFormat)
	}

	return tiger, err
}

func vSightingRule(curSight model.TigerSightings, tiger model.Tiger) error {
	var err error
	if curSight.LastSeenTimeStamp < tiger.LastSeenTimeStamp {
		err = fmt.Errorf(" this is an older sighting")
		return err
	}
	curCordinate := haversine.Coord{
		Lat: curSight.Lat,
		Lon: curSight.Long,
	}
	lastCordinate := haversine.Coord{
		Lat: tiger.Lat,
		Lon: tiger.Long,
	}

	_, distInKm := haversine.Distance(curCordinate, lastCordinate)
	if distInKm < config.TigerSightDistInKm {
		err = fmt.Errorf(" tiger is not withing %f km range", float64(5))
	}
	return err

}
