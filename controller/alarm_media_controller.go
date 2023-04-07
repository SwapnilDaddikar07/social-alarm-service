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
	"social-alarm-service/constants"
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
	fmt.Println("validating request bindings")
	request := request_model.GetMediaForAlarm{}

	bindingErr := ctx.ShouldBindWith(&request, binding.JSON)
	if bindingErr != nil {
		fmt.Printf("request binding failed %v", bindingErr)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	fmt.Printf("request binding successful. Calling service to fetch media associated with alarm id %s\n", request.AlarmId)
	alarmMedia, serviceErr := amc.service.GetMediaForAlarm(ctx, request.AlarmId, request.UserId)
	if serviceErr != nil {
		fmt.Printf("service error %v", serviceErr)
		ctx.AbortWithStatusJSON(serviceErr.HttpStatusCode, serviceErr)
		return
	}

	fmt.Printf("fetched %d associated media for alarm id %s\n", len(alarmMedia), request.AlarmId)
	ctx.JSON(http.StatusOK, alarmMedia)
}

func (amc alarmMediaController) UploadMedia(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, constants.MaxFileSizeInMB<<20)

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

	fmt.Println("request validation completed. creating tmp file.")
	tmpFileName, tmpFileError := amc.service.CreateTmpFile(ctx, file, filepath.Ext(fileHeader.Filename))
	if tmpFileError != nil {
		fmt.Printf("error creating tmp file %v", tmpFileError)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	defer amc.service.DeleteTmpFile(ctx, tmpFileName)

	serviceError := amc.service.UploadMedia(ctx, alarmId, senderId, tmpFileName)
	if serviceError != nil {
		fmt.Printf("service error when uploading media for alarm %v", serviceError)
		ctx.AbortWithStatusJSON(serviceError.HttpStatusCode, serviceError)
		return
	}

	fmt.Printf("successfully uploaded media for alarm id %s", alarmId)
	ctx.AbortWithStatus(http.StatusCreated)
}

func isValidContentType(file multipart.File) bool {
	buffIO := bufio.NewReader(file)
	sniffBytes, _ := buffIO.Peek(512)
	contentType := http.DetectContentType(sniffBytes)

	fmt.Printf("content type for file is %s \n", contentType)

	file.Seek(io.SeekStart, io.SeekStart)

	//TODO add relevant content types
	return contentType != "video/mp4" || contentType != "audio/wave"
}
