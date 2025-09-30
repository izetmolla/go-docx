package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/izetmolla/go-docx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	// MinIO configuration
	endpoint        = "localhost:9000"
	accessKeyID     = "minioadmin"
	secretAccessKey = "minioadmin"
	useSSL          = false
	bucketName      = "docx-templates"
)

func main() {
	// Initialize MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	// Create bucket if it doesn't exist
	ctx := context.Background()
	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check if bucket already exists
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists != nil || !exists {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}

	// Example 1: Basic serverless processing
	fmt.Println("=== Example 1: Basic Serverless Processing ===")
	err = basicServerlessProcessing(minioClient, ctx)
	if err != nil {
		log.Printf("Basic serverless processing failed: %v", err)
	}

	// Example 2: Serverless processing with custom functions
	fmt.Println("\n=== Example 2: Serverless Processing with Custom Functions ===")
	err = serverlessProcessingWithFuncs(minioClient, ctx)
	if err != nil {
		log.Printf("Serverless processing with funcs failed: %v", err)
	}
}

// basicServerlessProcessing demonstrates basic MinIO-to-MinIO processing
func basicServerlessProcessing(minioClient *minio.Client, ctx context.Context) error {
	// Download template bytes from MinIO
	templateBytes, err := minioClient.GetObject(ctx, bucketName, "templates/report-template.docx", minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get template from MinIO: %w", err)
	}
	defer templateBytes.Close()

	templateData, err := io.ReadAll(templateBytes)
	if err != nil {
		return fmt.Errorf("failed to read template data: %w", err)
	}

	// Prepare data for template processing
	data := docx.PlaceholderMap{
		"company_name":    "Acme Corporation",
		"report_date":     time.Now().Format("2006-01-02"),
		"total_revenue":   "$1,234,567.89",
		"employee_count":  "150",
		"quarter":         "Q4 2024",
		"ceo_name":        "John Smith",
		"department":      "Engineering",
		"project_name":    "Serverless DOCX Processing",
		"status":          "Completed",
		"notes":           "Successfully implemented serverless processing with MinIO",
	}

	// Process template bytes with data (no file system involved)
	processedBytes, err := docx.CompleteTemplateFromBytesToBytes(templateData, data)
	if err != nil {
		return fmt.Errorf("failed to process template: %w", err)
	}

	// Upload processed bytes back to MinIO
	_, err = minioClient.PutObject(
		ctx,
		bucketName,
		"processed/report-"+time.Now().Format("20060102-150405")+".docx",
		bytes.NewReader(processedBytes),
		int64(len(processedBytes)),
		minio.PutObjectOptions{
			ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload processed document to MinIO: %w", err)
	}

	fmt.Println("✅ Successfully processed template and uploaded to MinIO")
	return nil
}

// serverlessProcessingWithFuncs demonstrates serverless processing with custom functions
func serverlessProcessingWithFuncs(minioClient *minio.Client, ctx context.Context) error {
	// Download template bytes from MinIO
	templateBytes, err := minioClient.GetObject(ctx, bucketName, "templates/report-template.docx", minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get template from MinIO: %w", err)
	}
	defer templateBytes.Close()

	templateData, err := io.ReadAll(templateBytes)
	if err != nil {
		return fmt.Errorf("failed to read template data: %w", err)
	}

	// Prepare data for template processing
	data := docx.PlaceholderMap{
		"company_name":    "TechCorp Solutions",
		"report_date":     time.Now().Format("2006-01-02"),
		"total_revenue":   "$2,345,678.90",
		"employee_count":  "250",
		"quarter":         "Q4 2024",
		"ceo_name":        "Jane Doe",
		"department":      "Research & Development",
		"project_name":    "Advanced Serverless Architecture",
		"status":          "In Progress",
		"notes":           "Implementing advanced serverless processing with custom functions",
	}

	// Define custom functions for data processing
	funcMap := map[string]interface{}{
		"uppercase": func(s string) string {
			return strings.ToUpper(s)
		},
		"lowercase": func(s string) string {
			return strings.ToLower(s)
		},
		"format_currency": func(amount string) string {
			// Simple currency formatting (in a real app, you'd use proper currency formatting)
			return fmt.Sprintf("USD %s", amount)
		},
		"current_timestamp": func() string {
			return time.Now().Format("2006-01-02 15:04:05")
		},
	}

	// Process template bytes with data and custom functions (no file system involved)
	processedBytes, err := docx.CompleteTemplateFromBytesToBytesWithFuncs(templateData, data, funcMap)
	if err != nil {
		return fmt.Errorf("failed to process template with functions: %w", err)
	}

	// Upload processed bytes back to MinIO
	_, err = minioClient.PutObject(
		ctx,
		bucketName,
		"processed/report-with-funcs-"+time.Now().Format("20060102-150405")+".docx",
		bytes.NewReader(processedBytes),
		int64(len(processedBytes)),
		minio.PutObjectOptions{
			ContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload processed document to MinIO: %w", err)
	}

	fmt.Println("✅ Successfully processed template with custom functions and uploaded to MinIO")
	return nil
}
