table "customers" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "first_name" {
    null = false
    type = text
  }
  column "last_name" {
    null = false
    type = text
  }
  column "email" {
    null = false
    type = text
  }
  column "phone_number" {
    null = true
    type = text
  }
  column "address" {
    null = false
    type = text
  }
  column "city" {
    null = true
    type = text
  }
  column "state" {
    null = true
    type = text
  }
  column "zip_code" {
    null = true
    type = text
  }
  column "is_active" {
    null    = true
    type    = boolean
    default = true
  }
  column "property_type" {
    null = true
    type = text
  }
  column "roof_type" {
    null = true
    type = text
  }
  column "home_ownership_type" {
    null = true
    type = text
  }
  column "average_monthly_bill" {
    null = true
    type = numeric
  }
  column "utility_provider" {
    null = true
    type = text
  }
  column "lead_source" {
    null = true
    type = text
  }
  column "referral_code" {
    null = true
    type = text
  }
  column "status" {
    null    = true
    type    = text
    default = "prospect"
  }
  column "notes" {
    null = true
    type = text
  }
  column "preferred_contact_method" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_customers_email" {
    unique  = true
    columns = [column.email]
  }
}
table "leads" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "external_lead_id" {
    null = true
    type = bigint
  }
  column "sync_status" {
    null    = true
    type    = text
    default = "pending"
  }
  column "last_synced_at" {
    null = true
    type = timestamptz
  }
  column "state" {
    null    = false
    type    = bigint
    default = 0
  }
  column "customer_id" {
    null = false
    type = bigint
  }
  column "creator_id" {
    null = true
    type = bigint
  }
  column "latitude" {
    null = false
    type = numeric
  }
  column "longitude" {
    null = false
    type = numeric
  }
  column "address" {
    null = true
    type = text
  }
  column "source" {
    null    = false
    type    = bigint
    default = 0
  }
  column "promo_code" {
    null = true
    type = text
  }
  column "is_2d" {
    null    = true
    type    = boolean
    default = false
  }
  column "kwh_usage" {
    null = true
    type = numeric
  }
  column "kwh_per_kw_manual" {
    null = true
    type = bigint
  }
  column "electricity_cost_pre" {
    null = true
    type = bigint
  }
  column "electricity_cost_post" {
    null = true
    type = bigint
  }
  column "additional_incentive" {
    null = true
    type = bigint
  }
  column "system_size" {
    null = true
    type = numeric
  }
  column "panel_count" {
    null = true
    type = bigint
  }
  column "panel_id" {
    null = true
    type = bigint
  }
  column "inverter_id" {
    null = true
    type = bigint
  }
  column "inverter_count" {
    null    = true
    type    = bigint
    default = 1
  }
  column "battery_count" {
    null    = true
    type    = bigint
    default = 0
  }
  column "utility_id" {
    null = true
    type = bigint
  }
  column "tariff_id" {
    null = true
    type = bigint
  }
  column "roof_material" {
    null = true
    type = bigint
  }
  column "surface_id" {
    null = true
    type = bigint
  }
  column "annual_production" {
    null = true
    type = numeric
  }
  column "welcome_call_state" {
    null = true
    type = bigint
  }
  column "financing_state" {
    null = true
    type = bigint
  }
  column "utility_bill_state" {
    null = true
    type = bigint
  }
  column "design_approved_state" {
    null = true
    type = bigint
  }
  column "permitting_approved_state" {
    null = true
    type = bigint
  }
  column "site_photos_state" {
    null = true
    type = bigint
  }
  column "install_crew_state" {
    null = true
    type = bigint
  }
  column "installation_state" {
    null = true
    type = bigint
  }
  column "final_inspection_state" {
    null = true
    type = bigint
  }
  column "pto_state" {
    null = true
    type = bigint
  }
  column "installation_date" {
    null = true
    type = text
  }
  column "date_ntp" {
    null = true
    type = text
  }
  column "date_installed" {
    null = true
    type = text
  }
  column "lightfusion_3d_project_id" {
    null = true
    type = bigint
  }
  column "lightfusion_3d_house_id" {
    null = true
    type = bigint
  }
  column "model_3d_status" {
    null = true
    type = text
  }
  column "model_3d_created_at" {
    null = true
    type = timestamptz
  }
  column "model_3d_completed_at" {
    null = true
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_leads_external_lead_id" {
    unique  = true
    columns = [column.external_lead_id]
  }
}
table "projects" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "customer_id" {
    null = true
    type = bigint
  }
  column "user_id" {
    null = true
    type = bigint
  }
  column "name" {
    null = true
    type = text
  }
  column "description" {
    null = true
    type = text
  }
  column "status" {
    null    = true
    type    = text
    default = "draft"
  }
  column "address" {
    null = true
    type = text
  }
  primary_key {
    columns = [column.id]
  }
}
table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = bigserial
  }
  column "created_at" {
    null = true
    type = timestamptz
  }
  column "updated_at" {
    null = true
    type = timestamptz
  }
  column "firstname" {
    null = true
    type = text
  }
  column "lastname" {
    null = true
    type = text
  }
  column "email" {
    null = true
    type = text
  }
  column "password" {
    null = true
    type = text
  }
  column "type" {
    null = true
    type = smallint
  }
  column "phone_number" {
    null = true
    type = text
  }
  column "address" {
    null = true
    type = text
  }
  column "company_id" {
    null = true
    type = bigint
  }
  column "creator_id" {
    null = true
    type = bigint
  }
  column "picture_path" {
    null = true
    type = text
  }
  column "disabled" {
    null    = true
    type    = boolean
    default = false
  }
  column "is_manager" {
    null    = true
    type    = boolean
    default = false
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_users_email" {
    unique  = true
    columns = [column.email]
  }
}
schema "public" {
  comment = "standard public schema"
}
