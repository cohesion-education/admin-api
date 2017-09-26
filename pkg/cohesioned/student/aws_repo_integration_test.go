//+build integration

package student_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cohesion-education/api/fakes"
	"github.com/cohesion-education/api/pkg/cohesioned/config"
	"github.com/cohesion-education/api/pkg/cohesioned/student"
	"github.com/cohesion-education/api/testutils"
)

var (
	repo      student.Repo
	emptyTime = time.Time{}
)

func TestMain(m *testing.M) {
	awsConfig, err := config.NewAwsConfig()
	if err != nil {
		panic("Unexpected error initializing AwsConfig: " + err.Error())
	}

	db, err := awsConfig.DialRDS()
	if err != nil {
		panic("Failed to connect to db " + err.Error())
	}

	repo = student.NewAwsRepo(db)

	if err := testutils.SetupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	testResult := m.Run()

	if err := testutils.CleanupDB(db); err != nil {
		fmt.Println(err.Error())
	}

	os.Exit(testResult)
}

func TestRepoList(t *testing.T) {
	parentID := int64(1)
	list, err := repo.List(parentID)
	if err != nil {
		t.Errorf("Failed to list student: %v", err)
	}

	if len(list) == 0 {
		t.Errorf("student list is empty")
	}

	for _, student := range list {
		if len(student.Name) == 0 {
			t.Error("student name is empty")
		}

		if len(student.Grade) == 0 {
			t.Error("student grade is empty")
		}

		if len(student.School) == 0 {
			t.Error("student school is empty")
		}

		if student.ParentID != parentID {
			t.Errorf("parent id incorrect; expected %d received %d", parentID, student.ParentID)
		}

		if student.Created == emptyTime {
			t.Error("student created is empty")
		}

		if student.CreatedBy == 0 {
			t.Error("student created by is empty")
		}
	}
}

func TestRepoSave(t *testing.T) {
	student := fakes.FakeStudent()
	id, err := repo.Save(student)

	if err != nil {
		t.Errorf("Failed to save student: %v", err)
	}

	if id == 0 {
		t.Errorf("student ID was zero - expected db to generate a student id")
	}
}

func TestRepoUpdate(t *testing.T) {
	student := fakes.FakeStudent()
	student.Updated = time.Now()
	student.UpdatedBy = student.CreatedBy
	err := repo.Update(student)

	if err != nil {
		t.Errorf("Failed to update student: %v", err)
	}
}
