package color

import (
	"log"

	"github.com/devinmcgloin/clr/clr"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ColorCatagory string

const Shade = "shade"
const SpecificColor = "specific"

type SpriocColorTable struct {
	db   *sqlx.DB
	Type ColorCatagory
}

func (spc SpriocColorTable) Iterate() []clr.Color {
	hexColors := []struct{ Hex string }{}
	err := spc.db.Select(&hexColors, `SELECT hex FROM colors.clr WHERE type = $1`, spc.Type)

	if err != nil {
		log.Println(err)
		return []clr.Color{}
	}

	colors := make([]clr.Color, len(hexColors))
	for i, color := range hexColors {
		colors[i] = clr.Hex{Code: color.Hex}
	}

	return colors
}

func (spc SpriocColorTable) Lookup(hex string) clr.ColorSpace {
	var name string
	err := spc.db.Get(&name, `SELECT name FROM colors.clr WHERE hex = $1`, hex)
	if err != nil {
		log.Println(err)
		return ""
	}
	return clr.ColorSpace(name)
}

func RetrieveColorTable(db *sqlx.DB, t ColorCatagory) SpriocColorTable {
	return SpriocColorTable{
		db:   db,
		Type: t,
	}
}

func AddColor(db *sqlx.DB, name, hex, t string) error {
	_, err := db.Exec("INSERT INTO colors.clr(name, hex, type) VALUES($1, $2, $3)",
		name, hex, t)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func AddColors(db *sqlx.DB, colors map[string]string, t string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare(pq.CopyInSchema("colors", "clr", "name", "hex", "type"))
	if err != nil {
		log.Println(err)
		return err
	}

	for hex, name := range colors {
		_, err = stmt.Exec(name, hex, t)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Println(err)
		return err
	}

	err = stmt.Close()
	if err != nil {
		log.Println(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetColors(db *sqlx.DB, t string) (map[string]string, error) {
	var clrs map[string]string
	var hex []struct {
		Name string
		Hex  string
	}
	err := db.Select(&hex, "SELECT name, hex FROM colors.clr WHERE type = $1", t)
	if err != nil {
		log.Println(err)
		return clrs, err
	}

	for _, color := range hex {
		clrs[color.Hex] = color.Name
	}
	return clrs, nil
}