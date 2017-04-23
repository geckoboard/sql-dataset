CREATE TABLE builds (id INTEGER PRIMARY KEY, build_id INTEGER, build_cost FLOAT, app_name TEXT, app_sha TEXT, triggered_by TEXT, created_at TIMESTAMP NOT NULL, updated_at TIMESTAMP);

INSERT INTO builds(id, build_id, build_cost, app_name, app_sha, triggered_by, created_at, updated_at) VALUES(1, 420,0.54, "everdeen", "12345abcd", "maeve millay", "2017-03-21 11:12:00", "2017-04-23 11:14:00");
INSERT INTO builds(id, build_id, build_cost, app_name, app_sha, triggered_by, created_at, updated_at) VALUES(2, 563,1.11, "react", "23453abcd", "dr.robert ford", "2017-04-23 12:32:00", "2017-04-23 13:32:00");
INSERT INTO builds(id, build_id, build_cost, app_name, app_sha, triggered_by, created_at, updated_at) VALUES(3, 673,0.24, "geckoboard-ruby", "457890abcd", "bernard", "2017-04-23 13:42:00", "2017-04-23 13:42:00");

INSERT INTO builds(id, build_id, build_cost, app_name, app_sha, triggered_by, created_at, updated_at) VALUES(4, 123,2.64, "westworld", "654321abcd", "dolores", "2017-03-23 15:11:00", NULL);
INSERT INTO builds(id, build_id, build_cost, app_name, app_sha, triggered_by, created_at, updated_at) VALUES(5, NULL,0.24, "geckoboard-ruby", "457890abcd", "bernard", "2017-03-23 16:12:00", NULL);
INSERT INTO builds(id, build_id, build_cost, app_name, app_sha, triggered_by, created_at, updated_at) VALUES(6, 443,NULL, "geckoboard-ruby", "457890abcd", "bernard", "2017-03-23 16:22:00", "2017-03-23 16:55:00");
INSERT INTO builds(id, build_id, build_cost, app_name, app_sha, triggered_by, created_at, updated_at) VALUES(7, 943,11.32, "", "", "hector", "2017-03-23 16:44:00", "2017-03-23 16:45:00");
