//+build integration

package testutils

import (
	"database/sql"
	"fmt"
	"strings"
)

func SetupDB(db *sql.DB) error {
	setupSql := `
		insert into user (created, email, full_name, first_name, last_name) values (now(), 'test@cohesioned.io', 'Test User', 'Test', 'User');
		insert into taxonomy (name, created_by, created) values ('test-taxonomy', 1, now());
		insert into taxonomy (name, created_by, created) values ('test-parent', 1, now());
		insert into taxonomy (name, created_by, created, parent_id) values ('test-child', 1, now(), 2);
		insert into video (title, file_name, bucket, object_key, taxonomy_id, created_by, created) values ('test-video', 'test-file.fake', 'test-bucket', 'test-obj-key', 1, 1, now());
	`

	for _, sql := range strings.Split(setupSql, ";") {
		trimmedSql := strings.TrimSpace(sql)
		if len(trimmedSql) == 0 {
			continue
		}

		stmt, err := db.Prepare(trimmedSql)
		if err != nil {
			return fmt.Errorf("Failed to prepare statement %s: %v", trimmedSql, err)
		}

		if _, err := stmt.Exec(); err != nil {
			return fmt.Errorf("Failed to setup db %v", err)
		}
	}

	return nil
}

func CleanupDB(db *sql.DB) error {
	cleanupSql := `
		delete from video;
		delete from taxonomy where parent_id is not null;
		delete from taxonomy;
		delete from student;
		delete from user;
		alter table video auto_increment = 1;
		alter table taxonomy auto_increment = 1;
    alter table student auto_increment = 1;
		alter table user auto_increment = 1;
	`

	for _, sql := range strings.Split(cleanupSql, ";") {
		trimmedSql := strings.TrimSpace(sql)
		if len(trimmedSql) == 0 {
			continue
		}

		stmt, err := db.Prepare(trimmedSql)
		if err != nil {
			return fmt.Errorf("Failed to prepare statement %s: %v", trimmedSql, err)
		}

		if _, err := stmt.Exec(); err != nil {
			return fmt.Errorf("Failed to cleanup db %v", err)
		}
	}

	return nil
}
