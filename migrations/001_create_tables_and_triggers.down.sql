-- Down Migration for tables and related objects Spy Cat Agency

-- 1. Removal of triggers on the Targets table
DROP TRIGGER IF EXISTS trg_freeze_notes ON targets;
DROP TRIGGER IF EXISTS trg_check_min_targets ON targets;
DROP TRIGGER IF EXISTS trg_check_max_targets ON targets;

-- 2. Removing the trigger on the Missions table
DROP TRIGGER IF EXISTS trg_prevent_mission_delete ON missions;

-- 3. ВIdentification of functions
DROP FUNCTION IF EXISTS freeze_notes_if_completed();
DROP FUNCTION IF EXISTS check_min_targets_and_prevent_delete_completed();
DROP FUNCTION IF EXISTS check_max_targets();
DROP FUNCTION IF EXISTS prevent_mission_delete_if_assigned();

-- 4. Index removal (if not automatically deleted together with tables)
DROP INDEX IF EXISTS idx_unique_active_mission;
DROP INDEX IF EXISTS idx_missions_cat_id;
DROP INDEX IF EXISTS idx_missions_completed;
DROP INDEX IF EXISTS idx_spy_cats_name;
DROP INDEX IF EXISTS idx_spy_cats_breed;
DROP INDEX IF EXISTS idx_targets_mission_id;
DROP INDEX IF EXISTS idx_targets_complete;
DROP INDEX IF EXISTS idx_targets_country;

-- 5. Removing the tables in reverse order of addiction
DROP TABLE IF EXISTS targets;
DROP TABLE IF EXISTS missions;
DROP TABLE IF EXISTS spy_cats;
