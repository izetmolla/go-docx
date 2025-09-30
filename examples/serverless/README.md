# Serverless DOCX Processing Example

This example demonstrates how to use the go-docx library for serverless processing with MinIO, where no file system is involved.

## Features

- **MinIO-to-MinIO Processing**: Download templates from MinIO, process them in memory, and upload results back
- **No File System**: All processing happens in memory using byte arrays
- **Custom Functions**: Support for custom data processing functions
- **Serverless Ready**: Perfect for AWS Lambda, Google Cloud Functions, Azure Functions, etc.

## Prerequisites

1. **MinIO Server**: You need a running MinIO server. You can start one locally using Docker:

```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"
```

2. **Template File**: Upload a DOCX template file to MinIO at `templates/report-template.docx`

## Usage

### Basic Serverless Processing

```go
// Download template bytes from MinIO
templateBytes, err := minioClient.GetObject(bucketName, "templates/report-template.docx", minio.GetObjectOptions{})
if err != nil {
    log.Fatal(err)
}
defer templateBytes.Close()

templateData, err := ioutil.ReadAll(templateBytes)
if err != nil {
    log.Fatal(err)
}

// Process template bytes with data (no file system involved)
processedBytes, err := docx.CompleteTemplateFromBytesToBytes(templateData, data)
if err != nil {
    log.Fatal(err)
}

// Upload processed bytes back to MinIO
minioClient.PutObject(
    bucketName, 
    "processed/report.docx", 
    bytes.NewReader(processedBytes), 
    int64(len(processedBytes)), 
    minio.PutObjectOptions{},
)
```

### With Custom Functions

```go
// With custom functions
processedBytes, err := docx.CompleteTemplateFromBytesToBytesWithFuncs(templateData, data, funcMap)
```

## API Functions

### `CompleteTemplateFromBytesToBytes`

Processes a DOCX template from bytes and returns processed bytes.

```go
func CompleteTemplateFromBytesToBytes(templateBytes []byte, data PlaceholderMap) ([]byte, error)
```

**Parameters:**
- `templateBytes`: The DOCX template as a byte array
- `data`: A map of placeholder keys to replacement values

**Returns:**
- `[]byte`: The processed DOCX document as bytes
- `error`: Any error that occurred during processing

### `CompleteTemplateFromBytesToBytesWithFuncs`

Processes a DOCX template from bytes with custom functions and returns processed bytes.

```go
func CompleteTemplateFromBytesToBytesWithFuncs(templateBytes []byte, data PlaceholderMap, funcMap map[string]interface{}) ([]byte, error)
```

**Parameters:**
- `templateBytes`: The DOCX template as a byte array
- `data`: A map of placeholder keys to replacement values
- `funcMap`: A map of function names to function implementations

**Returns:**
- `[]byte`: The processed DOCX document as bytes
- `error`: Any error that occurred during processing

## Running the Example

1. **Start MinIO** (if not already running):
```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e "MINIO_ROOT_USER=minioadmin" \
  -e "MINIO_ROOT_PASSWORD=minioadmin" \
  minio/minio server /data --console-address ":9001"
```

2. **Upload a template** to MinIO:
   - Open MinIO console at http://localhost:9001
   - Create a bucket named `docx-templates`
   - Upload a DOCX template file to `templates/report-template.docx`

3. **Run the example**:
```bash
cd examples/serverless
go mod init serverless-example
go mod tidy
go run main.go
```

## Template Placeholders

The example uses the following placeholders in the DOCX template:

- `{company_name}` - Company name
- `{report_date}` - Report date
- `{total_revenue}` - Total revenue
- `{employee_count}` - Number of employees
- `{quarter}` - Quarter information
- `{ceo_name}` - CEO name
- `{department}` - Department name
- `{project_name}` - Project name
- `{status}` - Project status
- `{notes}` - Additional notes

## Serverless Deployment

This example is perfect for serverless deployment:

### AWS Lambda
- Package the binary and upload to Lambda
- Configure MinIO credentials via environment variables
- Set appropriate memory and timeout values

### Google Cloud Functions
- Deploy using `gcloud functions deploy`
- Configure MinIO credentials via environment variables
- Set appropriate memory and timeout values

### Azure Functions
- Deploy using Azure Functions Core Tools
- Configure MinIO credentials via application settings
- Set appropriate memory and timeout values

## Benefits

1. **No File System**: Perfect for serverless environments
2. **Memory Efficient**: Processes documents entirely in memory
3. **Scalable**: Can handle high concurrency in serverless environments
4. **Cost Effective**: Pay only for processing time
5. **Fast**: No disk I/O overhead
