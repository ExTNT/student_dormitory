export type Role = 'student' | 'repair_staff' | 'cleaning_staff' | 'dormitory_manager' | 'system_admin';

export interface User {
  id: number;
  username: string;
  role: Role;
  name: string;
  phone?: string;
  created_at: string;
  has_survey: boolean;
  has_bed: boolean;
  building_id?: number;
  building_name?: string;
  room_id?: number;
  room_number?: string;
  bed_id?: number;
  bed_label?: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}

export interface Building {
  id: number;
  name: string;
  location?: string;
}

export interface Room {
  id: number;
  building_id: number;
  room_number: string;
  floor: number;
  total_beds: number;
  water_balance: number;
  electricity_balance: number;
}

export interface AvailableBed {
  bed_id: number;
  room_id: number;
  room_number: string;
  bed_label: string;
  building_id: number;
  building_name: string;
  floor: number;
}

export interface LifestyleSurvey {
  id?: number;
  student_id?: number;
  sleep_time?: string;
  smoking: number;
  snoring: number;
  study_habit?: string;
  remarks?: string;
  submitted_at?: string;
}

export interface MyRequest {
  student_id: number;
  request_type: string;
  request_id: number;
  status: string;
  created_at: string;
  detail: string;
}

export interface Roommate {
  student_id: number;
  roommate_id: number;
  roommate_name: string;
  roommate_phone?: string;
  bed_label: string;
  avatar_attachment_id?: number;
}

export interface Notification {
  id: number;
  recipient_id: number;
  room_id?: number;
  message: string;
  type: string;
  is_read: number;
  created_at: string;
}

export interface PendingRepair {
  request_id: number;
  status: string;
  student_id: number;
  student_name: string;
  room_id: number;
  room_number: string;
  description: string;
  repair_staff_id?: number;
  created_at: string;
  repair_staff_name?: string;
  repair_description?: string;
  reviewer_id?: number;
  reviewer_name?: string;
  review_comment?: string;
  accepted_at?: string;
  repaired_at?: string;
  reviewed_at?: string;
}

export interface PendingCleaning {
  request_id: number;
  status: string;
  student_id: number;
  student_name: string;
  building_id: number;
  building_name: string;
  location_desc: string;
  cleaner_id?: number;
  created_at: string;
  cleaner_name?: string;
  reviewer_id?: number;
  reviewer_name?: string;
  review_comment?: string;
  accepted_at?: string;
  cleaned_at?: string;
  reviewed_at?: string;
}

export interface PendingLeave {
  id: number;
  student_id: number;
  student_name: string;
  type: string;
  destination: string;
  emergency_contact: string;
  return_time: string;
  reason: string;
  status: string;
  created_at: string;
}

export interface PendingLateReturn {
  id: number;
  student_id: number;
  student_name: string;
  return_date: string;
  reason: string;
  status: string;
  created_at: string;
}

export interface PendingRoomChange {
  id: number;
  student_id: number;
  student_name: string;
  from_bed_id: number;
  from_building_name: string;
  from_room_number: string;
  from_bed_label: string;
  target_room_id?: number;
  target_bed_id?: number;
  target_building_name?: string;
  target_room_number?: string;
  target_bed_label?: string;
  recommended_bed_id?: number;
  recommended_building_name?: string;
  recommended_room_number?: string;
  recommended_bed_label?: string;
  reason: string;
  status: string;
  created_at: string;
}

export interface PendingOffCampus {
  id: number;
  student_id: number;
  student_name: string;
  retain_bed: number;
  reason: string;
  destination?: string;
  status: string;
  created_at: string;
}

export interface AllocationRequest {
  id: number;
  student_id: number;
  recommended_room_id: number;
  recommended_bed_id: number;
  status: string;
  admin_id?: number;
  created_at: string;
  resolved_at?: string;
}

export interface AttachmentMeta {
  id: number;
  owner_type: string;
  owner_id: number;
  category: string;
  sort_order: number;
  content_type: string;
  file_name?: string;
  uploaded_at: string;
}

export interface DormitorySummary {
  building_id: number;
  building_name: string;
  total_rooms: number;
  total_beds: number;
  occupied_beds: number;
  free_beds: number;
}

export interface LowBalanceRoom {
  room_id: number;
  building_id: number;
  room_number: string;
  water_balance: number;
  electricity_balance: number;
}
