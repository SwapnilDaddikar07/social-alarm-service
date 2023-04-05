package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"social-alarm-service/aws_util"
	"social-alarm-service/controller"
	"social-alarm-service/db_helper"
	"social-alarm-service/repository"
	"social-alarm-service/repository/transaction_manager"
	"social-alarm-service/service"
)

func main() {
	r := gin.Default()

	registerRoutes(r)

	r.Run(":8080")
}

func registerRoutes(r *gin.Engine) {
	db, err := db_helper.Connect()
	if err != nil {
		fmt.Println(err)
		panic("DB connection error")
	}

	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Couldn't load default configuration. Have you set up your AWS account? %v \n", err)
		return
	}

	s3Client := s3.NewFromConfig(sdkConfig)
	awsUtil := aws_util.NewAWSUtil(s3Client)

	transactionManager := transaction_manager.NewTransactionManager(db)

	alarmRepository := repository.NewAlarmRepository(db)
	alarmService := service.NewAlarmService(alarmRepository, transactionManager)
	alarmController := controller.NewAlarmController(alarmService)

	alarmMediaRepository := repository.NewAlarmMediaRepository(db)
	alarmMediaService := service.NewAlarmMediaService(alarmRepository, alarmMediaRepository, awsUtil, transactionManager)
	alarmMediaController := controller.NewAlarmMediaController(alarmMediaService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	publicRoute := r.Group("/api/social-alarm")

	publicRoute.POST("/create/alarm", alarmController.CreateAlarm)
	publicRoute.GET("/eligible/alarms", alarmController.GetPublicNonExpiredAlarms)
	publicRoute.POST("/update/alarm-status", alarmController.UpdateAlarmStatus)
	publicRoute.GET("/my/alarms", alarmController.GetAllAlarms)

	publicRoute.GET("/media/alarm", alarmMediaController.GetMediaForAlarm)
	publicRoute.POST("/upload/media", alarmMediaController.UploadMedia)

	publicRoute.GET("/profiles", userController.GetProfiles)
}
