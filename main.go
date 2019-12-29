package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB
var s3Session, _ = session.NewSession()
var env = os.Getenv("ENV")
var (
	// S3Region the S3 region for images
	S3Region = os.Getenv("S3_REGION")
	// S3Bucket the S3 bucket for images
	S3Bucket = os.Getenv("S3_BUCKET")
	// S3Endpoint the S3 endpoint for localdev
	S3Endpoint = os.Getenv("S3_ENDPOINT")
)

var (
	dbName = os.Getenv("DB_NAME")
	dbPass = os.Getenv("DB_PASS")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
)

func connectToDb() {
	log.Println("Connecting to the databse")
	dbSource := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", dbPass, dbHost, dbPort, dbName)
	log.Println("Database source:", dbSource)

	var err error
	db, err = gorm.Open("mysql", dbSource)

	if err != nil {
		panic("Failed to connect to the database " + err.Error())
	}
	log.Println("Connection to the database established")
}

func migrateDb() {
	log.Println("Migrating the database to match model")
	db.AutoMigrate(&Note{}).AddUniqueIndex("idx_note_title_user", "title", "user_id")
}

func initialiseDb() {
	connectToDb()
	migrateDb()
	db.Debug()
}

func initialiseS3() {
	log.Printf("Initialising AWS")
	var err error
	AWSEndPointResolver := func(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
		if service == endpoints.S3ServiceID {
			return endpoints.ResolvedEndpoint{
				URL:           S3Endpoint,
				SigningRegion: S3Region,
			}, nil
		}

		return endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	}

	if S3Endpoint != "" {
		log.Printf("Setting up localstack endpoint")
		// TODO this doesn't work and I've wasted enough time here already
		s3Session, err = session.NewSession(&aws.Config{
			Region:           aws.String(S3Region),
			EndpointResolver: endpoints.ResolverFunc(AWSEndPointResolver),
			Endpoint:         &S3Endpoint,
		})
	} else {
		s3Session, err = session.NewSession(&aws.Config{
			Region: aws.String(S3Region),
		})
	}

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	initialiseDb()
	initialiseS3()
	startServer()
}
