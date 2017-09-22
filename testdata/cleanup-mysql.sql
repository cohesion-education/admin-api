delete from video;
delete from taxonomy;
delete from audit_info;
delete from user;
alter table video auto_increment = 1;
alter table taxonomy auto_increment = 1;
alter table audit_info auto_increment = 1;
alter table user auto_increment = 1;
