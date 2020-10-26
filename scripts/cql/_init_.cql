-- noinspection SqlNoDataSourceInspectionForFile

CREATE KEYSPACE lifetrack_dev
  WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3};

USE lifetrack_dev;

CREATE TABLE category (
  id text,
  user_id text,
  name text,
  description text,
  target_time int,
  picture text,
  create_time timestamp,
  update_time timestamp,
  active boolean,
  PRIMARY KEY(id, user_id)
);

CREATE TABLE activity (
  id text,
  category_id text,
  title text,
  picture text,
  create_time timestamp,
  update_time timestamp,
  active boolean,
  PRIMARY KEY(id, category_id)
);

CREATE TABLE occurrence (
  id text,
  activity_id text,
  start_time timestamp,
  end_time timestamp,
  total_duration int,
  create_time timestamp,
  update_time timestamp,
  PRIMARY KEY ((id, activity_id), start_time)
);
