DROP TABLE IF EXISTS builds;
CREATE TABLE builds(id INT PRIMARY KEY, build_cost FLOAT, percent_passed INT, run_time FLOAT, app_name VARCHAR(50), triggered_by VARCHAR(255), created_at DATETIME NOT NULL);

INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(1, 0.54, 80, 0.31882276212,  'everdeen', 'maeve millay', '2017-03-21 11:12:00');
INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(2, 1.11, 95, 118.18382961212, 'react', 'dr.robert ford', '2017-04-23 12:32:00');
INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(3, 0.24, 24, 0.21882232124, 'geckoboard-ruby', 'maeve millay', '2017-04-23 13:42:00');
INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(4, 1.44, 100, 144.31838122382, 'everdeen', 'bernard', '2017-03-21 11:13:00');
INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(5, 0.92, 55, 77.21381276421, 'geckoboard-ruby', 'bernard', '2017-04-23 13:43:00');

INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(6, 2.64, NULL, 321.93774373, 'westworld', 'dolores', '2017-03-23 15:11:00');
INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(7, NULL, NULL, NULL, 'geckoboard-ruby', 'bernard', '2017-03-23 16:12:00');
INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(8, NULL, 1, 0.12349876543, '', 'dr.robert ford', '2017-03-23 16:22:00');
INSERT INTO builds(id, build_cost, percent_passed, run_time, app_name, triggered_by, created_at) VALUES(9, 11.32, 34, 46.432763287, '', 'hector', '2017-03-23 16:44:00')
