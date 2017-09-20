-- -----------------------------------------------------
-- Table `cohesion`.`student`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `cohesion`.`student` (
  `id` INT NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `grade` INT NOT NULL,
  `school` VARCHAR(255) NOT NULL,
  `user_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC),
  INDEX `fk_student_user1_idx` (`user_id` ASC),
  CONSTRAINT `fk_student_user1`
    FOREIGN KEY (`user_id`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
