-- A table for spy cats
CREATE TABLE spy_cats (
                          id SERIAL PRIMARY KEY,
                          name TEXT NOT NULL,
                          years_of_experience INTEGER NOT NULL,
                          breed TEXT NOT NULL,
                          salary NUMERIC(10,2) NOT NULL,
                          created_at TIMESTAMPTZ DEFAULT now()
);

-- Indexes to quickly search cats
CREATE INDEX idx_spy_cats_name ON spy_cats(name);
CREATE INDEX idx_spy_cats_breed ON spy_cats(breed);

-- Таблиця для місій
CREATE TABLE missions (
                          id SERIAL PRIMARY KEY,
                          cat_id INTEGER REFERENCES spy_cats(id) ON DELETE SET NULL,
                          completed BOOLEAN NOT NULL DEFAULT FALSE,
                          created_at TIMESTAMPTZ DEFAULT now()
);

-- Unique index providing: one cat can only have one active mission
CREATE UNIQUE INDEX idx_unique_active_mission
    ON missions(cat_id)
    WHERE completed = false AND cat_id IS NOT NULL;

-- Indexes for a quick search for cat and condition
CREATE INDEX idx_missions_cat_id ON missions(cat_id);
CREATE INDEX idx_missions_completed ON missions(completed);

-- Trigger to prohibit the removal of a mission if it is already assigned a cat
CREATE OR REPLACE FUNCTION prevent_mission_delete_if_assigned()
RETURNS trigger AS $$
BEGIN
    IF OLD.cat_id IS NOT NULL THEN
        RAISE EXCEPTION 'Cannot delete mission % because it is assigned to a cat', OLD.id;
END IF;
RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_mission_delete
    BEFORE DELETE ON missions
    FOR EACH ROW
    EXECUTE FUNCTION prevent_mission_delete_if_assigned();

-- Table for mission purposes
CREATE TABLE targets (
                         id SERIAL PRIMARY KEY,
                         mission_id INTEGER NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
                         name TEXT NOT NULL,
                         country TEXT NOT NULL,
                         notes TEXT,
                         complete BOOLEAN NOT NULL DEFAULT FALSE,
                         created_at TIMESTAMPTZ DEFAULT now()
);

-- Indexes to quickly search for goals
CREATE INDEX idx_targets_mission_id ON targets(mission_id);
CREATE INDEX idx_targets_complete ON targets(complete);
CREATE INDEX idx_targets_country ON targets(country);

-- Trigger to verify: when adding a new target is checked that the mission is not completed and the number of targets does not exceed 3
CREATE OR REPLACE FUNCTION check_max_targets()
RETURNS trigger AS $$
DECLARE
target_count INTEGER;
    mission_completed BOOLEAN;
BEGIN
SELECT count(*) INTO target_count FROM targets WHERE mission_id = NEW.mission_id;
SELECT completed INTO mission_completed FROM missions WHERE id = NEW.mission_id;

IF mission_completed THEN
        RAISE EXCEPTION 'Cannot add target to mission % because the mission is already completed', NEW.mission_id;
END IF;

    IF target_count >= 3 THEN
        RAISE EXCEPTION 'Mission % already has maximum number of targets (3)', NEW.mission_id;
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_max_targets
    BEFORE INSERT ON targets
    FOR EACH ROW
    EXECUTE FUNCTION check_max_targets();

-- Trigger to check: Prohibition of deleting the completed target and providing a minimum (at least 1 target should remain)
CREATE OR REPLACE FUNCTION check_min_targets_and_prevent_delete_completed()
RETURNS trigger AS $$
DECLARE
target_count INTEGER;
BEGIN
    IF OLD.complete THEN
        RAISE EXCEPTION 'Cannot delete target % because it is already completed', OLD.id;
END IF;

SELECT count(*) INTO target_count FROM targets WHERE mission_id = OLD.mission_id;
IF (target_count - 1) < 1 THEN
        RAISE EXCEPTION 'Mission % must have at least one target', OLD.mission_id;
END IF;

RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_min_targets
    BEFORE DELETE ON targets
    FOR EACH ROW
    EXECUTE FUNCTION check_min_targets_and_prevent_delete_completed();

-- Trigger to "freeze" notes if either target or mission is completed
CREATE OR REPLACE FUNCTION freeze_notes_if_completed()
RETURNS trigger AS $$
DECLARE
mission_completed BOOLEAN;
BEGIN
SELECT completed INTO mission_completed FROM missions WHERE id = NEW.mission_id;

IF (OLD.complete = TRUE OR mission_completed = TRUE)
       AND (NEW.notes IS DISTINCT FROM OLD.notes) THEN
        RAISE EXCEPTION 'Cannot update notes because either the target or the mission is completed';
END IF;

RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_freeze_notes
    BEFORE UPDATE OF notes ON targets
    FOR EACH ROW
    EXECUTE FUNCTION freeze_notes_if_completed();
