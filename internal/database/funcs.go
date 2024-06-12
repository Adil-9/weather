package database

import (
	"database/sql"
	"fmt"
	"weather/internal/structs"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 3000
	user     = "postgres"
	password = "password"
	dbname   = "weather"
)

func ConnectPSQL() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host= %s port= %d user= %s password=%s dbname= %s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	return db, err
}

func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS Coordinates  (
    id SERIAL PRIMARY KEY,
    lon FLOAT,
    lat FLOAT
);

CREATE TABLE IF NOT EXISTS Weather (
    id SERIAL PRIMARY KEY,
    weather_data_id INTEGER,
    main VARCHAR(50),
    description VARCHAR(255),
    icon VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS Main (
    id SERIAL PRIMARY KEY,
    temp FLOAT,
    feels_like FLOAT,
    temp_min FLOAT,
    temp_max FLOAT,
    pressure INTEGER,
    humidity INTEGER,
    sea_level INTEGER,
    grnd_level INTEGER
);

CREATE TABLE IF NOT EXISTS Wind (
    id SERIAL PRIMARY KEY,
    speed FLOAT,
    deg INTEGER,
    gust FLOAT
);

CREATE TABLE IF NOT EXISTS Clouds (
    id SERIAL PRIMARY KEY,
    "all" INTEGER
);

CREATE TABLE IF NOT EXISTS Sys (
    id SERIAL PRIMARY KEY,
    type INTEGER,
    country VARCHAR(5),
    sunrise INTEGER,
    sunset INTEGER
);

CREATE TABLE IF NOT EXISTS WeatherData (
    id SERIAL PRIMARY KEY,
	city VARCHAR(100),
    coord_id INTEGER REFERENCES Coordinates(id),
	weather_id INTEGER REFERENCES Weather(id),
    base VARCHAR(50),
    main_id INTEGER REFERENCES Main(id),
    visibility INTEGER,
    wind_id INTEGER REFERENCES Wind(id),
    clouds_id INTEGER REFERENCES Clouds(id),
    dt INTEGER,
    sys_id INTEGER REFERENCES Sys(id),
    timezone INTEGER,
    name VARCHAR(100),
    cod INTEGER
);
`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetWeather(city string) {

}

func (m *Manager) PutWeather(cwd structs.WeatherData, city string) error {
	var coord_id int
	query := fmt.Sprintf(`INSERT INTO Coordinates (lon, lat) VALUES (%v, %v) RETURNING id`, cwd.Coord.Lon, cwd.Coord.Lat)
	err := m.db.QueryRow(query).Scan(&coord_id)
	if err != nil {
		return err
	}

	var weather_id int
	query = fmt.Sprintf(`INSERT INTO Weather (weather_data_id, main, description, icon) VALUES (%v, '%v', '%v', '%v') RETURNING id`, cwd.Weather[0].ID, cwd.Weather[0].Main, cwd.Weather[0].Description, cwd.Weather[0].Icon)
	err = m.db.QueryRow(query).Scan(&weather_id)
	if err != nil {
		return err
	}

	var main_id int
	query = fmt.Sprintf(`INSERT INTO Main (temp, feels_like, temp_min, temp_max, pressure, humidity, sea_level, grnd_level) VALUES (%v, %v, %v, %v, %v, %v, %v, %v) RETURNING id`, cwd.Main.Temp, cwd.Main.FeelsLike, cwd.Main.TempMin, cwd.Main.TempMax, cwd.Main.Pressure, cwd.Main.Humidity, cwd.Main.SeaLevel, cwd.Main.GrndLevel)
	err = m.db.QueryRow(query).Scan(&main_id)
	if err != nil {
		return err
	}

	var wind_id int
	query = fmt.Sprintf(`INSERT INTO Wind (speed, deg, gust) VALUES (%v, %v, %v) RETURNING id`, cwd.Wind.Speed, cwd.Wind.Deg, cwd.Wind.Gust)
	err = m.db.QueryRow(query).Scan(&wind_id)
	if err != nil {
		return err
	}

	var clouds_id int
	query = fmt.Sprintf(`INSERT INTO Clouds ("all") VALUES (%v) RETURNING id`, cwd.Clouds.All)
	err = m.db.QueryRow(query).Scan(&clouds_id)
	if err != nil {
		return err
	}

	var sys_id int
	query = fmt.Sprintf(`INSERT INTO Sys (type, country, sunrise, sunset) VALUES (%v, '%v', %v, %v) RETURNING id`, cwd.Sys.Type, cwd.Sys.Country, cwd.Sys.Sunrise, cwd.Sys.Sunset)
	err = m.db.QueryRow(query).Scan(&sys_id)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`INSERT INTO WeatherData (city, coord_id, weather_id, base, main_id, visibility, wind_id, clouds_id, dt, sys_id, timezone, name, cod)
VALUES ('%v', %v, %v, '%v', %v, %v, %v, %v, %v, %v, %v, '%v', %v);`, city, coord_id, weather_id, cwd.Base, main_id, cwd.Visibility, wind_id, clouds_id, cwd.Dt, sys_id, cwd.Timezone, cwd.Name, cwd.Cod)
	_, err = m.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
