package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"xm-exercise/internal/utils"
	// Replace with actual imports for your application's models
	"xm-exercise/pkg/models" // Assuming models package exists and contains Company structs
)

// makeHTTPRequest is a helper to simplify making HTTP requests.
// It returns the *http.Response and an error. Callers must close the response body.
func makeHTTPRequest(client *http.Client,
	method, url string,
	body interface{},
	token string,
) (*http.Response, error) {
	var reqBody bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&reqBody).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request %s %s: %w", method, url, err)
	}

	return resp, nil
}

// checkHealth check if the server is healthy and ready to accept connections
func checkHealth(client *http.Client, baseURL string) error {
	resp, err := makeHTTPRequest(client, "GET", baseURL+"/health", nil, "")
	if err != nil {
		return err
	}

	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server not ready yet")
	}
	return nil
}

// registerUser sends a POST request to register a new user.
func registerUser(client *http.Client, baseURL string, name, email, password string) error {
	registerReq := models.UserRegistration{Name: name, Email: email, Password: password}
	resp, err := makeHTTPRequest(client, "POST", baseURL+"/auth/register", registerReq, "")
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			log.Println("invalid or no response data")
		}
		return fmt.Errorf("failed to register user %s: received status %d, body: %v",
			name, resp.StatusCode, errorBody)
	}
	log.Printf("Successfully registered user: %s", name)
	return nil
}

// loginUser sends a POST request to login and returns the auth token.
func loginUser(client *http.Client, baseURL string, email, password string) (string, error) {
	loginReq := models.UserLogin{Email: email, Password: password}
	resp, err := makeHTTPRequest(client, "POST", baseURL+"/auth/login", loginReq, "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			log.Println("invalid or no response data")
		}
		return "", fmt.Errorf("failed to login user %s: received status %d, body: %v",
			email, resp.StatusCode, errorBody)
	}

	var loginRes models.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginRes); err != nil {
		return "", fmt.Errorf("failed to decode login response body: %w", err)
	}
	if loginRes.Token == "" {
		return "", fmt.Errorf("login response token is empty for user %s", email)
	}
	log.Printf("Successfully logged in user: %s, received token", email)
	return loginRes.Token, nil
}

// createCompany sends a POST request to create a company and returns the created company details.
func createCompany(
	client *http.Client,
	baseURL string, token string,
	companyReq models.CompanyCreateRequest,
) (models.CompanyResponse, error) {
	resp, err := makeHTTPRequest(client, "POST", baseURL+"/companies", companyReq, token)
	if err != nil {
		return models.CompanyResponse{}, err
	}
	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != http.StatusCreated {
		var errorBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			log.Println("invalid or no response data")
		}

		return models.CompanyResponse{},
			fmt.Errorf("failed to create company %s: received status %d, body: %v",
				companyReq.Name, resp.StatusCode, errorBody)
	}

	var companyRes models.CompanyResponse
	if err := json.NewDecoder(resp.Body).Decode(&companyRes); err != nil {
		return models.CompanyResponse{},
			fmt.Errorf("failed to decode create company response body: %w", err)
	}
	if companyRes.ID == "" {
		return models.CompanyResponse{},
			fmt.Errorf("created company ID is empty for company %s", companyReq.Name)
	}

	log.Printf("Successfully created company: %s (ID: %s)", companyRes.Name, companyRes.ID)
	return companyRes, nil
}

// getCompany sends a GET request to retrieve a company by ID.
func getCompany(
	client *http.Client,
	baseURL, token, companyID string,
) (models.CompanyResponse, error) {
	url := fmt.Sprintf("%s/companies/%s", baseURL, companyID)
	resp, err := makeHTTPRequest(client, "GET", url, nil, token)
	if err != nil {
		// Don't return error immediately, caller needs to check status code
		return models.CompanyResponse{}, err
	}
	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			log.Println("invalid or no response data")
		}
		return models.CompanyResponse{}, fmt.Errorf("failed to get company with ID %s: received status %d, body: %v",
			companyID, resp.StatusCode, errorBody)
	}

	var companyRes models.CompanyResponse
	if err := json.NewDecoder(resp.Body).Decode(&companyRes); err != nil {
		return models.CompanyResponse{}, fmt.Errorf("failed to decode get company response body for ID %s: %w",
			companyID, err)
	}
	if companyRes.ID != companyID {
		return models.CompanyResponse{}, fmt.Errorf("retrieved company ID mismatch for ID %s: expected %s, got %s",
			companyID, companyID, companyRes.ID)
	}
	log.Printf("Successfully retrieved company with ID: %s", companyID)
	return companyRes, nil
}

