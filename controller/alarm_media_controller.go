package controller

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	error2 "social-alarm-service/error"
	"social-alarm-service/request_model"
	"social-alarm-service/service"
	"strings"
)

type AlarmMediaController interface {
	GetMediaForAlarm(ctx *gin.Context)
	UploadMedia(ctx *gin.Context)
}

type alarmMediaController struct {
	service service.AlarmMediaService
}

func NewAlarmMediaController(service service.AlarmMediaService) AlarmMediaController {
	return alarmMediaController{service: service}
}

func (amc alarmMediaController) GetMediaForAlarm(ctx *gin.Context) {
	request := request_model.GetMediaForAlarm{}

	bindingErr := ctx.ShouldBindWith(&request, binding.JSON)
	if bindingErr != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	alarmMedia, serviceErr := amc.service.GetMediaForAlarm(ctx, request.AlarmId)
	if serviceErr != nil {
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	ctx.JSON(200, alarmMedia)
}

func (amc alarmMediaController) UploadMedia(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 14<<20)

	alarmId := ctx.Request.FormValue("alarm_id")
	if strings.TrimSpace(alarmId) == "" {
		fmt.Println("alarm id is missing")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	senderId := ctx.Request.FormValue("sender_id")
	if strings.TrimSpace(senderId) == "" {
		fmt.Println("sender id is missing")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	file, fileHeader, err := ctx.Request.FormFile("media_file")
	if err != nil {
		fmt.Printf("error validating media file %v ", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !isValidContentType(file) {
		fmt.Println("content type is invalid")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, error2.ContentTypeNotSupported)
		return
	}

	tmpFileName, tmpFileError := amc.service.CreateTmpFile(ctx, file, filepath.Ext(fileHeader.Filename))
	if tmpFileError != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	defer amc.service.DeleteTmpFile(ctx, tmpFileName)

	serviceError := amc.service.UploadMedia(ctx, alarmId, senderId, tmpFileName)
	if serviceError != nil {
		ctx.AbortWithStatusJSON(serviceError.HttpStatusCode, serviceError)
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}

func isValidContentType(file multipart.File) bool {
	buffIO := bufio.NewReader(file)
	sniffBytes, _ := buffIO.Peek(512)
	contentType := http.DetectContentType(sniffBytes)

	//TODO check this usage
	file.Seek(io.SeekStart, io.SeekStart)

	//TODO add relevant content types
	return contentType != "video/mp4" || contentType != "audio/wave"
}
