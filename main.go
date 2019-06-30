package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	squadxmlFile = "squad.xml"
	squadxmlLogo = "logo.paa"
	squadxmlHead = `<?xml version="1.0"?>
<!DOCTYPE squad SYSTEM "squad.dtd">
<?xml-stylesheet href="squad.xsl" type="text/xsl"?>`
	squadxmlSquad = `<squad nick="sTs">
<name>Special Tactical Service</name>
<email>noreply@sts.wtf</email>
<web>www.sts.wtf</web>
<picture>` + squadxmlLogo + `</picture>
<title>Special Tactical Service</title>`
	squadxmlEnd = `</squad>`
)

type Member struct {
	UserOption32 string         `db:"userOption32"`
	Username     string         `db:"username"`
	RankTitle    string         `db:"rankTitle"`
	UserOption33 sql.NullString `db:"userOption33"`
}

var (
	db *sqlx.DB
)

func connectToDB() {
	user := os.Getenv("SQUADXML_DB_USER")
	password := os.Getenv("SQUADXML_DB_PASSWORD")
	host := os.Getenv("SQUADXML_DB_HOST")
	database := os.Getenv("SQUADXML_DB")
	var err error
	db, err = sqlx.Connect("mysql", mysqlConnection(user, password, host, database))

	if err != nil {
		logrus.WithError(err).Error("Error connecting to database")
		return
	}

	if err := db.Ping(); err != nil {
		logrus.WithError(err).Error("Error pinging database")
	}
}

func mysqlConnection(user, password, host, database string) string {
	return user + ":" + password + "@" + host + "/" + database
}

func disconnectDB() {
	if err := db.Close(); err != nil {
		logrus.WithError(err).Error("Error closing database connection")
	}
}

func buildSquadXML() {
	for {
		logrus.Info("Rebuilding squad.xml...")
		connectToDB()
		member := getMember()

		if member != nil {
			writeSquadXMLToFile(member)
		}

		disconnectDB()
		time.Sleep(time.Minute * 10)
	}
}

func getMember() []Member {
	query := `SELECT wcf1_user_option_value.userOption32,
		wcf1_user.username,
		wcf1_user_rank.rankTitle,
		wcf1_user_option_value.userOption33
        FROM wcf1_user, wcf1_user_option_value, wcf1_user_to_group, wcf1_user_rank
        WHERE wcf1_user.userID = wcf1_user_to_group.userID AND wcf1_user.userID = wcf1_user_option_value.userID AND wcf1_user_to_group.groupID = 6 AND wcf1_user.rankID = wcf1_user_rank.rankID AND wcf1_user_option_value.userOption32 != 0
        ORDER BY wcf1_user.userID ASC`
	var member []Member

	if err := db.Select(&member, query); err != nil {
		logrus.WithError(err).Error("Error selecting members")
		return nil
	}

	logrus.WithField("count", len(member)).Info("Selected members")
	return member
}

func writeSquadXMLToFile(member []Member) {
	xml := squadxmlHead + squadxmlSquad

	for _, m := range member {
		motto := ""

		if m.UserOption33.Valid {
			motto = " - " + m.UserOption33.String
		}

		xml += `<member id="` + m.UserOption32 + `" nick="` + m.Username + `">
<name>` + m.Username + `</name>
<email></email>
<icq>N/A</icq>
<remark>` + m.RankTitle + motto + `</remark>
</member>`
	}

	xml += squadxmlEnd

	if err := ioutil.WriteFile(filepath.Join(os.Getenv("SQUADXML_PATH"), squadxmlFile), []byte(xml), 0777); err != nil {
		logrus.WithError(err).Error("Error writing squadxml to file")
		return
	}

	logrus.WithField("xml", xml).Info("Successfully build squad.xml")
}

func main() {
	go buildSquadXML()
	path := os.Getenv("SQUADXML_PATH")
	logrus.WithField("squadxml_path", path).Info("Starting server...")

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(path, squadxmlFile))
	}))

	if err := http.ListenAndServe(os.Getenv("SQUADXML_HOST"), nil); err != nil {
		logrus.WithError(err).Fatal("Error starting server")
	}
}
