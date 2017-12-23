DROP FUNCTION IF EXISTS calc_price_cchf;
DROP FUNCTION IF EXISTS calc_instance_discount_cchf;

DELIMITER //

CREATE FUNCTION calc_price_cchf(
  _pattern_id           INT(10) UNSIGNED,
  _length_mm            INT(10),
  _width_mm             INT(10),
  _price_per_pearl_cchf INT(10)
)
  RETURNS INT(10) UNSIGNED
READS SQL DATA
DETERMINISTIC
  BEGIN
    DECLARE ret_ INT(10) UNSIGNED;

    -- todo: implement case where length_mm > 180

    SELECT ROUND(
        (p.price_cchf                                       -- base price
         + p.price_cchf_cm2 * _length_mm * _width_mm / 100  -- price per area
         + COALESCE(
             _price_per_pearl_cchf                          -- price for pearls
             * FLOOR(
                 p.numb_pearls + -- base numb of pearls
                 p.numb_pearls_10cm * _length_mm / 100      -- numb of pearls per length
             )
             , 0                                            -- 0 if any pearl field is NULL
         )
        ),
        -1                                                 -- round to 10 rappen
    )
    INTO ret_
    FROM pattern p
    WHERE p.pattern_id = _pattern_id;

    RETURN ret_;

  END//

CREATE FUNCTION calc_instance_discount_cchf(
  _price_cchf INT(10) UNSIGNED
)
  RETURNS INT(10) UNSIGNED
READS SQL DATA
DETERMINISTIC
  BEGIN
    RETURN ROUND(
        _price_cchf * 0.1,
        -1 -- round to 10 rappen
    );
  END//
