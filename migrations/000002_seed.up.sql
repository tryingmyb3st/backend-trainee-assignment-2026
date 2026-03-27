INSERT INTO rooms(id, name, description, capacity)
VALUES
('c81753a7-a230-43ba-a6a6-3b08f9c9913e', 'meeting_room_1', 'small meeting room', 4),
('c9d4be37-f36f-4222-bca4-e7d4fa73f61c', 'meeting_room_2', 'room for team meetings', 6),
('c25256b8-cab7-4799-a332-369dc7376034', 'conference_room', 'spacious conference room', 12),
('95229adc-0b13-44d8-85c6-121d1adfea79', 'quiet_room', 'quiet space for focus', 1),
('a99c7dd2-73b3-4fcc-bc52-e8492acce854', 'work_room', 'room for daily work', 3);

INSERT INTO schedules(id, room_id, days_of_week, start_time, end_time)
VALUES
('bbf6db6d-cabb-48c6-adeb-7634c96a3681', 'c81753a7-a230-43ba-a6a6-3b08f9c9913e', ARRAY[1,2,3,4,5,7], '13:00:00', '14:00:00'),
('46e44f99-c2e2-4a98-b677-1476b6723f09', 'c9d4be37-f36f-4222-bca4-e7d4fa73f61c', ARRAY[1,3,5,7], '16:20:00', '17:40:00');

INSERT INTO slots(id, room_id, start_timestamp, end_timestamp)
VALUES
('37106485-a6c2-4b12-b6bd-f0e76eda294c', 'c9d4be37-f36f-4222-bca4-e7d4fa73f61c', '2026-03-29T16:20:00Z', '2026-03-29T16:50:00Z'),
('bc10a525-8544-4666-b781-26565afd0aac', 'c9d4be37-f36f-4222-bca4-e7d4fa73f61c', '2026-03-29T16:50:00Z', '2026-03-29T17:20:00Z');

INSERT INTO bookings(id, slot_id, user_id, status)
VALUES
('67f2bc09-edaa-43e8-a1e8-25e3f345e25f', '37106485-a6c2-4b12-b6bd-f0e76eda294c', '8794e589-0ddb-43ce-9f92-16faafcf4ee4', 'cancelled'),
('9b3b9cc8-3650-42b3-acd0-1ce7fd9dee27', 'bc10a525-8544-4666-b781-26565afd0aac', '8794e589-0ddb-43ce-9f92-16faafcf4ee4', 'active');