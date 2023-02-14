package main

import (
	"github.com/gin-gonic/gin"
	"social-alarm-service/controller"
	"social-alarm-service/dbHelper"
	"social-alarm-service/repository"
	"social-alarm-service/service"
)

func main() {
	r := gin.Default()

	registerRoutes(r)

	r.Run(":8080")
}

func registerRoutes(r *gin.Engine) {
	db, _ := dbHelper.Connect()

	alarmRepository := repository.NewAlarmRepository(db)
	alarmService := service.NewAlarmService(alarmRepository)
	alarmController := controller.NewAlarmController(alarmService)

	r.GET("/get-eligible-alarms", alarmController.GetPublicNonExpiredAlarms)
}
