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

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSystem(t *testing.T) {
	// Initialize test database
	db, err := gorm.Open(postgres.Open(os.Getenv("TEST_DATABASE_URL")), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up database before test
	if err := db.Exec("TRUNCATE users, linux_process, windows_process CASCADE").Error; err != nil {
		t.Fatalf("Failed to clean up database: %v", err)
	}

	// Read test data files
	linuxData, err := os.ReadFile("testdata/linux_ps.txt")
	if err != nil {
		t.Fatalf("Failed to read Linux test data: %v", err)
	}

	windowsData, err := os.ReadFile("testdata/windows_tasklist.txt")
	if err != nil {
		t.Fatalf("Failed to read Windows test data: %v", err)
	}

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
	if err != nil {
		t.Fatalf("Failed to marshal Linux user request: %v", err)
	}

	windowsJSON, err := json.Marshal(windowsUser)
	if err != nil {
		t.Fatalf("Failed to marshal Windows user request: %v", err)
	}

	// Send requests to the application
	baseURL := "http://localhost:8080"

	// Send Linux user data
	resp, err := http.Post(baseURL+"/upload", "application/json", bytes.NewBuffer(linuxJSON))
	if err != nil {
		t.Fatalf("Failed to send Linux user data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Linux user request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Send Windows user data
	resp, err = http.Post(baseURL+"/upload", "application/json", bytes.NewBuffer(windowsJSON))
	if err != nil {
		t.Fatalf("Failed to send Windows user data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Windows user request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Wait a bit for Kafka consumer to process messages
	time.Sleep(2 * time.Second)

	// Verify database state
	t.Run("Verify Users", func(t *testing.T) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			t.Fatalf("Failed to query users: %v", err)
		}

		if len(users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(users))
		}

		// Verify Linux user
		var linuxUser models.User
		if err := db.Where("user_id = ?", "john.doe@law.edu").First(&linuxUser).Error; err != nil {
			t.Errorf("Linux user not found: %v", err)
		} else {
			if linuxUser.Faculty != "Law" {
				t.Errorf("Expected faculty 'Law', got '%s'", linuxUser.Faculty)
			}
			if linuxUser.OSVersion != "Ubuntu" {
				t.Errorf("Expected OS 'Ubuntu', got '%s'", linuxUser.OSVersion)
			}
		}

		// Verify Windows user
		var windowsUser models.User
		if err := db.Where("user_id = ?", "jane.smith@med.edu").First(&windowsUser).Error; err != nil {
			t.Errorf("Windows user not found: %v", err)
		} else {
			if windowsUser.Faculty != "Medicine" {
				t.Errorf("Expected faculty 'Medicine', got '%s'", windowsUser.Faculty)
			}
			if windowsUser.OSVersion != "Windows 10" {
				t.Errorf("Expected OS 'Windows 10', got '%s'", windowsUser.OSVersion)
			}
		}
	})

	t.Run("Verify Linux Processes", func(t *testing.T) {
		var processes []models.LinuxProcess
		if err := db.Where("user_id = ?", "john.doe@law.edu").Find(&processes).Error; err != nil {
			t.Fatalf("Failed to query Linux processes: %v", err)
		}

		if len(processes) != 2 {
			t.Errorf("Expected 2 Linux processes, got %d", len(processes))
		}

		// Verify specific process details
		for _, proc := range processes {
			if proc.Command == "/usr/sbin/sshd" {
				if proc.User != "root" {
					t.Errorf("Expected user 'root' for sshd, got '%s'", proc.User)
				}
				if proc.CPU != 0.5 {
					t.Errorf("Expected CPU 0.5 for sshd, got %f", proc.CPU)
				}
			} else if proc.Command == "vim test.txt" {
				if proc.User != "user" {
					t.Errorf("Expected user 'user' for vim, got '%s'", proc.User)
				}
				if proc.CPU != 2.0 {
					t.Errorf("Expected CPU 2.0 for vim, got %f", proc.CPU)
				}
			}
		}
	})

	t.Run("Verify Windows Processes", func(t *testing.T) {
		var processes []models.WindowsProcess
		if err := db.Where("user_id = ?", "jane.smith@med.edu").Find(&processes).Error; err != nil {
			t.Fatalf("Failed to query Windows processes: %v", err)
		}

		if len(processes) != 2 {
			t.Errorf("Expected 2 Windows processes, got %d", len(processes))
		}

		// Verify specific process details
		for _, proc := range processes {
			if proc.ImageName == "chrome.exe" {
				if proc.SessionName != "Console" {
					t.Errorf("Expected session 'Console' for chrome, got '%s'", proc.SessionName)
				}
				if proc.MemUsage != "50,000 K" {
					t.Errorf("Expected memory '50,000 K' for chrome, got '%s'", proc.MemUsage)
				}
			} else if proc.ImageName == "notepad.exe" {
				if proc.SessionName != "Console" {
					t.Errorf("Expected session 'Console' for notepad, got '%s'", proc.SessionName)
				}
				if proc.MemUsage != "5,000 K" {
					t.Errorf("Expected memory '5,000 K' for notepad, got '%s'", proc.MemUsage)
				}
			}
		}
	})
}
