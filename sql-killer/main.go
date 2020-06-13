package main

import (
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	user, pwd, dbName, host, port *string
	timeout *int
)

type SlowSearch struct {
	Id int
}

func main(){
	initConfig()
	flag.Parse()
	checkConfig()

	// 连接数据库
	setting := fmt.Sprintf(
		"%s:%s@/%s?%s:%scharset=utf8&parseTime=True&loc=Local", *user, *pwd, *dbName, *host, *port)
	db, err := gorm.Open("mysql", setting)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	// 扫描慢查询
	fmt.Println("Start Job")
	var ss []SlowSearch

	db.Raw("SELECT id FROM information_schema.processlist WHERE Command != 'Sleep' AND Time > ?;", *timeout).Scan(&ss)
	if len(ss) == 0 {
		fmt.Println("Slow search is not found!")
	}

	for _, s := range ss{
		fmt.Printf("Kill %d", s.Id)
		if row := db.Exec("Kill ?;", s.Id).RowsAffected; row == 0{
			fmt.Println(" failed")
		}else{
			fmt.Println(" success")
		}
	}
	fmt.Println("Finished")
}

func initConfig(){
	user = flag.String("user", "", "database username")
	pwd = flag.String("pwd", "", "database password")
	dbName = flag.String("dbname", "", "database name")
	port = flag.String("port", "3306", "database port (default: 3306)")
	host = flag.String("host", "127.0.0.1", "database host (default: 127.0.0.1)")
	timeout = flag.Int("timeout", 0, "sql timeout")
}

func checkConfig(){
	if *user == "" ||  *pwd == "" || *dbName == "" || *timeout == 0{
		panic("Have not provide username or password or dbname or timeout")
	}
}

