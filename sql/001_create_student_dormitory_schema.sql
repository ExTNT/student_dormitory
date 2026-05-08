-- Dormitory management system schema
-- Target database: student_dormitory
-- Owner: turing
--
-- Usage:
--   PGPASSWORD=10928 psql -U turing -d student_dormitory -f sql/001_create_student_dormitory_schema.sql

BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    password_hash VARCHAR(128) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (
        role IN ('student', 'repair_staff', 'cleaning_staff', 'dormitory_manager', 'system_admin')
    ),
    name VARCHAR(32) NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash CHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS buildings (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    location VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS rooms (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    building_id INT NOT NULL REFERENCES buildings(id) ON DELETE RESTRICT,
    room_number VARCHAR(16) NOT NULL,
    floor SMALLINT NOT NULL,
    total_beds SMALLINT NOT NULL CHECK (total_beds > 0),
    water_balance NUMERIC(8,2) NOT NULL DEFAULT 0.00 CHECK (water_balance >= 0),
    electricity_balance NUMERIC(8,2) NOT NULL DEFAULT 0.00 CHECK (electricity_balance >= 0),
    CONSTRAINT uq_rooms_building_room UNIQUE (building_id, room_number)
);

CREATE TABLE IF NOT EXISTS beds (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    room_id INT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    bed_label VARCHAR(8) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'occupied')),
    student_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    occupied_since DATE,
    CONSTRAINT uq_beds_room_label UNIQUE (room_id, bed_label)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_beds_student_id_not_null
    ON beds(student_id)
    WHERE student_id IS NOT NULL;

