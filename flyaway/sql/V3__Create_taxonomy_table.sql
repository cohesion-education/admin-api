-- -----------------------------------------------------
-- Table `cohesion`.`taxonomy`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `cohesion`.`taxonomy` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NULL,
  `parent_id` INT NULL,
  `audit_info_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_taxonomy_taxonomy1_idx` (`parent_id` ASC),
  INDEX `fk_taxonomy_audit_info1_idx` (`audit_info_id` ASC),
  CONSTRAINT `fk_taxonomy_parent_id`
    FOREIGN KEY (`parent_id`)
    REFERENCES `cohesion`.`taxonomy` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_taxonomy_audit_info1`
    FOREIGN KEY (`audit_info_id`)
    REFERENCES `cohesion`.`audit_info` (`id`)
    ON DELETE CASCADE
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
