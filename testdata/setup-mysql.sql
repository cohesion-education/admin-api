insert into user (email, full_name, first_name, last_name) values ('test@cohesioned.io', 'Test User', 'Test', 'User');
insert into audit_info (created, created_by) values (now(), 1);
insert into audit_info (created, created_by) values (now(), 1);
insert into taxonomy (name, audit_info_id) values ('test-taxonomy', 1);
insert into video (title, file_name, bucket, object_key, taxonomy_id, audit_info_id) values ('test-video', 'test-file.fake', 'test-bucket', 'test-obj-key', 1, 2);
