package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"social-alarm-service/controller"
	"social-alarm-service/db_helper"
	"social-alarm-service/repository"
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

	alarmRepository := repository.NewAlarmRepository(db)
	alarmService := service.NewAlarmService(alarmRepository)
	alarmController := controller.NewAlarmController(alarmService)

	r.GET("/get-eligible-alarms", alarmController.GetPublicNonExpiredAlarms)
}
