package handlers

import (
	"database/sql"
	"goscraper/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

const (
	clearJobs = `DELETE FROM customJobs;`

	createJobsTable = `
		CREATE TABLE IF NOT EXISTS customJobs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			location TEXT,
			created_at INTEGER,
			company TEXT,
			apply_url TEXT,
			image_url TEXT,
			description TEXT,
			skills TEXT,
			expired BOOLEAN,
			salary TEXT
		);
	`

	insertJobs = `
		INSERT INTO customJobs (title, location, created_at, company, apply_url, image_url, description, skills, expired, salary)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	updateJob = `
		UPDATE customJobs
		SET title = ?, location = ?, created_at = ?, company = ?, apply_url = ?, image_url = ?, description = ?, skills = ?, expired = ?, salary = ?
		WHERE id = ?;
	`

	deleteJob = `DELETE FROM customJobs WHERE id = ?;`

	getallJobs = `
    SELECT id, title, description, skills, salary, location, created_at, company, apply_url, image_url
    FROM customJobs;
`
)

var newdb *sql.DB

func init() {
	var err error
	newdb, err = sql.Open("sqlite", "jobs.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Try creating the table
	_, err = newdb.Exec(createJobsTable)
	if err != nil {
		log.Fatalf("Error creating jobs table: %v", err)
	}
}

func closeDB() {
	if err := newdb.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}

func InsertJob(c *gin.Context, job models.UserDefinedJob) {
	skillsString := strings.Join(job.Skills, ",")

	_, err := newdb.Exec(insertJobs, job.Title, job.Location, job.CreatedAt, job.Company, job.ApplyURL, job.ImageUrl, job.Description, skillsString, job.Expired, job.Salary)
	if err != nil {
		log.Printf("Error inserting job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job added successfully!"})
}

func UpdateJob(c *gin.Context, job models.UserDefinedJob) {
	skillsString := strings.Join(job.Skills, ",")

	_, err := newdb.Exec(updateJob, job.Title, job.Location, job.CreatedAt, job.Company, job.ApplyURL, job.ImageUrl, job.Description, skillsString, job.Expired, job.Salary, job.ID)
	if err != nil {
		log.Printf("Error updating job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job updated successfully!"})
}

func DeleteJob(c *gin.Context, jobID string) {
	_, err := newdb.Exec(deleteJob, jobID)
	if err != nil {
		log.Printf("Error deleting job: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Job deleted successfully!"})
}

func GetAllJobs(c *gin.Context) {
	rows, err := newdb.Query(getallJobs)
	if err != nil {
		log.Printf("Error getting all jobs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var jobs []models.UserDefinedJob

	for rows.Next() {
		var job models.UserDefinedJob
		var skillsString string

		if err := rows.Scan(&job.ID, &job.Title, &job.Description, &skillsString, &job.Salary, &job.Location, &job.CreatedAt, &job.Company, &job.ApplyURL, &job.ImageUrl); err != nil {
			log.Printf("Error scanning row: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}
