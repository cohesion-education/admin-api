-- -----------------------------------------------------
-- Table `cohesion`.`taxonomy`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `cohesion`.`taxonomy` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(255) NULL,
  `parent_id` INT NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL,
  `created_by` INT NOT NULL,
  `updated_by` INT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_taxonomy_taxonomy_idx` (`parent_id` ASC),
  INDEX `fk_taxonomy_created_by_idx` (`created_by` ASC),
  INDEX `fk_taxonomy_updated_by_idx` (`updated_by` ASC),
  CONSTRAINT `ffk_taxonomy_created_by`
    FOREIGN KEY (`created_by`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_taxonomy_updated_by`
    FOREIGN KEY (`updated_by`)
    REFERENCES `cohesion`.`user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_taxonomy_parent_id`
    FOREIGN KEY (`parent_id`)
    REFERENCES `cohesion`.`taxonomy` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
