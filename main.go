package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

	transactionManager := transaction_manager.NewTransactionManager(db)

	alarmRepository := repository.NewAlarmRepository(db)
	alarmService := service.NewAlarmService(alarmRepository, transactionManager)
	alarmController := controller.NewAlarmController(alarmService)

	alarmMediaRepository := repository.NewAlarmMediaRepository(db)
	alarmMediaService := service.NewAlarmMediaService(alarmMediaRepository)
	alarmMediaController := controller.NewAlarmMediaController(alarmMediaService)

	r.POST("/create/alarm", alarmController.CreateAlarm)
	r.GET("/eligible/alarms", alarmController.GetPublicNonExpiredAlarms)

	r.GET("/media/alarm", alarmMediaController.GetMediaForAlarm)
}
