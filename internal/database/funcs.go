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
    coord_id INTEGER,
	FOREIGN KEY (coord_id) REFERENCES Coordinates(id) ON DELETE CASCADE,
    base VARCHAR(50),
    main_id INTEGER,
	FOREIGN KEY (main_id) REFERENCES Main(id) ON DELETE CASCADE,
    visibility INTEGER,
    wind_id INTEGER,
	FOREIGN KEY (wind_id) REFERENCES Wind(id) ON DELETE CASCADE,
    clouds_id INTEGER,
	FOREIGN KEY (clouds_id) REFERENCES Clouds(id) ON DELETE CASCADE,
    dt INTEGER,
    sys_id INTEGER,
	FOREIGN KEY (sys_id) REFERENCES Sys(id) ON DELETE CASCADE,
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

func (m *Manager) GetWeather(city string) (structs.WeatherData, error) { //this is a mess
	var ret structs.WeatherData
	query := fmt.Sprintf(`
	SELECT
	wd.base,
	wd.visibility,
	wd.dt,
	wd.timezone,
	wd.name,
	wd.cod,
	c.lon,
	c.lat,
	m.temp,
	m.feels_like,
	m.temp_min,
	m.temp_max,
	m.pressure,
	m.humidity,
	m.sea_level,
	m.grnd_level,
	wn.speed,
	wn.deg,
	wn.gust,
	cl.all,
	s.type,
	s.country,
	s.sunrise,
	s.sunset
		FROM WeatherData wd
		JOIN Coordinates c ON wd.coord_id = c.id
		JOIN Main m ON wd.main_id = m.id
		JOIN Wind wn ON wd.wind_id = wn.id
		JOIN Clouds cl ON wd.clouds_id = cl.id
		JOIN Sys s ON wd.sys_id = s.id
		WHERE wd.city = '%s'; -- assuming you are filtering by WeatherData ID`, city)
	row := m.db.QueryRow(query)
	err := row.Scan(
		&ret.Base, &ret.Visibility, &ret.Dt, &ret.Timezone, &ret.Name, &ret.Cod,
		&ret.Coord.Lon, &ret.Coord.Lat,
		&ret.Main.Temp, &ret.Main.FeelsLike, &ret.Main.TempMin, &ret.Main.TempMax, &ret.Main.Pressure, &ret.Main.Humidity, &ret.Main.SeaLevel, &ret.Main.GrndLevel,
		&ret.Wind.Speed, &ret.Wind.Deg, &ret.Wind.Gust,
		&ret.Clouds.All,
		&ret.Sys.Type, &ret.Sys.Country, &ret.Sys.Sunrise, &ret.Sys.Sunset,
	)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (m *Manager) PutWeather(cwd structs.WeatherData, city string) error {
	var city_id int
	query := fmt.Sprintf(`SELECT id FROM weatherdata WHERE city = '%s'`, city)
	err := m.db.QueryRow(query).Scan(&city_id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	var coord_id int
	query = fmt.Sprintf(`INSERT INTO Coordinates (lon, lat) VALUES (%v, %v) RETURNING id`, cwd.Coord.Lon, cwd.Coord.Lat)
	err = m.db.QueryRow(query).Scan(&coord_id)
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

	query = fmt.Sprintf(`INSERT INTO WeatherData (city, coord_id, base, main_id, visibility, wind_id, clouds_id, dt, sys_id, timezone, name, cod)
VALUES ('%v', %v, '%v', %v, %v, %v, %v, %v, %v, %v, '%v', %v);`, city, coord_id, cwd.Base, main_id, cwd.Visibility, wind_id, clouds_id, cwd.Dt, sys_id, cwd.Timezone, cwd.Name, cwd.Cod)
	_, err = m.db.Exec(query)
	if err != nil {
		return err
	}

	m.deleteCityInfo(city_id)

	return nil
}

func (m *Manager) deleteCityInfo(city_id int) error {
	query := fmt.Sprintf(`DELETE FROM weatherdata WHERE id = %d`, city_id)
	_, err := m.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
