package customers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	//"github.com/nareenuch/finalexam/customers"
	"github.com/nareenuch/finalexam/database"

	_ "github.com/lib/pq"
)

type response struct {
    Message   string `json:"message"`
}

func createCustomersHandler(c *gin.Context){
	var item Customers
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bind error" : err.Error()})
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"connect database error : " : err.Error()})
		return
	}
	// db := database.Conn()
	defer db.Close()
	row := db.QueryRow("INSERT INTO customers (name, email, status) values ($1, $2, $3) RETURNING id", item.Name, item.Email, item.Status )
	err = row.Scan(&item.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"can't scan id " : err, "id " : item.ID})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func getCustomersByIdHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"connect database error : " : err.Error()})
		return
	}
	// db := database.Conn()
	defer db.Close()
	stmt, err := db.Prepare("SELECT id, name, email, status FROM customers Where id = $1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "can't prepare query statement")
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	row := stmt.QueryRow(id)
	var myid int
	var myname, myemail, mystatus string
	err = row.Scan(&myid, &myname, &myemail, &mystatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "can't scan row into variable")
		return
		//fmt.Println(c.Param("id"))
		//log.Fatal("can't scan row into variable", err)
	}
	t := Customers{ID: myid, Name: myname, Email: myemail, Status: mystatus,}
	c.JSON(http.StatusOK, t)
}


func getCustomersHandler(c *gin.Context){
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"connect database error : " : err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM customers")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error Prepare :" : err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError,"error query step")
		return
	}

	cust := []Customers{}
	for rows.Next() {
		t := Customers{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Error scan row")
			return
		}
		t = Customers{ID: t.ID, Name: t.Name, Status: t.Status,}
		cust = append(cust, t)
	}
	c.JSON(http.StatusOK, cust)
}


func updateCustomersHandler(c *gin.Context) {
	var item Customers
	id, _ := strconv.Atoi(c.Param("id"))
	err := c.ShouldBindJSON(&item)
	item.ID = id
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "connect to database error")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "can't prepare statement update")
		return
	}

	row:=stmt.QueryRow(item.ID,item.Name,item.Email,item.Status)
	if err != nil {
	 	c.JSON(http.StatusInternalServerError, "error execute update")
	 	return
	}

	c.JSON(http.StatusOK, item)
	fmt.Println(row)
}

func deleteCustomersHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "connect to database error")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM customers WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, "can't prepare statement delete")
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	item, err := stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "execute error")
		return
	}
	fmt.Println("delete success ", item)
	resp := response{Message:"customer deleted"}
	c.JSON(http.StatusOK, resp)
}

func CreateTable() {
	ctb := `
	CREATE TABLE IF NOT EXISTS customers (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
		);`

	_, err := database.Conn().Exec(ctb)
	if err != nil {
		log.Fatal("can't create table", err)
	}
}

func loginMiddleware(c *gin.Context) {
	//log.Println("starting middleware")
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Abort()
		return
	}
	c.Next()
	//log.Println("ending middleware")
}


func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(loginMiddleware)
	r.GET("/customers", getCustomersHandler)
	r.GET("/customers/:id", getCustomersByIdHandler)
	r.PUT("/customers/:id", updateCustomersHandler) 
	r.POST("/customers", createCustomersHandler)
	r.DELETE("/customers/:id", deleteCustomersHandler)
	return r
}
