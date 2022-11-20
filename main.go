package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"log"
	"net/http"
)

type Company struct {
	ID uuid.UUID `json:"id", field:"id"`
	// Use pointers to distinguish between null
	// and empty/zero/default values in input data
	// For example, a company that is registered: false
	// vs registered: null
	Name              *string `json:"name", db:"name"`
	Description       *string `json:"description", db:"description"`
	AmountOfEmployees *int    `json:"amount_of_employees", db:"amount_of_employees"`
	Registered        *bool   `json:"registered", db:"registered"`
	Type              *string `json:"company_type", db:"company_type"`
}

func NewCompany() *Company {
	var s1 string
	var s2 string
	var s3 string
	var b bool
	var i int
	return &Company{
		Name:              &s1,
		Description:       &s2,
		AmountOfEmployees: &i,
		Registered:        &b,
		Type:              &s3,
	}
}

const dsn = "postgres://postgres:xm@localhost:5432/xm"

type api struct {
	r  *gin.Engine
	db *sql.DB
}

func main() {
	fmt.Println("XM Golang test")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("could not open db: %s", err.Error())
	}

	a := api{
		r:  gin.Default(),
		db: db,
	}

	a.r.GET("/:id", a.ReadCompany)
	auth := a.r.Group("/", gin.BasicAuth(gin.Accounts{
		"admin": "xm",
	}))
	auth.POST("/", a.CreateCompany)
	auth.PATCH("/:id", a.UpdateCompany)
	auth.DELETE("/:id", a.DeleteCompany)

	a.r.Run()
}

func (a *api) CreateCompany(c *gin.Context) {
	data := &Company{}

	if err := c.BindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("error unmarshaling data: %s", err.Error()),
			})
		log.Println(err.Error())
		return
	}

	if err := checkFields(data); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("error with data: %s", err.Error()),
			})
		log.Println(err.Error())
		return
	}

	data.ID = uuid.New()

	execStr := "INSERT INTO companies (id, name, description, amount_of_employees, registered, company_type)" +
		"VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := a.db.Exec(execStr, data.ID, data.Name, data.Description, data.AmountOfEmployees,
		data.Registered, data.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("error inserting database: %s", err.Error()),
			})
		return
	}

	c.JSON(http.StatusOK,
		gin.H{
			"id":   data.ID,
			"name": data.Name,
		})
}

func (a *api) ReadCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("could not parse ID: %s", err.Error()),
			})
		return
	}

	row := a.db.QueryRow("select * from companies where id = $1", id.String())

	data := NewCompany()
	var desc sql.NullString
	err = row.Scan(&data.ID, data.Name, &desc, data.AmountOfEmployees,
		data.Registered, data.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("error parsing database query: %s", err.Error()),
			})
		return
	}

	if desc.Valid {
		data.Description = &desc.String
	}

	c.JSON(http.StatusOK, data)
}

func (a *api) UpdateCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("could not parse ID: %s", err.Error()),
			})
		return
	}

	data := &Company{}
	if err := c.BindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("error unmarshaling data: %s", err.Error()),
			})
		return
	}

	execStr := "UPDATE companies SET (id, name, description, amount_of_employees, registered, company_type) =" +
		"($1, $2, $3, $4, $5, $6) WHERE id = $7"
	_, err = a.db.Exec(execStr, id, data.Name, data.Description, data.AmountOfEmployees, data.Registered,
		data.Type, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("error updating database: %s", err.Error()),
			})
		return
	}

	c.String(http.StatusOK, "account updated")
}

func (a *api) DeleteCompany(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{
				"error": fmt.Sprintf("could not parse ID: %s", err.Error()),
			})
		return
	}

	_, err = a.db.Exec("DELETE FROM companies WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("error updating database: %s", err.Error()),
			})
		return
	}

	c.String(http.StatusOK, "account deleted")
}

func checkFields(c *Company) error {
	if c.Name == nil {
		return errors.New("'name' should not be null")
	}
	if c.AmountOfEmployees == nil {
		return errors.New("'amount_of_employees' should not be null")
	}
	if c.Registered == nil {
		return errors.New("'registered' should not be null")
	}
	if c.Type == nil {
		return errors.New("'company_type' should not be null")
	}
	valid := map[string]bool{
		"Corporations":        true,
		"NonProfit":           true,
		"Cooperative":         true,
		"Sole Proprietorship": true,
	}
	if !valid[*c.Type] {
		return errors.New("invalid company type")
	}
	return nil
}
