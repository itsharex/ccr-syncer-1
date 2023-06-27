package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(dbPath string) (DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// create table info && progress, if not exists
	// all is tuple (string, string)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS jobs (job_name TEXT PRIMARY KEY, job_info TEXT)")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS progresses (job_name TEXT PRIMARY KEY, progress TEXT)")
	if err != nil {
		return nil, err
	}

	return &SQLiteDB{db: db}, nil
}

func (s *SQLiteDB) AddJob(jobName string, jobInfo string) error {
	// check job name exists, if exists, return error
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE job_name = ?", jobName).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrJobExists
	}

	// insert job info
	_, err = s.db.Exec("INSERT INTO jobs (job_name, job_info) VALUES (?, ?)", jobName, jobInfo)
	return err
}

// Update Job
func (s *SQLiteDB) UpdateJob(jobName string, jobInfo string) error {
	// check job name exists, if not exists, return error
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE job_name = ?", jobName).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrJobNotExists
	}

	// update job info
	_, err = s.db.Exec("UPDATE jobs SET job_info = ? WHERE job_name = ?", jobInfo, jobName)
	return err
}

func (s *SQLiteDB) IsJobExist(jobName string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE job_name = ?", jobName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SQLiteDB) GetAllJobs() (map[string]string, error) {
	rows, err := s.db.Query("SELECT job_name, job_info FROM jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var jobName, jobInfo string
		err = rows.Scan(&jobName, &jobInfo)
		if err != nil {
			return nil, err
		}
		result[jobName] = jobInfo
	}
	return result, nil
}

func (s *SQLiteDB) UpdateProgress(jobName string, progress string) error {
	// check job name exists, if not exists, return error
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE job_name = ?", jobName).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrJobNotExists
	}

	// update progress
	_, err = s.db.Exec("UPDATE progresses SET progress = ? WHERE job_name = ?", progress, jobName)
	return err
}

func (s *SQLiteDB) IsProgressExist(jobName string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM progresses WHERE job_name = ?", jobName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *SQLiteDB) GetProgress(jobName string) (string, error) {
	var progress string
	err := s.db.QueryRow("SELECT progress FROM progresses WHERE job_name = ?", jobName).Scan(&progress)
	if err != nil {
		return "", err
	}
	return progress, nil
}
