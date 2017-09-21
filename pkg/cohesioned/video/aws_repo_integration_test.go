// +build integration

package video_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/video"
	"github.com/joho/godotenv"
)

var (
	repo        video.Repo
	testVideoID int64
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		panic("Failed to load .env file: " + err.Error())
	}

	awsConfig, err := config.NewAwsConfig()
	if err != nil {
		panic("Unexpected error initializing AwsConfig: " + err.Error())
	}

	db, _ := awsConfig.DialRDS()
	repo, err = video.NewAwsRepo(awsConfig, "test-bucket")
	if err != nil {
		panic("Failed to connect to db " + err.Error())
	}

	stmt, _ := db.Prepare("insert into user (email, full_name, first_name, last_name) values ('test@cohesioned.io', 'Test User', 'Test', 'User')")
	result, err := stmt.Exec()
	if err != nil {
		panic("Failed to insert user " + err.Error())
	}

	userID, _ := result.LastInsertId()

	stmt, _ = db.Prepare("insert into audit_info (created, created_by) values (now(), ?)")
	result, err = stmt.Exec(userID)
	if err != nil {
		panic("Failed to insert audit_info " + err.Error())
	}

	auditInfoID, _ := result.LastInsertId()

	stmt, _ = db.Prepare("insert into taxonomy (name, audit_info_id) values ('test-taxonomy', ?)")
	result, err = stmt.Exec(auditInfoID)
	if err != nil {
		panic("Failed to insert taxonomy " + err.Error())
	}

	taxonomyID, _ := result.LastInsertId()

	stmt, _ = db.Prepare("insert into video (title, file_name, bucket, object_key, taxonomy_id, audit_info_id) values ('test-video', 'test-file.fake', 'test-bucket', 'test-obj-key', ?, ?)")

	result, err = stmt.Exec(taxonomyID, auditInfoID)
	if err != nil {
		panic("Failed to insert video " + err.Error())
	}

	testVideoID, _ = result.LastInsertId()
	fmt.Printf("inserted video ID: %d\n", testVideoID)

	testResult := m.Run()

	cleanupSql := `
		delete from video;
		delete from taxonomy;
		delete from audit_info;
		delete from user;
	`

	for _, sql := range strings.Split(cleanupSql, ";") {
		trimmedSql := strings.TrimSpace(sql)
		if len(trimmedSql) == 0 {
			continue
		}

		stmt, err = db.Prepare(trimmedSql)
		if err != nil {
			panic("Failed to prepare db cleanup statement " + err.Error())
		}

		if _, err := stmt.Exec(); err != nil {
			fmt.Errorf("Failed to cleanup %v\n", err)
		}
	}

	os.Exit(testResult)
}

func TestGet(t *testing.T) {

	video, err := repo.Get(testVideoID)
	if err != nil {
		t.Errorf("Failed to get video by ID: %v", err)
	}

	if video == nil {
		t.Error("Video by ID was null")
	}
	// db, err := config.DialRDS()
	// if err != nil {
	// 	t.Errorf("Failed to dial RDS: %v", err)
	// }
	//
	// tx, err := db.Begin()
	// if err != nil {
	// 	t.Errorf("Failed to start transaction: %v", err)
	// }
	//
	// defer tx.Rollback()
	//
	// repo.Add(video)
}

func TestList(t *testing.T) {
	list, err := repo.List()
	if err != nil {
		t.Errorf("Failed to List videos: %v", err)
	}

	fmt.Printf("video list: %v\n", list)
}
