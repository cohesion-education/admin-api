-- -----------------------------------------------------
-- Table `student`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `payment_detail` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `token_created` INT NOT NULL,
  `token_id` VARCHAR(255) NOT NULL,
  `token_client_ip` VARCHAR(255) NULL,
  `token_used` TINYINT NULL,
  `token_live_mode` TINYINT NULL,
  `token_type` VARCHAR(255) NULL,
  `card_id` VARCHAR(255) NOT NULL,
  `card_brand` VARCHAR(255) NULL,
  `card_funding` VARCHAR(255) NULL,
  `card_tokenization_method` VARCHAR(255) NULL,
  `card_fingerprint` VARCHAR(255) NULL,
  `card_name` VARCHAR(255) NULL,
  `card_country` VARCHAR(255) NULL,
  `card_currency` VARCHAR(255) NULL,
  `card_exp_month` VARCHAR(255) NULL,
  `card_exp_year` VARCHAR(255) NULL,
  `card_cvc_check` VARCHAR(255) NULL,
  `card_last4` VARCHAR(255) NULL,
  `card_dynamic_last4` VARCHAR(255) NULL,
  `card_address_line1` VARCHAR(255) NULL,
  `card_address_line1_check` VARCHAR(255) NULL,
  `card_address_line2` VARCHAR(255) NULL,
  `card_address_line2_check` VARCHAR(255) NULL,
  `card_address_country` VARCHAR(255) NULL,
  `card_address_state` VARCHAR(255) NULL,
  `card_address_zip` VARCHAR(255) NULL,
  `card_address_zip_check` VARCHAR(255) NULL,
  `created` DATETIME NOT NULL,
  `updated` DATETIME NULL,
  `created_by` INT NOT NULL,
  `updated_by` INT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `token_id_UNIQUE` (`token_id` ASC),
  INDEX `fk_payment_detail_created_by_idx` (`created_by` ASC),
  INDEX `fk_payment_detail_updated_by_idx` (`updated_by` ASC),
  CONSTRAINT `fk_payment_detail_created_by`
    FOREIGN KEY (`created_by`)
    REFERENCES `user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_payment_detail_updated_by`
    FOREIGN KEY (`updated_by`)
    REFERENCES `user` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