// updateCompany sends a PATCH request to update a company.
func updateCompany(
	client *http.Client,
	baseURL, token, companyID string,
	updates models.CompanyUpdateRequest,
	status int,
) error {
	url := fmt.Sprintf("%s/companies/%s", baseURL, companyID)
	resp, err := makeHTTPRequest(client, "PATCH", url, updates, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != status {
		var errorBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return fmt.Errorf("invalid or no response data")
		}
		return fmt.Errorf("failed to update company with ID %s: received status %d, body: %v",
			companyID, resp.StatusCode, errorBody)
	}
	if status >= 200 && status <= 300 {
		log.Printf("Successfully updated company with ID: %s", companyID)
	} else {
		log.Printf("Successfully checked error condition on company with ID: %s", companyID)
	}

	return nil
}

// deleteCompany sends a DELETE request to delete a company.
func deleteCompany(client *http.Client, baseURL string, token string, companyID string) error {
	url := fmt.Sprintf("%s/companies/%s", baseURL, companyID)
	resp, err := makeHTTPRequest(client, "DELETE", url, nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != http.StatusOK {
		var errorBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return fmt.Errorf("invalid or no response data")
		}
		return fmt.Errorf("failed to delete company with ID %s: received status %d, body: %v",
			companyID, resp.StatusCode, errorBody)
	}
	log.Printf("Successfully deleted company with ID: %s", companyID)
	return nil
}

// confirmCompanyDeleted sends a GET request and expects a 404 Not Found status.
func confirmCompanyDeleted(client *http.Client, baseURL string, token string, companyID string) error {
	url := fmt.Sprintf("%s/companies/%s", baseURL, companyID)
	resp, err := makeHTTPRequest(client, "GET", url, nil, token)
	if err != nil {
		return fmt.Errorf("unexpected error fetching company ID %s after deletion: %w", companyID, err)
	}
	defer resp.Body.Close() //nolint:errcheck // Closing errors are typically unrecoverable.

	if resp.StatusCode != http.StatusNotFound {
		var errorBody map[string]string
		err = json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			return fmt.Errorf("invalid or no response data")
		}
		return fmt.Errorf("company with ID %s should not be found after deletion: received status %d, body: %v",
			companyID, resp.StatusCode, errorBody)
	}
	log.Printf("Successfully confirmed company with ID: %s is deleted (received 404)", companyID)
	return nil
}

