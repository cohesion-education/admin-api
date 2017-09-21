-- -----------------------------------------------------
-- Table `cohesion`.`video`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `cohesion`.`video` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NULL,
  `audit_info_id` INT NULL,
  `taxonomy_id` INT NOT NULL,
  `file_name` VARCHAR(255) NOT NULL,
  `bucket` VARCHAR(255) NOT NULL,
  `object_key` VARCHAR(255) NOT NULL,
  `key_terms` LONGTEXT NULL,
  `file_type` VARCHAR(255) NULL,
  `file_size` INT NULL,
  `state_standards` LONGTEXT NULL,
  `common_core_standards` LONGTEXT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_video_audit_info1_idx` (`audit_info_id` ASC),
  INDEX `fk_video_taxonomy1_idx` (`taxonomy_id` ASC),
  CONSTRAINT `fk_video_audit_info`
    FOREIGN KEY (`audit_info_id`)
    REFERENCES `cohesion`.`audit_info` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_video_taxonomy1`
    FOREIGN KEY (`taxonomy_id`)
    REFERENCES `cohesion`.`taxonomy` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
