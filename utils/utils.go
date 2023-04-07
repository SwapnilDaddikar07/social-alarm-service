package utils

import uuid2 "github.com/google/uuid"

type Utils interface {
	GenerateUUID() string
}

type utils struct {
}

func NewUtils() Utils {
	return utils{}
}

func (u utils) GenerateUUID() string {
	uuid, _ := uuid2.NewUUID()
	return uuid.String()
}
