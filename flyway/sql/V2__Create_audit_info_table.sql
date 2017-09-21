-- -----------------------------------------------------
-- Table `cohesion`.`audit_info`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `cohesion`.`audit_info` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL,
  `created_by` INT NOT NULL,
  `updated_by` INT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_audit_info_user1_idx` (`created_by` ASC),
  INDEX `fk_audit_info_user2_idx` (`updated_by` ASC),
  CONSTRAINT `fk_audit_info_created_by`
    FOREIGN KEY (`created_by`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_audit_info_updated_by`
    FOREIGN KEY (`updated_by`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
