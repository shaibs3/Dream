package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"dream/internal/models"
	"dream/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSystem(t *testing.T) {
	if os.Getenv("TEST_DATABASE_URL") == "" {
		os.Setenv("TEST_DATABASE_URL", "postgres://dreamuser:dreampass@localhost:5432/dreamdb?sslmode=disable")
	}
	// Initialize test database
	db, err := gorm.Open(postgres.Open(os.Getenv("TEST_DATABASE_URL")), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to test database")

	// Clean up database before test
	require.NoError(t, db.Exec("TRUNCATE users, linux_processes, windows_processes CASCADE").Error, "Failed to clean up database")

	// Read test data files
	linuxData, err := os.ReadFile("testdata/linux_ps.txt")
	require.NoError(t, err, "Failed to read Linux test data")

	windowsData, err := os.ReadFile("testdata/windows_tasklist.txt")
	require.NoError(t, err, "Failed to read Windows test data")

	// Create test requests for two different users
	linuxUser := types.MessageRequest{
		MachineID:     "linux-pc-001",
		MachineName:   "Linux Workstation",
		OSVersion:     "Ubuntu",
		Timestamp:     time.Now(),
		CommandType:   "ps",
		UserName:      "john.doe",
		UserID:        "john.doe@law.edu",
		Faculty:       "Law",
		CommandOutput: string(linuxData),
	}

	windowsUser := types.MessageRequest{
		MachineID:     "win-pc-001",
		MachineName:   "Windows Workstation",
		OSVersion:     "Windows 10",
		Timestamp:     time.Now(),
		CommandType:   "tasklist",
		UserName:      "jane.smith",
		UserID:        "jane.smith@med.edu",
		Faculty:       "Medicine",
		CommandOutput: string(windowsData),
	}

	// Convert requests to JSON
	linuxJSON, err := json.Marshal(linuxUser)
	require.NoError(t, err, "Failed to marshal Linux user request")

	windowsJSON, err := json.Marshal(windowsUser)
	require.NoError(t, err, "Failed to marshal Windows user request")

	// Send requests to the application
	baseURL := "http://localhost:8080"

	// Send Linux user data
	resp, err := http.Post(baseURL+"/upload", "application/json", bytes.NewBuffer(linuxJSON))
	require.NoError(t, err, "Failed to send Linux user data")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, http.StatusOK, resp.StatusCode, "Linux user request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Send Windows user data
	resp, err = http.Post(baseURL+"/upload", "application/json", bytes.NewBuffer(windowsJSON))
	require.NoError(t, err, "Failed to send Windows user data")
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		require.Equal(t, http.StatusOK, resp.StatusCode, "Windows user request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Wait a bit for Kafka consumer to process messages
	time.Sleep(2 * time.Second)

	// Verify database state
	t.Run("Verify Users", func(t *testing.T) {
		var users []models.User
		require.NoError(t, db.Find(&users).Error, "Failed to query users")

		assert.Equal(t, 2, len(users), "Expected 2 users, got %d", len(users))

		// Verify Linux user
		var linuxUser models.User
		err := db.Where("user_id = ?", "john.doe@law.edu").First(&linuxUser).Error
		assert.NoError(t, err, "Linux user not found")
		if err == nil {
			assert.Equal(t, "Law", linuxUser.Faculty, "Expected faculty 'Law', got '%s'", linuxUser.Faculty)
			assert.Equal(t, "Ubuntu", linuxUser.OSVersion, "Expected OS 'Ubuntu', got '%s'", linuxUser.OSVersion)
		}

		// Verify Windows user
		var windowsUser models.User
		err = db.Where("user_id = ?", "jane.smith@med.edu").First(&windowsUser).Error
		assert.NoError(t, err, "Windows user not found")
		if err == nil {
			assert.Equal(t, "Medicine", windowsUser.Faculty, "Expected faculty 'Medicine', got '%s'", windowsUser.Faculty)
			assert.Equal(t, "Windows 10", windowsUser.OSVersion, "Expected OS 'Windows 10', got '%s'", windowsUser.OSVersion)
		}
	})

	t.Run("Verify Linux Processes", func(t *testing.T) {
		var processes []models.LinuxProcess
		require.NoError(t, db.Where("user_id = ?", "john.doe@law.edu").Find(&processes).Error, "Failed to query Linux processes")

		assert.Equal(t, 2, len(processes), "Expected 2 Linux processes, got %d", len(processes))

		// Verify specific process details
		for _, proc := range processes {
			if proc.Command == "/usr/sbin/sshd" {
				assert.Equal(t, "root", proc.User, "Expected user 'root' for sshd, got '%s'", proc.User)
				assert.Equal(t, 0.5, proc.CPU, "Expected CPU 0.5 for sshd, got %f", proc.CPU)
			} else if proc.Command == "vim test.txt" {
				assert.Equal(t, "user", proc.User, "Expected user 'user' for vim, got '%s'", proc.User)
				assert.Equal(t, 2.0, proc.CPU, "Expected CPU 2.0 for vim, got %f", proc.CPU)
			}
		}
	})

	t.Run("Verify Windows Processes", func(t *testing.T) {
		var processes []models.WindowsProcess
		require.NoError(t, db.Where("user_id = ?", "jane.smith@med.edu").Find(&processes).Error, "Failed to query Windows processes")

		assert.Equal(t, 2, len(processes), "Expected 2 Windows processes, got %d", len(processes))

		// Verify specific process details
		for _, proc := range processes {
			if proc.ImageName == "chrome.exe" {
				assert.Equal(t, "Console", proc.SessionName, "Expected session 'Console' for chrome, got '%s'", proc.SessionName)
				assert.Equal(t, "50,000 K", proc.MemUsage, "Expected memory '50,000 K' for chrome, got '%s'", proc.MemUsage)
			} else if proc.ImageName == "notepad.exe" {
				assert.Equal(t, "Console", proc.SessionName, "Expected session 'Console' for notepad, got '%s'", proc.SessionName)
				assert.Equal(t, "5,000 K", proc.MemUsage, "Expected memory '5,000 K' for notepad, got '%s'", proc.MemUsage)
			}
		}
	})
}