CREATE TABLE IF NOT EXISTS lifestyle_surveys (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sleep_time TIME,
    smoking SMALLINT NOT NULL DEFAULT 0 CHECK (smoking IN (0, 1)),
    snoring SMALLINT NOT NULL DEFAULT 0 CHECK (snoring IN (0, 1)),
    study_habit VARCHAR(255),
    remarks TEXT,
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS allocation_requests (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recommended_room_id INT NOT NULL REFERENCES rooms(id) ON DELETE RESTRICT,
    recommended_bed_id INT NOT NULL REFERENCES beds(id) ON DELETE RESTRICT,
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    admin_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS leave_applications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL DEFAULT 'normal' CHECK (type IN ('normal', 'holiday')),
    destination VARCHAR(128) NOT NULL,
    emergency_contact VARCHAR(64) NOT NULL,
    return_time TIMESTAMPTZ NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    manager_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS utility_payments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    room_id INT NOT NULL REFERENCES rooms(id) ON DELETE RESTRICT,
    payer_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    amount NUMERIC(8,2) NOT NULL CHECK (amount > 0),
    payment_type VARCHAR(12) NOT NULL CHECK (payment_type IN ('water', 'electricity', 'both')),
    paid_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS repair_requests (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    room_id INT NOT NULL REFERENCES rooms(id) ON DELETE RESTRICT,
    description TEXT NOT NULL,
    status VARCHAR(12) NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'accepted', 'repaired', 'completed', 'rejected')
    ),
    repair_staff_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    repair_description TEXT,
    reviewer_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    review_comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    accepted_at TIMESTAMPTZ,
    repaired_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS cleaning_requests (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    building_id INT NOT NULL REFERENCES buildings(id) ON DELETE RESTRICT,
    location_desc VARCHAR(255) NOT NULL,
    status VARCHAR(12) NOT NULL DEFAULT 'pending' CHECK (
        status IN ('pending', 'accepted', 'cleaned', 'completed', 'rejected')
    ),
    cleaner_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    reviewer_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    review_comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    accepted_at TIMESTAMPTZ,
    cleaned_at TIMESTAMPTZ,
    reviewed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS late_return_records (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    return_date DATE NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    manager_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS room_change_requests (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_bed_id INT NOT NULL REFERENCES beds(id) ON DELETE RESTRICT,
    target_room_id INT REFERENCES rooms(id) ON DELETE RESTRICT,
    target_bed_id INT REFERENCES beds(id) ON DELETE RESTRICT,
    recommended_bed_id INT REFERENCES beds(id) ON DELETE RESTRICT,
    reason TEXT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    manager_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS notifications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    recipient_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    room_id INT REFERENCES rooms(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    type VARCHAR(12) NOT NULL DEFAULT 'general' CHECK (type IN ('low_balance', 'general')),
    is_read SMALLINT NOT NULL DEFAULT 0 CHECK (is_read IN (0, 1)),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS attachments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_type VARCHAR(32) NOT NULL CHECK (owner_type IN ('user_avatar', 'repair', 'cleaning')),
    owner_id BIGINT NOT NULL,
    category VARCHAR(32) NOT NULL CHECK (category IN ('avatar', 'before', 'after')),
    sort_order SMALLINT NOT NULL DEFAULT 0,
    file_name VARCHAR(255),
    content_type VARCHAR(100) NOT NULL CHECK (content_type LIKE 'image/%'),
    file_data BYTEA NOT NULL,
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_attachments_owner_category CHECK (
        (owner_type = 'user_avatar' AND category = 'avatar')
        OR (owner_type = 'repair' AND category = 'after')
        OR (owner_type = 'cleaning' AND category IN ('before', 'after'))
    )
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_attachments_user_avatar
    ON attachments(owner_id)
    WHERE owner_type = 'user_avatar' AND category = 'avatar';

CREATE TABLE IF NOT EXISTS off_campus_living_applications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    retain_bed SMALLINT NOT NULL DEFAULT 0 CHECK (retain_bed IN (0, 1)),
    reason TEXT NOT NULL,
    destination VARCHAR(255),
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    manager_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    review_comment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_rooms_building_id ON rooms(building_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_beds_room_id ON beds(room_id);
CREATE INDEX IF NOT EXISTS idx_lifestyle_surveys_student_id ON lifestyle_surveys(student_id);
CREATE INDEX IF NOT EXISTS idx_allocation_requests_student_id ON allocation_requests(student_id);
CREATE INDEX IF NOT EXISTS idx_allocation_requests_status ON allocation_requests(status);
CREATE INDEX IF NOT EXISTS idx_leave_applications_student_id ON leave_applications(student_id);
CREATE INDEX IF NOT EXISTS idx_leave_applications_status ON leave_applications(status);
CREATE INDEX IF NOT EXISTS idx_utility_payments_room_id ON utility_payments(room_id);
CREATE INDEX IF NOT EXISTS idx_utility_payments_payer_id ON utility_payments(payer_id);
CREATE INDEX IF NOT EXISTS idx_repair_requests_student_id ON repair_requests(student_id);
CREATE INDEX IF NOT EXISTS idx_repair_requests_room_id ON repair_requests(room_id);
CREATE INDEX IF NOT EXISTS idx_repair_requests_status_created_at ON repair_requests(status, created_at);
CREATE INDEX IF NOT EXISTS idx_cleaning_requests_student_id ON cleaning_requests(student_id);
CREATE INDEX IF NOT EXISTS idx_cleaning_requests_building_id ON cleaning_requests(building_id);
CREATE INDEX IF NOT EXISTS idx_cleaning_requests_status_created_at ON cleaning_requests(status, created_at);
CREATE INDEX IF NOT EXISTS idx_late_return_records_student_id ON late_return_records(student_id);
CREATE INDEX IF NOT EXISTS idx_late_return_records_status ON late_return_records(status);
CREATE INDEX IF NOT EXISTS idx_room_change_requests_student_id ON room_change_requests(student_id);
CREATE INDEX IF NOT EXISTS idx_room_change_requests_status ON room_change_requests(status);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient_read ON notifications(recipient_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_room_id ON notifications(room_id);
CREATE INDEX IF NOT EXISTS idx_attachments_owner ON attachments(owner_type, owner_id, category);
CREATE INDEX IF NOT EXISTS idx_off_campus_student_id ON off_campus_living_applications(student_id);
CREATE INDEX IF NOT EXISTS idx_off_campus_status ON off_campus_living_applications(status);

CREATE OR REPLACE VIEW v_dormitory_summary AS
WITH room_stats AS (
    SELECT
        building_id,
        COUNT(*) AS total_rooms,
        COALESCE(SUM(total_beds), 0)::BIGINT AS total_beds
    FROM rooms
    GROUP BY building_id
),
bed_stats AS (
    SELECT
        r.building_id,
        COUNT(bed.id) FILTER (WHERE bed.status = 'occupied') AS occupied_beds,
        COUNT(bed.id) FILTER (WHERE bed.status = 'available') AS free_beds
    FROM rooms r
    LEFT JOIN beds bed ON bed.room_id = r.id
    GROUP BY r.building_id
)
SELECT
    b.id AS building_id,
    b.name AS building_name,
    COALESCE(rs.total_rooms, 0) AS total_rooms,
    COALESCE(rs.total_beds, 0) AS total_beds,
    COALESCE(bs.occupied_beds, 0) AS occupied_beds,
    COALESCE(bs.free_beds, 0) AS free_beds
FROM buildings b
LEFT JOIN room_stats rs ON rs.building_id = b.id
LEFT JOIN bed_stats bs ON bs.building_id = b.id;

CREATE OR REPLACE VIEW v_available_beds AS
SELECT
    bed.id AS bed_id,
    r.id AS room_id,
    r.room_number,
    bed.bed_label,
    b.id AS building_id,
    b.name AS building_name,
    r.floor
FROM beds bed
JOIN rooms r ON r.id = bed.room_id
JOIN buildings b ON b.id = r.building_id
WHERE bed.status = 'available';

CREATE OR REPLACE VIEW v_student_roommates AS
SELECT
    self_bed.student_id AS student_id,
    roommate.id AS roommate_id,
    roommate.name AS roommate_name,
    roommate.phone AS roommate_phone,
    roommate_bed.bed_label,
    avatar.id AS avatar_attachment_id
FROM beds self_bed
JOIN beds roommate_bed
    ON roommate_bed.room_id = self_bed.room_id
   AND roommate_bed.status = 'occupied'
   AND roommate_bed.student_id IS NOT NULL
   AND roommate_bed.student_id <> self_bed.student_id
JOIN users roommate ON roommate.id = roommate_bed.student_id
LEFT JOIN attachments avatar
    ON avatar.owner_type = 'user_avatar'
   AND avatar.category = 'avatar'
   AND avatar.owner_id = roommate.id
WHERE self_bed.status = 'occupied'
  AND self_bed.student_id IS NOT NULL;

CREATE OR REPLACE VIEW v_low_balance_rooms AS
SELECT
    id AS room_id,
    building_id,
    room_number,
    water_balance,
    electricity_balance
FROM rooms
WHERE water_balance < 5 OR electricity_balance < 5;

CREATE OR REPLACE VIEW v_pending_repairs AS
SELECT
    rr.id AS request_id,
    rr.status,
    student.name AS student_name,
    r.room_number,
    rr.description,
    rr.created_at,
    repair_staff.name AS repair_staff_name,
    reviewer.name AS reviewer_name
FROM repair_requests rr
JOIN users student ON student.id = rr.student_id
JOIN rooms r ON r.id = rr.room_id
LEFT JOIN users repair_staff ON repair_staff.id = rr.repair_staff_id
LEFT JOIN users reviewer ON reviewer.id = rr.reviewer_id
WHERE rr.status IN ('pending', 'accepted', 'repaired');

CREATE OR REPLACE VIEW v_pending_cleanings AS
SELECT
    cr.id AS request_id,
    cr.status,
    student.name AS student_name,
    b.name AS building_name,
    cr.location_desc,
    cr.created_at,
    cleaner.name AS cleaner_name,
    reviewer.name AS reviewer_name
FROM cleaning_requests cr
JOIN users student ON student.id = cr.student_id
JOIN buildings b ON b.id = cr.building_id
LEFT JOIN users cleaner ON cleaner.id = cr.cleaner_id
LEFT JOIN users reviewer ON reviewer.id = cr.reviewer_id
WHERE cr.status IN ('pending', 'accepted', 'cleaned');

CREATE OR REPLACE VIEW v_my_requests AS
SELECT student_id, '维修'::TEXT AS request_type, id AS request_id, status, created_at, left(description, 120) AS detail
FROM repair_requests
UNION ALL
SELECT student_id, '保洁'::TEXT AS request_type, id AS request_id, status, created_at, location_desc AS detail
FROM cleaning_requests
UNION ALL
SELECT student_id, '离校'::TEXT AS request_type, id AS request_id, status, created_at, destination AS detail
FROM leave_applications
UNION ALL
SELECT student_id, '换寝'::TEXT AS request_type, id AS request_id, status, created_at, left(reason, 120) AS detail
FROM room_change_requests
UNION ALL
SELECT student_id, '晚归'::TEXT AS request_type, id AS request_id, status, created_at, return_date::TEXT AS detail
FROM late_return_records
UNION ALL
SELECT student_id, '校外居住'::TEXT AS request_type, id AS request_id, status, created_at, COALESCE(destination, left(reason, 120)) AS detail
FROM off_campus_living_applications
UNION ALL
SELECT
    ar.student_id,
    '分配'::TEXT AS request_type,
    ar.id AS request_id,
    ar.status,
    ar.created_at,
    b.name || ' ' || r.room_number || ' ' || bed.bed_label || '床' AS detail
FROM allocation_requests ar
JOIN rooms r ON r.id = ar.recommended_room_id
JOIN buildings b ON b.id = r.building_id
JOIN beds bed ON bed.id = ar.recommended_bed_id
UNION ALL
SELECT
    payer_id AS student_id,
    '缴费'::TEXT AS request_type,
    id AS request_id,
    'paid'::TEXT AS status,
    paid_at AS created_at,
    payment_type || ' ' || amount::TEXT || '元' AS detail
FROM utility_payments;

CREATE OR REPLACE VIEW v_attachment_metadata AS
SELECT
    id,
    owner_type,
    owner_id,
    category,
    sort_order,
    content_type,
    file_name,
    uploaded_at
FROM attachments;

CREATE OR REPLACE FUNCTION fn_repair_status_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IS DISTINCT FROM OLD.status THEN
        IF NEW.status = 'accepted' AND NEW.accepted_at IS NULL THEN
            NEW.accepted_at := now();
        ELSIF NEW.status = 'repaired' AND NEW.repaired_at IS NULL THEN
            NEW.repaired_at := now();
        ELSIF NEW.status IN ('completed', 'rejected') AND NEW.reviewed_at IS NULL THEN
            NEW.reviewed_at := now();
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_repair_status_timestamp ON repair_requests;
CREATE TRIGGER trg_repair_status_timestamp
BEFORE UPDATE OF status ON repair_requests
FOR EACH ROW
EXECUTE FUNCTION fn_repair_status_timestamp();

CREATE OR REPLACE FUNCTION fn_cleaning_status_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IS DISTINCT FROM OLD.status THEN
        IF NEW.status = 'accepted' AND NEW.accepted_at IS NULL THEN
            NEW.accepted_at := now();
        ELSIF NEW.status = 'cleaned' AND NEW.cleaned_at IS NULL THEN
            NEW.cleaned_at := now();
        ELSIF NEW.status IN ('completed', 'rejected') AND NEW.reviewed_at IS NULL THEN
            NEW.reviewed_at := now();
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_cleaning_status_timestamp ON cleaning_requests;
CREATE TRIGGER trg_cleaning_status_timestamp
BEFORE UPDATE OF status ON cleaning_requests
FOR EACH ROW
EXECUTE FUNCTION fn_cleaning_status_timestamp();

CREATE OR REPLACE FUNCTION fn_bed_consistency()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = 'occupied' THEN
        IF NEW.student_id IS NULL OR NEW.occupied_since IS NULL THEN
            RAISE EXCEPTION 'occupied bed requires student_id and occupied_since';
        END IF;
    ELSIF NEW.status = 'available' THEN
        IF NEW.student_id IS NOT NULL OR NEW.occupied_since IS NOT NULL THEN
            RAISE EXCEPTION 'available bed requires student_id and occupied_since to be NULL';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_bed_consistency ON beds;
CREATE TRIGGER trg_bed_consistency
BEFORE INSERT OR UPDATE ON beds
FOR EACH ROW
EXECUTE FUNCTION fn_bed_consistency();

CREATE OR REPLACE FUNCTION fn_low_balance_notification()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.water_balance < 5
       AND (TG_OP = 'INSERT' OR OLD.water_balance IS NULL OR OLD.water_balance >= 5 OR OLD.water_balance IS DISTINCT FROM NEW.water_balance) THEN
        INSERT INTO notifications (recipient_id, room_id, message, type)
        SELECT
            bed.student_id,
            NEW.id,
            '宿舍水费余额低于5元，当前余额：' || NEW.water_balance::TEXT || '元',
            'low_balance'
        FROM beds bed
        WHERE bed.room_id = NEW.id
          AND bed.status = 'occupied'
          AND bed.student_id IS NOT NULL;
    END IF;

    IF NEW.electricity_balance < 5
       AND (TG_OP = 'INSERT' OR OLD.electricity_balance IS NULL OR OLD.electricity_balance >= 5 OR OLD.electricity_balance IS DISTINCT FROM NEW.electricity_balance) THEN
        INSERT INTO notifications (recipient_id, room_id, message, type)
        SELECT
            bed.student_id,
            NEW.id,
            '宿舍电费余额低于5元，当前余额：' || NEW.electricity_balance::TEXT || '元',
            'low_balance'
        FROM beds bed
        WHERE bed.room_id = NEW.id
          AND bed.status = 'occupied'
          AND bed.student_id IS NOT NULL;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_low_balance_notification ON rooms;
CREATE TRIGGER trg_low_balance_notification
AFTER UPDATE OF water_balance, electricity_balance ON rooms
FOR EACH ROW
EXECUTE FUNCTION fn_low_balance_notification();

CREATE OR REPLACE FUNCTION fn_application_resolved_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status IS DISTINCT FROM OLD.status
       AND NEW.status IN ('approved', 'rejected')
       AND NEW.resolved_at IS NULL THEN
        NEW.resolved_at := now();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_leave_resolved_at ON leave_applications;
CREATE TRIGGER trg_leave_resolved_at
BEFORE UPDATE OF status ON leave_applications
FOR EACH ROW
EXECUTE FUNCTION fn_application_resolved_at();

DROP TRIGGER IF EXISTS trg_late_return_resolved_at ON late_return_records;
CREATE TRIGGER trg_late_return_resolved_at
BEFORE UPDATE OF status ON late_return_records
FOR EACH ROW
EXECUTE FUNCTION fn_application_resolved_at();

DROP TRIGGER IF EXISTS trg_room_change_resolved_at ON room_change_requests;
CREATE TRIGGER trg_room_change_resolved_at
BEFORE UPDATE OF status ON room_change_requests
FOR EACH ROW
EXECUTE FUNCTION fn_application_resolved_at();

CREATE OR REPLACE FUNCTION fn_room_change_target_consistency()
RETURNS TRIGGER AS $$
DECLARE
    from_building_id INT;
    target_building_id INT;
    target_bed_status VARCHAR(10);
BEGIN
    IF NEW.target_room_id IS NULL AND NEW.target_bed_id IS NULL THEN
        RETURN NEW;
    END IF;

    IF NEW.target_room_id IS NULL OR NEW.target_bed_id IS NULL THEN
        RAISE EXCEPTION 'target_room_id and target_bed_id must be provided together';
    END IF;

    SELECT r.building_id
    INTO from_building_id
    FROM beds b
    JOIN rooms r ON r.id = b.room_id
    WHERE b.id = NEW.from_bed_id;

    SELECT r.building_id, b.status
    INTO target_building_id, target_bed_status
    FROM beds b
    JOIN rooms r ON r.id = b.room_id
    WHERE b.id = NEW.target_bed_id
      AND b.room_id = NEW.target_room_id;

    IF target_building_id IS NULL THEN
        RAISE EXCEPTION 'target_bed_id must belong to target_room_id';
    END IF;

    IF from_building_id IS DISTINCT FROM target_building_id THEN
        RAISE EXCEPTION 'target bed must be in the same building as the original bed';
    END IF;

    IF target_bed_status <> 'available' THEN
        RAISE EXCEPTION 'target bed must be available';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_room_change_target_consistency ON room_change_requests;
CREATE TRIGGER trg_room_change_target_consistency
BEFORE INSERT OR UPDATE OF from_bed_id, target_room_id, target_bed_id ON room_change_requests
FOR EACH ROW
EXECUTE FUNCTION fn_room_change_target_consistency();

DROP TRIGGER IF EXISTS trg_off_campus_resolved_at ON off_campus_living_applications;
CREATE TRIGGER trg_off_campus_resolved_at
BEFORE UPDATE OF status ON off_campus_living_applications
FOR EACH ROW
EXECUTE FUNCTION fn_application_resolved_at();

DROP TRIGGER IF EXISTS trg_allocation_resolved_at ON allocation_requests;
CREATE TRIGGER trg_allocation_resolved_at
BEFORE UPDATE OF status ON allocation_requests
FOR EACH ROW
EXECUTE FUNCTION fn_application_resolved_at();

CREATE OR REPLACE FUNCTION fn_apply_utility_payment()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.payment_type = 'water' THEN
        UPDATE rooms
        SET water_balance = water_balance + NEW.amount
        WHERE id = NEW.room_id;
    ELSIF NEW.payment_type = 'electricity' THEN
        UPDATE rooms
        SET electricity_balance = electricity_balance + NEW.amount
        WHERE id = NEW.room_id;
    ELSIF NEW.payment_type = 'both' THEN
        UPDATE rooms
        SET water_balance = water_balance + NEW.amount / 2,
            electricity_balance = electricity_balance + NEW.amount / 2
        WHERE id = NEW.room_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_apply_utility_payment ON utility_payments;
CREATE TRIGGER trg_apply_utility_payment
AFTER INSERT ON utility_payments
FOR EACH ROW
EXECUTE FUNCTION fn_apply_utility_payment();

COMMIT;
