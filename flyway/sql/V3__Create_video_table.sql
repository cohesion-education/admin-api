-- -----------------------------------------------------
-- Table `video`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `video` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NULL,
  `taxonomy_id` INT NOT NULL,
  `file_name` VARCHAR(255) NOT NULL,
  `bucket` VARCHAR(255) NOT NULL,
  `object_key` VARCHAR(255) NOT NULL,
  `key_terms` LONGTEXT NULL,
  `file_type` VARCHAR(255) NULL,
  `file_size` INT NULL,
  `state_standards` LONGTEXT NULL,
  `common_core_standards` LONGTEXT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL,
  `created_by` INT NOT NULL,
  `updated_by` INT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_video_taxonomy_idx` (`taxonomy_id` ASC),
  INDEX `fk_video_created_by_idx` (`created_by` ASC),
  INDEX `fk_video_updated_by_idx` (`updated_by` ASC),
  CONSTRAINT `ffk_video_created_by`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_video_updated_by`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_video_taxonomy`
    FOREIGN KEY (`taxonomy_id`)
    REFERENCES `taxonomy` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
