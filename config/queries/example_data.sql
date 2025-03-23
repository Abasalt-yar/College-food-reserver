DELETE FROM students;
ALTER SEQUENCE students_id_seq RESTART WITH 1;

DELETE FROM foods;
ALTER SEQUENCE foods_id_seq RESTART WITH 1;

DELETE FROM admins;
ALTER SEQUENCE admins_id_seq RESTART WITH 1;

DELETE FROM colleges;
ALTER SEQUENCE colleges_id_seq RESTART WITH 1;

INSERT INTO colleges (college_name) VALUES ('دانشگاه فنی و حرفه ای شماره یک');
INSERT INTO admins (username,password,first_name,last_name,position,college_id) VALUES ('imadminboy','$2a$10$Y3JC7A24z7pFZkKd1Fz/c.ygLbNPWHkUYzWDdd6IsO29TJ5ISutW.','ربات','رباتی','COLLEGE_OWNER',1);
INSERT INTO college_locations (name,college_id) VALUES ('سلف یک',1);
INSERT INTO college_locations (name,college_id) VALUES ('سلف دو',1);
INSERT INTO students (username,password,first_name,last_name,phone_number,position,college_id,added_by) VALUES ('03121095705031','$2a$10$Y3JC7A24z7pFZkKd1Fz/c.ygLbNPWHkUYzWDdd6IsO29TJ5ISutW.','اباصلت','یارمحمدزئی','09101234567','NORMAL_STUDENT',1,1) , ('03121095705032','$2a$10$Y3JC7A24z7pFZkKd1Fz/c.ygLbNPWHkUYzWDdd6IsO29TJ5ISutW.','علی','علیزاده','09101234567','DORM_STUDENT',1,1);
-- Password = 12345678

INSERT INTO foods (food_name,college_id,added_by,expensive) VALUES ('غذای یک',1,1,true), ('غذای دو',1,1,false), ('غذای سه',1,1,true), ('غذای چهار',1,1,false)