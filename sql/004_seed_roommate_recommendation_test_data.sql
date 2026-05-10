-- Roommate recommendation test data
-- Usage:
--   PGPASSWORD=10928 psql -U turing -d student_dormitory -f sql/004_seed_roommate_recommendation_test_data.sql
--
-- Notes:
--   Default password for all seeded users: 123456
--   This script creates an isolated recommendation scenario.
--   Login as student101, then call POST /api/allocations.
--   Expected recommendation: 自动推荐测试楼 701 C床.

BEGIN;

WITH seed_users(username, role, name, phone) AS (
    VALUES
        ('student101', 'student', '推荐测试新生', '13810000101'),
        ('student102', 'student', '兼容室友甲', '13810000102'),
        ('student103', 'student', '兼容室友乙', '13810000103'),
        ('student104', 'student', '不兼容室友甲', '13810000104'),
        ('student105', 'student', '不兼容室友乙', '13810000105'),
        ('student106', 'student', '无问卷室友', '13810000106')
)
INSERT INTO users (username, password_hash, role, name, phone)
SELECT
    username,
    '$2a$12$9.fLKJqiQlU1EY08CILApe.mFLNCovH2bW.6fGp8eyWsfLaUbk7GC',
    role,
    name,
    phone
FROM seed_users
ON CONFLICT (username) DO UPDATE
SET password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    name = EXCLUDED.name,
    phone = EXCLUDED.phone;

INSERT INTO buildings (name, location)
SELECT '自动推荐测试楼', '推荐算法测试区'
WHERE NOT EXISTS (
    SELECT 1 FROM buildings WHERE name = '自动推荐测试楼'
);

WITH seed_rooms(room_number, floor, total_beds, water_balance, electricity_balance) AS (
    VALUES
        ('701', 7, 4, 30.00, 30.00),
        ('702', 7, 4, 30.00, 30.00),
        ('703', 7, 4, 30.00, 30.00)
)
INSERT INTO rooms (building_id, room_number, floor, total_beds, water_balance, electricity_balance)
SELECT b.id, sr.room_number, sr.floor, sr.total_beds, sr.water_balance, sr.electricity_balance
FROM seed_rooms sr
JOIN buildings b ON b.name = '自动推荐测试楼'
ON CONFLICT (building_id, room_number) DO UPDATE
SET floor = EXCLUDED.floor,
    total_beds = EXCLUDED.total_beds,
    water_balance = EXCLUDED.water_balance,
    electricity_balance = EXCLUDED.electricity_balance;

WITH labels(label) AS (
    VALUES ('A'), ('B'), ('C'), ('D')
),
test_rooms AS (
    SELECT r.id AS room_id
    FROM rooms r
    JOIN buildings b ON b.id = r.building_id
    WHERE b.name = '自动推荐测试楼'
      AND r.room_number IN ('701', '702', '703')
)
INSERT INTO beds (room_id, bed_label, status)
SELECT room_id, label, 'available'
FROM test_rooms
JOIN labels ON true
WHERE NOT EXISTS (
    SELECT 1 FROM beds existing
    WHERE existing.room_id = test_rooms.room_id
      AND existing.bed_label = labels.label
);

WITH test_users AS (
    SELECT id FROM users WHERE username IN (
        'student101',
        'student102',
        'student103',
        'student104',
        'student105',
        'student106'
    )
),
test_beds AS (
    SELECT bed.id
    FROM beds bed
    JOIN rooms r ON r.id = bed.room_id
    JOIN buildings b ON b.id = r.building_id
    WHERE b.name = '自动推荐测试楼'
      AND r.room_number IN ('701', '702', '703')
)
DELETE FROM allocation_requests ar
WHERE ar.student_id IN (SELECT id FROM test_users)
   OR ar.recommended_bed_id IN (SELECT id FROM test_beds);

WITH test_users AS (
    SELECT id FROM users WHERE username IN (
        'student101',
        'student102',
        'student103',
        'student104',
        'student105',
        'student106'
    )
),
test_rooms AS (
    SELECT r.id
    FROM rooms r
    JOIN buildings b ON b.id = r.building_id
    WHERE b.name = '自动推荐测试楼'
      AND r.room_number IN ('701', '702', '703')
)
UPDATE beds
SET status = 'available',
    student_id = NULL,
    occupied_since = NULL
WHERE room_id IN (SELECT id FROM test_rooms)
   OR student_id IN (SELECT id FROM test_users);

WITH test_users AS (
    SELECT id FROM users WHERE username IN (
        'student101',
        'student102',
        'student103',
        'student104',
        'student105',
        'student106'
    )
)
DELETE FROM lifestyle_surveys
WHERE student_id IN (SELECT id FROM test_users);

WITH assignments(username, room_number, bed_label, occupied_since) AS (
    VALUES
        ('student102', '701', 'A', DATE '2025-09-01'),
        ('student103', '701', 'B', DATE '2025-09-01'),
        ('student104', '702', 'A', DATE '2025-09-01'),
        ('student105', '702', 'B', DATE '2025-09-01'),
        ('student106', '703', 'A', DATE '2025-09-01')
),
resolved AS (
    SELECT
        bed.id AS bed_id,
        u.id AS student_id,
        a.occupied_since
    FROM assignments a
    JOIN users u ON u.username = a.username
    JOIN buildings b ON b.name = '自动推荐测试楼'
    JOIN rooms r ON r.building_id = b.id AND r.room_number = a.room_number
    JOIN beds bed ON bed.room_id = r.id AND bed.bed_label = a.bed_label
)
UPDATE beds bed
SET status = 'occupied',
    student_id = resolved.student_id,
    occupied_since = resolved.occupied_since
FROM resolved
WHERE bed.id = resolved.bed_id;

WITH seed_surveys(username, sleep_time, smoking, snoring, study_habit, remarks) AS (
    VALUES

        ('student102', TIME '23:00', 1, 0, '考研 编程 安静', '与测试新生作息和学习习惯高度匹配'),
        ('student103', TIME '23:30', 1, 0, '考研 夜间自习 安静', '与测试新生作息和学习习惯高度匹配'),
        ('student104', TIME '02:30', 0, 1, '游戏 外放 音乐', '用于制造低分候选房间'),
        ('student105', TIME '01:50', 0, 1, '游戏 社交 外放', '用于制造低分候选房间')
)
INSERT INTO lifestyle_surveys (student_id, sleep_time, smoking, snoring, study_habit, remarks)
SELECT u.id, ss.sleep_time, ss.smoking, ss.snoring, ss.study_habit, ss.remarks
FROM seed_surveys ss
JOIN users u ON u.username = ss.username;

COMMIT;

