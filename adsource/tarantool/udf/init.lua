
-- create game object
start = function(self)
  -- create spaces and indexes
  box.once('init', function()
  end
  return false
ebd


box.schema.user.create('rotator', { password = 'secret' })
box.schema.user.grant('rotator', 'read,execute', 'universe')

-- ----------------------------------------------------------------------------
-- Campaigns
-- ----------------------------------------------------------------------------

cmp = box.schema.space.create('campaigns')

cmp:format({
  {name = 'id',                 type = 'unsigned'},
  {name = 'pricing_model',      type = 'unsigned'},
  {name = 'active',             type = 'boolean'},
  {name = 'weight',             type = 'unsigned'},
  -- Budget Limits
  {name = 'daily_budget',       type = 'unsigned'},
  {name = 'budget',             type = 'unsigned'},
  {name = 'daily_test_budget',  type = 'unsigned'},
  {name = 'total_test_budget',  type = 'unsigned'},
  -- Targetting
  {name = 'formats',            type = 'array'},
  {name = 'zones',              type = 'array'},
  {name = 'domains',            type = 'array'},
  {name = 'categories',         type = 'array'},
  {name = 'countries',          type = 'array'},
  {name = 'geos',               type = 'array'},
  {name = 'languages',          type = 'array'},
  {name = 'device_types',       type = 'array'},
  {name = 'devices',            type = 'array'},
  {name = 'os',                 type = 'array'},
  {name = 'browsers',           type = 'array'},
  {name = 'keywords',           type = 'array'},
  {name = 'sex',                type = 'array'},
  {name = 'hours',              type = 'string'},
  -- DEBUG
  {name = 'trace',              type = 'array'},
  {name = 'trace_percent',      type = 'unsigned'},
  -- Time
  {name = 'updated_at',         type = 'integer'}
})

cmp:create_index('primary', {type = 'tree', unique = true, parts = {
  1, 'unsigned',
  2, 'unsigned',
  3, 'boolean'
}})

function campaignFilter(pricing_model)

end

-- ----------------------------------------------------------------------------
-- Advertisements
-- ----------------------------------------------------------------------------

ads = box.schema.space.create('ads')

ads:format({
  {name = 'id',                 type = 'unsigned'},
  {name = 'pricing_model',      type = 'unsigned'},
  {name = 'campaign_id',        type = 'unsigned'},
  {name = 'active',             type = 'boolean'},
  {name = 'format_id',          type = 'unsigned'},
  {name = 'weight',             type = 'unsigned'},
  {name = 'frequency_capping',  type = 'unsigned'},
  -- Items
  {name = 'link',               type = 'string'},
  {name = 'content',            type = 'map'},
  -- Budget Limits
  {name = 'bid_price',          type = 'unsigned'},
  {name = 'price',              type = 'unsigned'},
  {name = 'lead_price',         type = 'unsigned'},
  {name = 'daily_budget',       type = 'unsigned'},
  {name = 'budget',             type = 'unsigned'},
  {name = 'daily_test_budget',  type = 'unsigned'},
  {name = 'total_test_budget',  type = 'unsigned'},
  -- Targetting
  {name = 'hours',              type = 'string'},
  -- Time
  {name = 'updated_at',         type = 'integer'}
})

ads:create_index('primary', {type = 'tree', unique=true, parts = {
  1, 'unsigned',
  2, 'unsigned',
  4, 'boolean',
  5, 'unsigned'
}})
