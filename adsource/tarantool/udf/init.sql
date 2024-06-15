
CREATE TABLE campaigns
( id                  unsigned
, pricing_model       unsigned
, active              boolean
, weight              unsigned
  -- Budget Limits
, daily_budget        unsigned
, budget              unsigned
, daily_test_budget   unsigned
, total_test_budget   unsigned
  -- Targetting
, formats             array
, zones               array
, domains             array
, categories          array
, countries           array
, geos                array
, languages           array
, device_types        array
, devices             array
, os                  array
, browsers            array
, keywords            array
, sex                 array
, hours               string
  -- DEBUG
, trace               array
, trace_percent       unsigned
  -- Time
, updated_at          integer

, PRIMARY KEY (id, pricing_model, active)
);
