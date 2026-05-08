-- Dormitory management system test data
-- Usage:
--   PGPASSWORD=10928 psql -U turing -d student_dormitory -f sql/002_seed_student_dormitory_test_data.sql
--
-- Notes:
--   Default password for all seeded users: 123456
--   This script intentionally does not create repair/cleaning work orders or image attachments.

BEGIN;

WITH seed_users(username, role, name, phone) AS (
    VALUES
        ('admin001', 'system_admin', '系统管理员', '13000000001'),
        ('manager001', 'dormitory_manager', '宿舍管理员王老师', '13000000002'),
        ('repair001', 'repair_staff', '维修人员李师傅', '13000000003'),
        ('cleaner001', 'cleaning_staff', '保洁人员赵阿姨', '13000000004'),
        ('student001', 'student', '张三', '13800000001'),
        ('student002', 'student', '李四', '13800000002'),
        ('student003', 'student', '王五', '13800000003'),
        ('student004', 'student', '赵六', '13800000004'),
        ('student005', 'student', '钱七', '13800000005'),
        ('student006', 'student', '孙八', '13800000006'),
        ('student007', 'student', '周九', '13800000007'),
        ('student008', 'student', '吴十', '13800000008')
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

WITH seed_buildings(name, location) AS (
    VALUES
        ('一号宿舍楼', '东区生活区'),
        ('二号宿舍楼', '西区生活区'),
        ('留学生公寓', '南门附近')
)
INSERT INTO buildings (name, location)
SELECT name, location
FROM seed_buildings
WHERE NOT EXISTS (
    SELECT 1 FROM buildings b WHERE b.name = seed_buildings.name
);

WITH seed_rooms(building_name, room_number, floor, total_beds, water_balance, electricity_balance) AS (
    VALUES
        ('一号宿舍楼', '101', 1, 4, 28.50, 36.00),
        ('一号宿舍楼', '102', 1, 4, 3.20, 18.00),
        ('一号宿舍楼', '201', 2, 4, 15.00, 4.50),
        ('一号宿舍楼', '202', 2, 4, 45.00, 60.00),
        ('二号宿舍楼', '301', 3, 4, 20.00, 22.00),
        ('二号宿舍楼', '302', 3, 4, 4.80, 3.90),
        ('留学生公寓', '501', 5, 2, 50.00, 50.00)
)
INSERT INTO rooms (building_id, room_number, floor, total_beds, water_balance, electricity_balance)
SELECT b.id, sr.room_number, sr.floor, sr.total_beds, sr.water_balance, sr.electricity_balance
FROM seed_rooms sr
JOIN buildings b ON b.name = sr.building_name
ON CONFLICT (building_id, room_number) DO UPDATE
SET floor = EXCLUDED.floor,
    total_beds = EXCLUDED.total_beds,
    water_balance = EXCLUDED.water_balance,
    electricity_balance = EXCLUDED.electricity_balance;

WITH labels(label) AS (
    VALUES ('A'), ('B'), ('C'), ('D')
),
rooms_for_beds AS (
    SELECT r.id AS room_id, r.total_beds
    FROM rooms r
    JOIN buildings b ON b.id = r.building_id
    WHERE b.name IN ('一号宿舍楼', '二号宿舍楼')
)
INSERT INTO beds (room_id, bed_label, status)
SELECT room_id, label, 'available'
FROM rooms_for_beds
JOIN labels ON true
WHERE
    (total_beds >= 4)
    AND NOT EXISTS (
        SELECT 1 FROM beds existing
        WHERE existing.room_id = rooms_for_beds.room_id
          AND existing.bed_label = labels.label
    );

WITH labels(label) AS (
    VALUES ('A'), ('B')
),
rooms_for_beds AS (
    SELECT r.id AS room_id
    FROM rooms r
    JOIN buildings b ON b.id = r.building_id
    WHERE b.name = '留学生公寓'
      AND r.room_number = '501'
)
INSERT INTO beds (room_id, bed_label, status)
SELECT room_id, label, 'available'
FROM rooms_for_beds
JOIN labels ON true
WHERE NOT EXISTS (
    SELECT 1 FROM beds existing
    WHERE existing.room_id = rooms_for_beds.room_id
      AND existing.bed_label = labels.label
);

WITH assignments(username, building_name, room_number, bed_label, occupied_since) AS (
    VALUES
        ('student001', '一号宿舍楼', '101', 'A', DATE '2025-09-01'),
        ('student002', '一号宿舍楼', '101', 'B', DATE '2025-09-01'),
        ('student003', '一号宿舍楼', '101', 'C', DATE '2025-09-01'),
        ('student004', '一号宿舍楼', '102', 'A', DATE '2025-09-02'),
        ('student005', '一号宿舍楼', '102', 'B', DATE '2025-09-02'),
        ('student006', '一号宿舍楼', '201', 'A', DATE '2025-09-03'),
        ('student007', '二号宿舍楼', '301', 'A', DATE '2025-09-03')
),
resolved AS (
    SELECT
        bed.id AS bed_id,
        u.id AS student_id,
        a.occupied_since
    FROM assignments a
    JOIN users u ON u.username = a.username
    JOIN buildings b ON b.name = a.building_name
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
        ('student001', TIME '23:00', 0, 0, '晚上学习，偏安静', '希望室友作息规律'),
        ('student002', TIME '23:30', 0, 1, '白天学习为主', '可接受轻微噪音'),
        ('student003', TIME '00:00', 0, 0, '夜间学习较多', '希望靠窗床位'),
        ('student004', TIME '22:30', 0, 0, '早睡早起', '对卫生要求较高'),
        ('student008', TIME '23:45', 0, 0, '弹性作息', '新生待分配')
)
INSERT INTO lifestyle_surveys (student_id, sleep_time, smoking, snoring, study_habit, remarks)
SELECT u.id, ss.sleep_time, ss.smoking, ss.snoring, ss.study_habit, ss.remarks
FROM seed_surveys ss
JOIN users u ON u.username = ss.username
WHERE NOT EXISTS (
    SELECT 1 FROM lifestyle_surveys existing
    WHERE existing.student_id = u.id
      AND existing.study_habit = ss.study_habit
);

WITH target AS (
    SELECT
        u.id AS student_id,
        r.id AS room_id,
        bed.id AS bed_id
    FROM users u
    JOIN buildings b ON b.name = '二号宿舍楼'
    JOIN rooms r ON r.building_id = b.id AND r.room_number = '302'
    JOIN beds bed ON bed.room_id = r.id AND bed.bed_label = 'A'
    WHERE u.username = 'student008'
)
INSERT INTO allocation_requests (student_id, recommended_room_id, recommended_bed_id, status)
SELECT student_id, room_id, bed_id, 'pending'
FROM target
WHERE NOT EXISTS (
    SELECT 1 FROM allocation_requests ar
    WHERE ar.student_id = target.student_id
      AND ar.status = 'pending'
);

WITH student AS (
    SELECT id FROM users WHERE username = 'student001'
),
manager AS (
    SELECT id FROM users WHERE username = 'manager001'
)
INSERT INTO leave_applications (student_id, type, destination, emergency_contact, return_time, reason, status, manager_id)
SELECT student.id, 'normal', '上海市浦东新区', '父亲 13900000001', now() + interval '3 days', '周末回家', 'approved', manager.id
FROM student, manager
WHERE NOT EXISTS (
    SELECT 1 FROM leave_applications la
    WHERE la.student_id = student.id
      AND la.reason = '周末回家'
);

WITH student AS (
    SELECT id FROM users WHERE username = 'student004'
)
INSERT INTO leave_applications (student_id, type, destination, emergency_contact, return_time, reason)
SELECT student.id, 'holiday', '南京市', '母亲 13900000004', now() + interval '7 days', '节假日离校'
FROM student
WHERE NOT EXISTS (
    SELECT 1 FROM leave_applications la
    WHERE la.student_id = student.id
      AND la.reason = '节假日离校'
);

WITH student AS (
    SELECT id FROM users WHERE username = 'student002'
)
INSERT INTO late_return_records (student_id, return_date, reason)
SELECT student.id, CURRENT_DATE, '社团活动结束较晚'
FROM student
WHERE NOT EXISTS (
    SELECT 1 FROM late_return_records lr
    WHERE lr.student_id = student.id
      AND lr.return_date = CURRENT_DATE
      AND lr.reason = '社团活动结束较晚'
);

WITH student_bed AS (
    SELECT
        u.id AS student_id,
        from_bed.id AS from_bed_id,
        target_room.id AS target_room_id,
        target_bed.id AS target_bed_id
    FROM users u
    JOIN beds from_bed ON from_bed.student_id = u.id
    JOIN rooms from_room ON from_room.id = from_bed.room_id
    JOIN rooms target_room ON target_room.building_id = from_room.building_id AND target_room.room_number = '202'
    JOIN beds target_bed ON target_bed.room_id = target_room.id AND target_bed.bed_label = 'A'
    WHERE u.username = 'student003'
      AND target_bed.status = 'available'
)
INSERT INTO room_change_requests (student_id, from_bed_id, target_room_id, target_bed_id, reason)
SELECT student_id, from_bed_id, target_room_id, target_bed_id, '希望与同专业同学同住'
FROM student_bed
WHERE NOT EXISTS (
    SELECT 1 FROM room_change_requests rcr
    WHERE rcr.student_id = student_bed.student_id
      AND rcr.status = 'pending'
);

WITH student AS (
    SELECT id FROM users WHERE username = 'student006'
)
INSERT INTO off_campus_living_applications (student_id, retain_bed, reason, destination)
SELECT student.id, 1, '短期在校外实习，申请保留床位', '学校附近实习单位宿舍'
FROM student
WHERE NOT EXISTS (
    SELECT 1 FROM off_campus_living_applications app
    WHERE app.student_id = student.id
      AND app.status = 'pending'
);

WITH room_payer AS (
    SELECT r.id AS room_id, u.id AS payer_id
    FROM users u
    JOIN beds bed ON bed.student_id = u.id
    JOIN rooms r ON r.id = bed.room_id
    WHERE u.username = 'student001'
)
INSERT INTO utility_payments (room_id, payer_id, amount, payment_type)
SELECT room_id, payer_id, 30.00, 'both'
FROM room_payer
WHERE NOT EXISTS (
    SELECT 1 FROM utility_payments p
    WHERE p.payer_id = room_payer.payer_id
      AND p.amount = 30.00
      AND p.payment_type = 'both'
);

UPDATE rooms r
SET water_balance = 43.50,
    electricity_balance = 51.00
FROM buildings b
WHERE b.id = r.building_id
  AND b.name = '一号宿舍楼'
  AND r.room_number = '101';

WITH student AS (
    SELECT id FROM users WHERE username = 'student001'
)
INSERT INTO notifications (recipient_id, message, type, is_read)
SELECT student.id, '欢迎使用宿舍管理系统，请及时完善个人信息。', 'general', 0
FROM student
WHERE NOT EXISTS (
    SELECT 1 FROM notifications n
    WHERE n.recipient_id = student.id
      AND n.message = '欢迎使用宿舍管理系统，请及时完善个人信息。'
);

COMMIT;
