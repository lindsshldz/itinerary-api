DROP DATABASE IF EXISTS itinerary_app;

CREATE DATABASE itinerary_app;

USE itinerary_app;

CREATE TABLE trips (
    id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    trip_name varchar(255) NOT NULL,
    budget float(20,2),
    start_date date NOT NULL,
    end_date date NOT NULL
);

CREATE TABLE details (
    id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    date date NOT NULL,
    day_location varchar(255),
    activities varchar(255),
    restaurants varchar(255),
    hotel varchar(255),
    trip_id int,

    FOREIGN KEY (trip_id) REFERENCES trips(id)
);

INSERT INTO trips (id, trip_name, budget, start_date, end_date) VALUES (1, 'Japan', 10000.75, '2018-08-10', '2018-08-12');
INSERT INTO trips (id, trip_name, budget, start_date, end_date) VALUES (2, 'Tahiti', 5000.50, '2018-12-31', '2019-01-02');
INSERT INTO trips (id, trip_name, budget, start_date, end_date) VALUES (3, 'Vegas', 15000.00, '2018-10-01', '2018-10-03');


INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (1, '2018-08-10', 'tokyo', 'play vr games and dress up like hello kitty', 'kfc', 'hotel-san', 1);
INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (2, '2018-08-11', 'kyoto', 'visit temples and have tea ceremony', 'ramen shop', 'futon palace', 1);
INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (3, '2018-08-12', 'hiroshima', 'visit peace park', 'gyoza and pan for my tumble', 'domo domo in my tumo', 1);

INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (4, '2018-12-31', 'papenoa', 'beach it', 'shaved ice', 'fancy-huts-r-us', 2);
INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (5, '2019-01-01', 'papenoa',  'scuba it', 'roasted piggly wiggly', 'fancy-huts-r-us', 2);
INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (6, '2019-01-02', 'papenoa', 'hike it', 'fresh fishies', 'fancy-huts-r-us', 2);

INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (7, '2018-10-01', 'strip', 'gamble and zombie burlesque', 'adult taco bell', 'aria', 3);
INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (8, '2018-10-02', 'red rocks', 'hike and explore', 'shake shack', 'aria', 3);
INSERT INTO details (id, date, day_location, activities, restaurants, hotel, trip_id) VALUES (9, '2018-10-03', 'strip', 'cry three times', 'momofuku', 'aria', 3);
