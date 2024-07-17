package model

import "github.com/golang-jwt/jwt/v4"

type ModelTenantDevices struct {
	TenantDeviceID string `json:"tenantDeviceID"`
	CreateDate     string `json:"createDate"`
	DeviceID       string `json:"deviceID"`
	DeviceType     string `json:"deviceType"`
	Solution       string `json:"solution"`
	TenantID       string `json:"tenantID"`
}

type DBdata struct {
	AuthStatus bool     `json:"authStatus"`
	Email      string   `json:"email"`
	IsProduct  []string `json:"isProduct"`
	Tenan      string   `json:"tenan"`
	Type       string   `json:"type"`
}

type Claims struct {
	Data DBdata `json:"data"`
	jwt.RegisteredClaims
}
