-- Down migration для таблиць і супутніх об'єктів Spy Cat Agency

-- 1. Видалення тригерів на таблиці targets
DROP TRIGGER IF EXISTS trg_freeze_notes ON targets;
DROP TRIGGER IF EXISTS trg_check_min_targets ON targets;
DROP TRIGGER IF EXISTS trg_check_max_targets ON targets;

-- 2. Видалення тригера на таблиці missions
DROP TRIGGER IF EXISTS trg_prevent_mission_delete ON missions;

-- 3. Видалення функцій
DROP FUNCTION IF EXISTS freeze_notes_if_completed();
DROP FUNCTION IF EXISTS check_min_targets_and_prevent_delete_completed();
DROP FUNCTION IF EXISTS check_max_targets();
DROP FUNCTION IF EXISTS prevent_mission_delete_if_assigned();

-- 4. Видалення індексів (якщо вони не видаляються автоматично разом з таблицями)
DROP INDEX IF EXISTS idx_unique_active_mission;
DROP INDEX IF EXISTS idx_missions_cat_id;
DROP INDEX IF EXISTS idx_missions_completed;
DROP INDEX IF EXISTS idx_spy_cats_name;
DROP INDEX IF EXISTS idx_spy_cats_breed;
DROP INDEX IF EXISTS idx_targets_mission_id;
DROP INDEX IF EXISTS idx_targets_complete;
DROP INDEX IF EXISTS idx_targets_country;

-- 5. Видалення таблиць у зворотному порядку залежностей
DROP TABLE IF EXISTS targets;
DROP TABLE IF EXISTS missions;
DROP TABLE IF EXISTS spy_cats;
