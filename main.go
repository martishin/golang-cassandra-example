package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
)

func main() {
	cassandraHost := os.Getenv("CASSANDRA_CONTACT_POINT")
	cassandraUsername := os.Getenv("CASSANDRA_USERNAME")
	cassandraPassword := os.Getenv("CASSANDRA_PASSWORD")

	// connect to the cluster
	cluster := gocql.NewCluster(cassandraHost)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cassandraUsername, Password: cassandraPassword, AllowedAuthenticators: []string{"com.instaclustr.cassandra.auth.InstaclustrPasswordAuthenticator"}}

	session, err := cluster.CreateSession()
	if err != nil {
		log.Println("Error creating session:", err)
		return
	}
	defer session.Close()

	// create keyspace
	err = session.Query("CREATE KEYSPACE IF NOT EXISTS sleep_centre WITH REPLICATION = {'class' : 'SimpleStrategy', 'replication_factor' : 1};").Exec()
	if err != nil {
		log.Println("Error creating keyspace:", err)
		return
	}

	// create table
	err = session.Query("CREATE TABLE IF NOT EXISTS sleep_centre.sleep_study (name text, study_date date, sleep_time_hours float, PRIMARY KEY (name, study_date));").Exec()
	if err != nil {
		log.Println("Error creating table:", err)
		return
	}

	// insert some practice data
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-07', 8.2);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-08', 6.4);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('James', '2018-01-09', 7.5);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-07', 6.6);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-08', 6.3);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Bob', '2018-01-09', 6.7);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Emily', '2018-01-07', 7.2);").Exec()
	err = session.Query("INSERT INTO sleep_centre.sleep_study (name, study_date, sleep_time_hours) VALUES ('Emily', '2018-01-09', 7.5);").Exec()
	if err != nil {
		log.Println("Error inserting data:", err)
		return
	}

	// Return average sleep time for James
	var sleepTimeHours float32

	sleepTimeOutput := session.Query("SELECT avg(sleep_time_hours) FROM sleep_centre.sleep_study WHERE name = 'James';").Iter()
	sleepTimeOutput.Scan(&sleepTimeHours)
	fmt.Println("Average sleep time for James was: ", sleepTimeHours, "h")

	// return average sleep time for group
	sleepTimeOutput = session.Query("SELECT avg(sleep_time_hours) FROM sleep_centre.sleep_study;").Iter()
	sleepTimeOutput.Scan(&sleepTimeHours)
	fmt.Println("Average sleep time for the group was: ", sleepTimeHours, "h")
}
