-- project: gorchat
-- author: wens
-- date: 2025/07/09

-- create im
create database im;

-- create 'user'@'%' as defualt remote access user
create user 'gochat'@'%' identified by 'gochat';
grant all privileges on im.* to 'gochat'@'%';
flush privileges;