func RunE2ETest() {
	log.Println("Starting E2E test runner...")

	client := &http.Client{Timeout: 30 * time.Second}

	baseURL := os.Getenv("E2E_BASE_URL")
	if baseURL == "" {
		baseURL = "localhost:8080"
	}
	if !strings.HasPrefix(baseURL, "http") {
		baseURL = fmt.Sprintf("http://%s", baseURL)
	}

	log.Printf("Running E2E test against base URL: %s", baseURL)

	// --- Test Sequence ---
	// 0. Check health
	log.Println("Step 1: Waiting for server to be ready...")
	for retry := 0; retry < 10; retry++ {
		if err := checkHealth(client, baseURL); err != nil {
			log.Println(err.Error())
			retry++
			time.Sleep(time.Second)
			continue
		}
		log.Println("Server is ready")
		break
	}

	baseURL = fmt.Sprintf("%s/api/v1", strings.TrimSuffix(baseURL, "/"))

	name := fmt.Sprintf("John Doe _%d", time.Now().UnixNano())
	email := fmt.Sprintf("jd_%d@email.com", time.Now().UnixNano())
	password := "e2e_password123"

	// 1. Register User
	log.Println("Step 1: Registering user...")
	if err := registerUser(client, baseURL, name, email, password); err != nil {
		log.Fatalf("E2E Test Failed at Register User: %v", err)
	}

	// 2. Login User and collect token
	log.Println("Step 2: Logging in user...")
	token, err := loginUser(client, baseURL, email, password)
	if err != nil {
		log.Fatalf("E2E Test Failed at Login User: %v", err)
	}
	if token == "" {
		log.Fatal("E2E Test Failed: Authentication token not obtained after login")
	}

	// 3. Create Company
	log.Println("Step 3: Creating company...")
	createReq := models.CompanyCreateRequest{
		Name:          fmt.Sprintf("Company %s", utils.GenerateRandomString(5)),
		Description:   aws.String("An end-to-end test company description."),
		EmployeeCount: 42,
		Registered:    aws.Bool(true),
		Type:          models.TypeNonProfit,
	}
	createdCompany, err := createCompany(client, baseURL, token, createReq)
	if err != nil {
		log.Fatalf("E2E Test Failed at Create Company: %v", err)
	}

	// 4. Get Company
	log.Println("Step 4: Getting company...")
	retrievedCompany, err := getCompany(client, baseURL, token, createdCompany.ID)
	if err != nil {
		log.Fatalf("E2E Test Failed at Get Company (first time): %v", err)
	}
	// Basic assertions
	if retrievedCompany.Name != createReq.Name {
		log.Fatalf("E2E Test Failed: Retrieved company name mismatch. "+
			"Expected: %s, Got: %s", createReq.Name, retrievedCompany.Name)
	}
	if *retrievedCompany.Description != *createReq.Description {
		log.Fatalf("E2E Test Failed: Retrieved company description mismatch. "+
			"Expected: %s, Got: %s", *createReq.Description, *retrievedCompany.Description)
	}
	if retrievedCompany.EmployeeCount != createReq.EmployeeCount {
		log.Fatalf("E2E Test Failed: Retrieved company employee count mismatch. Expected: %d, Got: %d",
			createReq.EmployeeCount, retrievedCompany.EmployeeCount)
	}
	if *retrievedCompany.Registered != *createReq.Registered {
		log.Fatalf("E2E Test Failed: Retrieved company registered status mismatch. Expected: %v, Got: %v",
			*createReq.Registered, *retrievedCompany.Registered)
	}
	if string(retrievedCompany.Type) != string(createReq.Type) {
		log.Fatalf("E2E Test Failed: Retrieved company type mismatch."+
			" Expected: %s, Got: %s", createReq.Type, retrievedCompany.Type)
	}
	if retrievedCompany.CreatedAt.IsZero() {
		log.Fatalf("E2E Test Failed: Retrieved company CreatedAt is zero")
	}
	if retrievedCompany.UpdatedAt.IsZero() {
		log.Fatalf("E2E Test Failed: Retrieved company UpdatedAt is zero")
	}

	// 5. Update Company
	log.Println("Step 5: Updating company...")
	newName := fmt.Sprintf("Updated %s", utils.GenerateRandomString(5))
	newEmployeeCount := 100
	updateReq := models.CompanyUpdateRequest{
		Name:          &newName,
		EmployeeCount: &newEmployeeCount,
	}

	if err := updateCompany(
		client,
		baseURL,
		token,
		createdCompany.ID,
		updateReq,
		http.StatusOK,
	); err != nil {
		log.Fatalf("E2E Test Failed at Update Company: %v", err)
	}

	// 6. Get Company again and confirm updates
	log.Println("Step 6: Getting company again to confirm updates...")
	updatedCompany, err := getCompany(client, baseURL, token, createdCompany.ID)
	if err != nil {
		log.Fatalf("E2E Test Failed at Get Company (after update): %v", err)
	}
	// Assert updated fields
	if updatedCompany.Name != newName {
		log.Fatalf("E2E Test Failed: Updated company name mismatch."+
			" Expected: %s, Got: %s", newName, updatedCompany.Name)
	}
	if updatedCompany.EmployeeCount != newEmployeeCount {
		log.Fatalf("E2E Test Failed: Updated company employee count mismatch."+
			" Expected: %d, Got: %d", newEmployeeCount, updatedCompany.EmployeeCount)
	}
	// Assert other fields are unchanged if not in updateReq
	if *updatedCompany.Description != *retrievedCompany.Description {
		log.Fatalf("E2E Test Failed: Updated company description changed unexpectedly."+
			" Expected: %s, Got: %s", *retrievedCompany.Description, *updatedCompany.Description)
	}
	if *updatedCompany.Registered != *retrievedCompany.Registered {
		log.Fatalf("E2E Test Failed: Updated company registered status changed unexpectedly."+
			" Expected: %v, Got: %v", *retrievedCompany.Registered, *updatedCompany.Registered)
	}
	if updatedCompany.Type != retrievedCompany.Type {
		log.Fatalf("E2E Test Failed: Updated company type changed unexpectedly."+
			" Expected: %s, Got: %s", retrievedCompany.Type, updatedCompany.Type)
	}
	// Check UpdatedAt timestamp
	if !updatedCompany.UpdatedAt.After(retrievedCompany.UpdatedAt) {
		log.Fatalf("E2E Test Failed: UpdatedAt timestamp did not update after PATCH")
	}

	// 7. Update Company
	log.Println("Step 7: Updating company with name longer than 15 characters...")
	newName = fmt.Sprintf("Updated %s", utils.GenerateRandomString(15))
	updateReq = models.CompanyUpdateRequest{
		Name: &newName,
	}
	if err := updateCompany(
		client,
		baseURL,
		token,
		createdCompany.ID,
		updateReq,
		http.StatusBadRequest,
	); err != nil {
		log.Fatalf("E2E Test Failed at Update Company: %v", err)
	}

	// 7. Delete Company
	log.Println("Step 8: Deleting company...")
	if err := deleteCompany(client, baseURL, token, createdCompany.ID); err != nil {
		log.Fatalf("E2E Test Failed at Delete Company: %v", err)
	}

	// 8. Confirm it is deleted again
	log.Println("Step 9: Confirming company is deleted...")
	if err := confirmCompanyDeleted(client, baseURL, token, createdCompany.ID); err != nil {
		// Note: confirmCompanyDeleted specifically checks for a 404, so a non-404 is an error
		log.Fatalf("E2E Test Failed at Confirm Company Deleted: %v", err)
	}

	log.Println("E2E Test Sequence Completed SUCCESSFULLY!")
}
