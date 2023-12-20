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
	utils2 "social-alarm-service/utils"
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
	presignClient := s3.NewPresignClient(s3Client)
	awsUtil := aws_util.NewAWSUtil(s3Client, presignClient)

	transactionManager := transaction_manager.NewTransactionManager(db)

	utils := utils2.NewUtils()
	userRepo := repository.NewUserRepository(db)
	alarmRepository := repository.NewAlarmRepository(db)
	alarmMediaRepository := repository.NewAlarmMediaRepository(db)

	alarmService := service.NewAlarmService(alarmRepository, userRepo, alarmMediaRepository, transactionManager)
	alarmController := controller.NewAlarmController(alarmService)

	alarmMediaService := service.NewAlarmMediaService(alarmRepository, alarmMediaRepository, userRepo, awsUtil, utils, transactionManager)
	alarmMediaController := controller.NewAlarmMediaController(alarmMediaService)

	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	publicRoute := r.Group("/api/social-alarm-api")

	publicRoute.POST("/create/alarm", alarmController.CreateAlarm)
	publicRoute.POST("/eligible/alarms", alarmController.GetPublicNonExpiredAlarms)
	publicRoute.POST("/update/alarm-status", alarmController.UpdateAlarmStatus)
	publicRoute.POST("/my/alarms", alarmController.GetAllAlarms)
	publicRoute.POST("/delete/alarm", alarmController.Delete)

	publicRoute.POST("/media/alarm", alarmMediaController.GetMediaForAlarm)
	publicRoute.POST("/upload/media", alarmMediaController.UploadMedia)

	publicRoute.POST("/profiles", userController.GetProfiles)

}
