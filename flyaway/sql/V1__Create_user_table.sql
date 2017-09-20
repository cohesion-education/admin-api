-- -----------------------------------------------------
-- Table `cohesion`.`user`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `cohesion`.`user` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `email` VARCHAR(45) NOT NULL,
  `full_name` VARCHAR(255) NOT NULL,
  `last_name` VARCHAR(255) NOT NULL,
  `enabled` TINYINT NULL,
  `first_name` VARCHAR(255) NOT NULL,
  `profile_pic_url` MEDIUMTEXT NULL,
  `nickname` VARCHAR(255) NULL,
  `email_verified` TINYINT NULL,
  `state` VARCHAR(255) NULL,
  `county` MEDIUMTEXT NULL,
  `sub` VARCHAR(45) NULL,
  `user_id` VARCHAR(45) NULL,
  `newsletter` TINYINT NULL,
  `beta_program` TINYINT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC))
ENGINE = InnoDB;
