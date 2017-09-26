-- -----------------------------------------------------
-- Table `cohesion`.`student`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `cohesion`.`student` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NOT NULL,
  `grade` VARCHAR(255) NOT NULL,
  `school` VARCHAR(255) NOT NULL,
  `user_id` INT NOT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL,
  `created_by` INT NOT NULL,
  `updated_by` INT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC),
  INDEX `fk_student_user_idx` (`user_id` ASC),
  INDEX `fk_student_created_by_idx` (`created_by` ASC),
  INDEX `fk_student_updated_by_idx` (`updated_by` ASC),
  CONSTRAINT `ffk_student_created_by`
    FOREIGN KEY (`created_by`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_student_updated_by`
    FOREIGN KEY (`updated_by`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_student_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
