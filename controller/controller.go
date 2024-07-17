package controller

import (
	"context"
	"fmt"
	"io"
	"strings"

	model "service_go_fetch_device_tenant/model"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/golang-jwt/jwt/v4"
)

const COLLECTION_DEVICE_TENANT = "TenantDeviceProfile"

func getFileFromS3(bucket, key string, region string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := s3.NewFromConfig(cfg)

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := client.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		return "", fmt.Errorf("failed to get file from S3, %v", err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read file content, %v", err)
	}

	return string(body), nil
}

func ValidateToken(tokens string) (int, string, string, error) {
	// fmt.Println("in ValidateToken")
	var REGION = "ap-southeast-1"
	var BUCKET = "cdk-hnb659fds-assets-058264531773-ap-southeast-1"
	var KEYFILE = "token.txt"
	setKey, err := getFileFromS3(BUCKET, KEYFILE, REGION)
	jwtKey := []byte(setKey)
	if err != nil {
		return 500, "Internal server error", "Internal server error", err
	}
	tokenString := strings.TrimPrefix(tokens, "Bearer ")
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		// fmt.Println("err ====> ", err)
		if err == jwt.ErrSignatureInvalid {
			return 401, "unauthorized", "unauthorized", err
		}
		return 401, "unauthorized", "unauthorized", err
	}

	if !token.Valid {
		return 401, "unauthorized", "unauthorized", err
	}

	return 200, claims.Data.Tenan, claims.Data.Type, nil
}

func CheckSuperAdmin(tenan string, userType string) bool {
	var SUPERADMIN = "super_admin"
	if tenan == SUPERADMIN && userType == SUPERADMIN {
		return true
	}
	return false
}

func QueryTenantData(isSuperAdmin bool, tenant string) ([]model.ModelTenantDevices, error) {

	var deviceTenant []model.ModelTenantDevices

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	if isSuperAdmin {
		params := &dynamodb.ScanInput{
			TableName: aws.String(COLLECTION_DEVICE_TENANT),
		}
		result, err := svc.Scan(params)
		for _, item := range result.Items {
			el := model.ModelTenantDevices{}
			err = dynamodbattribute.UnmarshalMap(item, &el)
			if err != nil {
				return deviceTenant, err
			}

			var setData = model.ModelTenantDevices{
				TenantDeviceID: el.TenantDeviceID,
				CreateDate:     el.CreateDate,
				DeviceID:       el.DeviceID,
				DeviceType:     el.DeviceType,
				Solution:       el.Solution,
				TenantID:       el.TenantID,
			}

			deviceTenant = append(deviceTenant, setData)
		}
	} else {
		params := &dynamodb.ScanInput{
			TableName:        aws.String(COLLECTION_DEVICE_TENANT),
			FilterExpression: aws.String("#tenantID = :tenantIDVal"),
			ExpressionAttributeNames: map[string]*string{
				"#tenantID": aws.String("tenantID"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":tenantIDVal": {
					S: aws.String(tenant),
				},
			},
		}

		result, err := svc.Scan(params)
		for _, item := range result.Items {
			el := model.ModelTenantDevices{}
			err = dynamodbattribute.UnmarshalMap(item, &el)
			if err != nil {
				return deviceTenant, err
			}

			var setData = model.ModelTenantDevices{
				TenantDeviceID: el.TenantDeviceID,
				CreateDate:     el.CreateDate,
				DeviceID:       el.DeviceID,
				DeviceType:     el.DeviceType,
				Solution:       el.Solution,
				TenantID:       el.TenantID,
			}

			deviceTenant = append(deviceTenant, setData)
		}
	}

	return deviceTenant, nil
}

func HaddleFetchData(tenan string, userType string) ([]model.ModelTenantDevices, error) {
	var dataOut []model.ModelTenantDevices
	isSuperAdmin := CheckSuperAdmin(tenan, userType)
	deviceTenantOut, err := QueryTenantData(isSuperAdmin, tenan)
	if err != nil {
		return dataOut, err
	}
	return deviceTenantOut, nil
}
