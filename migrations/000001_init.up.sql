CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(100) NOT NULL UNIQUE CHECK (email ~ '^[a-zA-Z0-9.]+@[a-zA-Z0-9.]+\.[a-zA-Z]{2,}$'),
  password TEXT,
  role VARCHAR(15) NOT NULL CHECK (role IN ('admin', 'user')),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE rooms (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  capacity INTEGER CHECK (capacity >= 0),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE schedules (
  id VARCHAR(36) PRIMARY KEY,
  room_id VARCHAR(36) NOT NULL UNIQUE REFERENCES rooms(id),
  days_of_week INTEGER[] NOT NULL CHECK (1 <= ALL(days_of_week) AND 7 >= ALL(days_of_week)),
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,

  CHECK (start_time < end_time)
);

CREATE TABLE slots (
  id VARCHAR(36) PRIMARY KEY,
  room_id VARCHAR(36) NOT NULL REFERENCES rooms(id),
  start_timestamp TIMESTAMP NOT NULL,
  end_timestamp TIMESTAMP NOT NULL,

  CHECK (start_timestamp < end_timestamp)
);

CREATE TABLE bookings (
  id VARCHAR(36) PRIMARY KEY,
  slot_id VARCHAR(36) NOT NULL UNIQUE REFERENCES slots(id),
  user_id VARCHAR(36) NOT NULL REFERENCES users(id),
  status VARCHAR(15) NOT NULL CHECK(status IN ('active', 'cancelled')),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_schedules_room_id ON schedules(room_id);
CREATE INDEX idx_slots_room_id ON slots(room_id);
CREATE INDEX idx_bookings_slot_id ON bookings(slot_id);
CREATE INDEX idx_bookings_user_id ON bookings(user_id);

INSERT INTO users(id, email, role)
VALUES('8794e589-0ddb-43ce-9f92-16faafcf4ee4', 'test.user@email.com', 'user');

INSERT INTO users(id, email, role)
VALUES('249be7cf-d419-4c54-97f2-d04107806e36', 'test.admin@email.com', 'admin');