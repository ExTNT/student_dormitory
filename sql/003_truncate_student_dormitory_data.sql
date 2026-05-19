-- Clear all dormitory management system data while keeping schema objects.
-- Usage:
--   PGPASSWORD=passwd psql -U admin -d student_dormitory -f sql/003_truncate_student_dormitory_data.sql
--
-- This script truncates all application tables and restarts identity sequences.
-- Views, functions, triggers, indexes, and table definitions are preserved.

BEGIN;

TRUNCATE TABLE
    refresh_tokens,
    attachments,
    notifications,
    off_campus_living_applications,
    room_change_requests,
    late_return_records,
    cleaning_requests,
    repair_requests,
    utility_payments,
    leave_applications,
    allocation_requests,
    lifestyle_surveys,
    beds,
    rooms,
    buildings,
    users,
    refresh_tokens
RESTART IDENTITY CASCADE;

COMMIT;
