package main

import (
	nodes "SRSC-Client/type"
	file "SRSC-Client/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	// AWS SDK v1 imports
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr" // Import for AWS error handling
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	// Wails imports (unchanged)
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 导出上下文到全局变量以便获取上下文
var ContextX context.Context

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 导出上下文到全局变量以便获取上下文
	ContextX = a.ctx
}

// S3Manager struct using AWS SDK v1 S3 client
type S3Manager struct {
	client *s3.S3 // Changed from v2 s3.Client to v1 s3.S3
}

// NewS3Client connects to object storage using AWS SDK v1
func (a *S3Manager) NewS3Client(endpoint, region, accessKey, secretKey string) (*S3Manager, error) {
	ctx := context.Background() // Use background context for initial connection setup
	runtime.LogDebug(ContextX, "Attempting to connect using v1 SDK: "+endpoint)

	// Configure AWS session for v1
	awsCfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true), // Use path-style addressing, common for non-AWS S3
		// DisableSSL:      aws.Bool(true), // Uncomment if endpoint uses HTTP instead of HTTPS
	}

	// Create a new session
	sess, err := session.NewSession(awsCfg)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to create session: "+err.Error())
		return nil, fmt.Errorf("v1: unable to create AWS session: %v", err)
	}
	runtime.LogDebug(ContextX, "v1: Session created successfully for region: "+region)

	// Create S3 service client
	client := s3.New(sess)
	runtime.LogDebug(ContextX, "v1: S3 client created")

	// List buckets to verify connection (using WithContext variant)
	listInput := &s3.ListBucketsInput{}
	resp, err := client.ListBucketsWithContext(ctx, listInput)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to list buckets: "+err.Error())
		// Log AWS specific error details if available
		if aerr, ok := err.(awserr.Error); ok {
			runtime.LogError(ContextX, fmt.Sprintf("v1: AWS Error Code: %s, Message: %s", aerr.Code(), aerr.Message()))
		}
		return nil, fmt.Errorf("v1: failed to list buckets: %v", err)
	}

	runtime.LogDebug(ContextX, "v1: === S3 Buckets ===")
	for _, bucket := range resp.Buckets {
		if bucket.Name != nil && bucket.CreationDate != nil {
			logMsg := fmt.Sprintf("- %s (Created at: %v)\n", *bucket.Name, *bucket.CreationDate)
			runtime.LogDebug(ContextX, logMsg)
			fmt.Printf(logMsg) // Also print to console as before
		}
	}

	return &S3Manager{client: client}, nil
}

// GetAllS3NodesInfo remains the same as it doesn't interact with AWS SDK
func (a *S3Manager) GetAllS3NodesInfo() []nodes.Node {
	fileContent, err := os.ReadFile("./config.json")
	runtime.LogDebug(ContextX, string(fileContent))
	if err != nil {
		runtime.LogError(ContextX, err.Error())
		return nil
	}
	nodeList, err := nodes.GetNodes(fileContent)
	if err != nil {
		runtime.LogError(ContextX, err.Error())
		return nil
	}
	runtime.LogDebug(ContextX, "成功获取节点信息"+nodeList[0].NodeName)
	return nodeList
}

// GetNodeBucketInfo Get bucket information for a node using provided parameters (v1 SDK)
func (a *S3Manager) GetNodeBucketInfo(endpoint, accessKey, secretKey, region string) []nodes.NodeBucketInfo {
	runtime.LogDebug(ContextX, "v1: Getting node bucket info with provided parameters")

	var allNodesBucketInfo []nodes.NodeBucketInfo

	s3Manager, err := a.NewS3Client(endpoint, region, accessKey, secretKey)
	if err != nil {
		runtime.LogError(ContextX, "v1: Connection failed: "+err.Error())
		return nil // Return empty slice on connection failure
	}

	nodeBucketInfo := nodes.NodeBucketInfo{
		NodeName: "S3节点", // Default name as before
		EndPoint: endpoint,
		Buckets:  []nodes.BucketInfo{},
	}

	ctx := context.Background()
	listInput := &s3.ListBucketsInput{}
	resp, err := s3Manager.client.ListBucketsWithContext(ctx, listInput)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to list buckets: "+err.Error())
		return nil // Return empty slice if listing fails
	}

	for _, bucket := range resp.Buckets {
		if bucket.Name == nil || bucket.CreationDate == nil {
			runtime.LogDebug(ContextX, "v1: Skipping bucket with nil name or creation date")
			continue
		}

		bucketInfo := nodes.BucketInfo{
			Name:                *bucket.Name,
			CreationDate:        *bucket.CreationDate,
			UsedSpace:           0,
			TotalObjects:        0,
			VersioningEnabled:   false,
			PublicAccessBlocked: false,
			HasPolicy:           false,
			EncryptionEnabled:   false,
			EncryptionType:      "",
			HasLifecycleRules:   false,
			LifecycleRulesCount: 0,
			Region:              region, // Use the provided region
			WebsiteEnabled:      false,
		}

		// Get object count and used space using pagination
		listObjectsInput := &s3.ListObjectsV2Input{Bucket: bucket.Name}
		err := s3Manager.client.ListObjectsV2PagesWithContext(ctx, listObjectsInput,
			func(page *s3.ListObjectsV2Output, lastPage bool) bool {
				if page == nil {
					return false // Stop if page is nil
				}
				bucketInfo.TotalObjects += int64(len(page.Contents))
				for _, obj := range page.Contents {
					if obj.Size != nil {
						bucketInfo.UsedSpace += *obj.Size
					}
				}
				return !lastPage // Continue if not the last page
			})
		if err != nil {
			runtime.LogError(ContextX, fmt.Sprintf("v1: Failed getting objects for bucket %s: %s", *bucket.Name, err.Error()))
			// Continue to next bucket or property, depending on desired robustness
		}

		// Get Versioning status
		verInput := &s3.GetBucketVersioningInput{Bucket: bucket.Name}
		verOutput, err := s3Manager.client.GetBucketVersioningWithContext(ctx, verInput)
		if err == nil && verOutput.Status != nil && *verOutput.Status == s3.BucketVersioningStatusEnabled {
			bucketInfo.VersioningEnabled = true
			runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket %s has versioning enabled", *bucket.Name))
		} else if err != nil {
			// Log non-critical errors
			runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get versioning for %s: %s", *bucket.Name, err.Error()))
		}

		// Get Public Access Block status
		pubInput := &s3.GetPublicAccessBlockInput{Bucket: bucket.Name}
		pubOutput, err := s3Manager.client.GetPublicAccessBlockWithContext(ctx, pubInput)
		if err == nil && pubOutput.PublicAccessBlockConfiguration != nil {
			config := pubOutput.PublicAccessBlockConfiguration
			// Check if all block flags are true (using aws.BoolValue to safely dereference)
			if aws.BoolValue(config.BlockPublicAcls) && aws.BoolValue(config.BlockPublicPolicy) &&
				aws.BoolValue(config.IgnorePublicAcls) && aws.BoolValue(config.RestrictPublicBuckets) {
				bucketInfo.PublicAccessBlocked = true
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket %s has public access blocked", *bucket.Name))
			}
		} else if err != nil {
			// Check if the error is because no configuration exists (expected case)
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "NoSuchPublicAccessBlockConfiguration" {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get public access block for %s: %s", *bucket.Name, err.Error()))
			} else if !ok {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get public access block for %s: %s", *bucket.Name, err.Error()))
			}
		}

		// Get Bucket Policy status
		polInput := &s3.GetBucketPolicyInput{Bucket: bucket.Name}
		_, err = s3Manager.client.GetBucketPolicyWithContext(ctx, polInput)
		if err == nil {
			bucketInfo.HasPolicy = true
			runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket %s has a policy", *bucket.Name))
		} else {
			// Check if the error is NoSuchBucketPolicy (expected if no policy exists)
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "NoSuchBucketPolicy" {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get policy for %s: %s", *bucket.Name, err.Error()))
			} else if !ok {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get policy for %s: %s", *bucket.Name, err.Error()))
			}
		}

		// Get Encryption status
		encInput := &s3.GetBucketEncryptionInput{Bucket: bucket.Name}
		encOutput, err := s3Manager.client.GetBucketEncryptionWithContext(ctx, encInput)
		if err == nil && encOutput.ServerSideEncryptionConfiguration != nil && len(encOutput.ServerSideEncryptionConfiguration.Rules) > 0 {
			bucketInfo.EncryptionEnabled = true
			// Get encryption type from the first rule
			rule := encOutput.ServerSideEncryptionConfiguration.Rules[0]
			if rule.ApplyServerSideEncryptionByDefault != nil && rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm != nil {
				bucketInfo.EncryptionType = *rule.ApplyServerSideEncryptionByDefault.SSEAlgorithm
			}
			runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket %s has encryption enabled, type: %s", *bucket.Name, bucketInfo.EncryptionType))
		} else if err != nil {
			// Check if the error is expected configuration not found error
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "NoSuchBucketPolicy" {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get encryption for %s: %s", *bucket.Name, err.Error()))
			} else if !ok {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get encryption for %s: %s", *bucket.Name, err.Error()))
			}
		}

		// Get Lifecycle rules status
		lcInput := &s3.GetBucketLifecycleConfigurationInput{Bucket: bucket.Name}
		lcOutput, err := s3Manager.client.GetBucketLifecycleConfigurationWithContext(ctx, lcInput)
		if err == nil && len(lcOutput.Rules) > 0 {
			bucketInfo.HasLifecycleRules = true
			bucketInfo.LifecycleRulesCount = len(lcOutput.Rules)
			runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket %s has %d lifecycle rules", *bucket.Name, bucketInfo.LifecycleRulesCount))
		} else if err != nil {
			// Check if the error is expected configuration not found error
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "NoSuchBucketPolicy" {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get lifecycle rules for %s: %s", *bucket.Name, err.Error()))
			} else if !ok {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get lifecycle rules for %s: %s", *bucket.Name, err.Error()))
			}
		}

		// Get Website hosting status
		webInput := &s3.GetBucketWebsiteInput{Bucket: bucket.Name}
		_, err = s3Manager.client.GetBucketWebsiteWithContext(ctx, webInput)
		if err == nil {
			bucketInfo.WebsiteEnabled = true
			runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket %s has static website hosting enabled", *bucket.Name))
		} else if err != nil {
			// Check if the error is expected configuration not found error
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() != "NoSuchBucketPolicy" {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get website config for %s: %s", *bucket.Name, err.Error()))
			} else if !ok {
				runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed to get website config for %s: %s", *bucket.Name, err.Error()))
			}
		}

		nodeBucketInfo.Buckets = append(nodeBucketInfo.Buckets, bucketInfo)
		runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket info: %s, Objects: %d, Size: %d bytes",
			bucketInfo.Name, bucketInfo.TotalObjects, bucketInfo.UsedSpace))
	}

	allNodesBucketInfo = append(allNodesBucketInfo, nodeBucketInfo)
	jsonData, err := json.MarshalIndent(allNodesBucketInfo, "", "  ")
	if err != nil {
		runtime.LogError(ContextX, "v1: JSON serialization failed: "+err.Error())
		// Still return the data collected so far
		return allNodesBucketInfo
	}
	fmt.Println("v1: allNodesBucketInfo content:")
	fmt.Println(string(jsonData))
	runtime.LogDebug(ContextX, "v1: Complete node bucket info:\n"+string(jsonData))

	return allNodesBucketInfo
}

// ObjectInfo struct remains the same
type ObjectInfo struct {
	Key          string            `json:"key"`          // Object key
	Size         int64             `json:"size"`         // Object size (bytes)
	LastModified time.Time         `json:"lastModified"` // Last modified time
	ETag         string            `json:"etag"`         // ETag value
	ContentType  string            `json:"contentType"`  // Content type
	StorageClass string            `json:"storageClass"` // Storage class
	Metadata     map[string]string `json:"metadata"`     // User-defined metadata
	VersionId    string            `json:"versionId"`    // Version ID (if versioning enabled)
}

// UploadObject uploads a file to the specified bucket using AWS SDK v1
func (a *S3Manager) UploadObject(endpoint, accessKey, secretKey, region, bucketName string) (string, error) {
	filePath := file.GetFilePath(ContextX)
	if filePath == "" {
		return "", fmt.Errorf("file selection cancelled or no file chosen")
	}
	runtime.LogDebug(ContextX, "v1: Upload file path => "+filePath)

	fileName := filepath.Base(filePath)
	objectKey := fileName // Use filename as object key

	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Uploading %s to bucket %s with key %s", filePath, bucketName, objectKey))

	// Create S3 client specifically for this operation (could reuse if manager instance persists)
	// Using the same configuration logic as NewS3Client
	awsCfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(awsCfg)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to create session for upload: "+err.Error())
		return "", fmt.Errorf("v1: failed to create upload session: %v", err)
	}
	client := s3.New(sess)

	// Open the file
	fileHandle, err := os.Open(filePath)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to open file: "+err.Error())
		return "", fmt.Errorf("v1: unable to open file %s: %v", filePath, err)
	}
	defer fileHandle.Close()

	// Get file info for size
	fileInfo, err := fileHandle.Stat()
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to get file info: "+err.Error())
		return "", fmt.Errorf("v1: unable to get file info: %v", err)
	}

	// Detect content type
	buffer := make([]byte, 512)
	_, err = fileHandle.Read(buffer)
	// We need to reset the file pointer after reading the buffer
	_, seekErr := fileHandle.Seek(0, 0)
	if seekErr != nil {
		runtime.LogError(ContextX, "v1: Failed to seek file pointer: "+seekErr.Error())
		return "", fmt.Errorf("v1: unable to reset file pointer: %v", seekErr)
	}
	// Handle potential read error *after* seeking back
	if err != nil && err != io.EOF {
		runtime.LogError(ContextX, "v1: Failed to read file for content type detection: "+err.Error())
		return "", fmt.Errorf("v1: unable to read file for content type: %v", err)
	}
	contentType := http.DetectContentType(buffer)

	// Prepare PutObject input
	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          fileHandle, // Pass the file handle (implements io.ReadSeeker)
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(contentType),
		Expires:       aws.Time(time.Now().Add(24 * time.Hour)), // Set expiry as before
		// Add other parameters like ACL, Metadata if needed
		// Metadata: map[string]*string{"custom-key": aws.String("custom-value")},
	}

	// Upload the file
	ctx := context.Background()
	runtime.LogDebug(ContextX, "v1: --------------------Uploading--------------------")
	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Bucket: %s, Key: %s, ContentType: %s", bucketName, objectKey, contentType))

	result, err := client.PutObjectWithContext(ctx, putInput)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to upload file: "+err.Error())
		if aerr, ok := err.(awserr.Error); ok {
			runtime.LogError(ContextX, fmt.Sprintf("v1: AWS Error Code: %s, Message: %s", aerr.Code(), aerr.Message()))
		}
		return "", fmt.Errorf("v1: upload failed: %v", err)
	}

	etag := ""
	if result.ETag != nil {
		etag = *result.ETag
	}
	runtime.LogDebug(ContextX, fmt.Sprintf("v1: File uploaded successfully, ETag: %s", etag))
	return etag, nil
}

// DownloadObject downloads an object from the specified bucket using AWS SDK v1
func (a *S3Manager) DownloadObject(endpoint, accessKey, secretKey, region, bucketName, objectKey string) error {
	savePathDir := file.GetDirPath(ContextX)
	if savePathDir == "" {
		return fmt.Errorf("v1: save directory selection cancelled or failed")
	}
	savePath := filepath.Join(savePathDir, objectKey) // Use filepath.Join for cross-platform paths

	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Downloading object %s from bucket %s to %s", objectKey, bucketName, savePath))

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		runtime.LogError(ContextX, "v1: Failed to create directory: "+err.Error())
		return fmt.Errorf("v1: unable to create directory %s: %v", filepath.Dir(savePath), err)
	}

	// Create S3 client
	awsCfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess, err := session.NewSession(awsCfg)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to create session for download: "+err.Error())
		return fmt.Errorf("v1: failed to create download session: %v", err)
	}
	client := s3.New(sess)

	// Create the output file
	outFile, err := os.Create(savePath)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to create file: "+err.Error())
		return fmt.Errorf("v1: unable to create file %s: %v", savePath, err)
	}
	defer outFile.Close()

	// Prepare GetObject input
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	// Download the object
	ctx := context.Background()
	result, err := client.GetObjectWithContext(ctx, getInput)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to download object: "+err.Error())
		if aerr, ok := err.(awserr.Error); ok {
			runtime.LogError(ContextX, fmt.Sprintf("v1: AWS Error Code: %s, Message: %s", aerr.Code(), aerr.Message()))
		}
		// Attempt to remove the potentially empty/partially created file on error
		_ = os.Remove(savePath)
		return fmt.Errorf("v1: download failed: %v", err)
	}
	defer result.Body.Close() // Ensure the response body is closed

	// Write the body to the file
	bytesCopied, err := io.Copy(outFile, result.Body)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to write to file: "+err.Error())
		// Attempt to remove the potentially incomplete file on error
		_ = os.Remove(savePath)
		return fmt.Errorf("v1: file write failed: %v", err)
	}

	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Object downloaded successfully, %d bytes written", bytesCopied))
	return nil
}

// ListObjects retrieves object information from a bucket using AWS SDK v1
func (a *S3Manager) ListObjects(endpoint, accessKey, secretKey, region, bucketName string) ([]ObjectInfo, error) {
	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Listing objects in bucket %s", bucketName))

	s3Manager, err := a.NewS3Client(endpoint, region, accessKey, secretKey)
	if err != nil {
		runtime.LogError(ContextX, "v1: Connection failed for ListObjects: "+err.Error())
		return nil, err
	}

	var objectInfoList []ObjectInfo
	ctx := context.Background()
	listInput := &s3.ListObjectsV2Input{Bucket: aws.String(bucketName)}

	err = s3Manager.client.ListObjectsV2PagesWithContext(ctx, listInput,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			if page == nil {
				return false // Stop if page is nil
			}
			for _, obj := range page.Contents {
				// Basic info available directly from ListObjectsV2Output
				info := ObjectInfo{
					Key:          aws.StringValue(obj.Key),        // Use aws.StringValue for safe dereference
					Size:         aws.Int64Value(obj.Size),        // Use aws.Int64Value
					LastModified: aws.TimeValue(obj.LastModified), // Use aws.TimeValue
					ETag:         aws.StringValue(obj.ETag),
				}
				if obj.StorageClass != nil {
					info.StorageClass = *obj.StorageClass
				}

				// Getting detailed info (ContentType, Metadata, VersionId) requires HeadObject
				// This part remains optional as it significantly increases API calls
				if false { // Set to true to enable HeadObject calls per object
					headInput := &s3.HeadObjectInput{
						Bucket: aws.String(bucketName),
						Key:    obj.Key, // Key is already a pointer here
					}
					headOutput, headErr := s3Manager.client.HeadObjectWithContext(ctx, headInput)
					if headErr == nil {
						info.ContentType = aws.StringValue(headOutput.ContentType)
						if headOutput.Metadata != nil {
							info.Metadata = make(map[string]string)
							for k, v := range headOutput.Metadata {
								info.Metadata[k] = aws.StringValue(v) // Metadata values are *string in v1
							}
						}
						info.VersionId = aws.StringValue(headOutput.VersionId)
					} else {
						runtime.LogDebug(ContextX, fmt.Sprintf("v1: Failed HeadObject for %s: %s", info.Key, headErr.Error()))
					}
				}

				objectInfoList = append(objectInfoList, info)
			}
			return !lastPage // Continue pagination if not the last page
		})

	if err != nil {
		runtime.LogError(ContextX, "v1: Failed to list objects pages: "+err.Error())
		return nil, fmt.Errorf("v1: failed listing objects: %v", err)
	}

	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Successfully listed %d objects in bucket %s", len(objectInfoList), bucketName))
	return objectInfoList, nil
}

// GetObjectInfo retrieves detailed information for a single object using AWS SDK v1
func (a *S3Manager) GetObjectInfo(endpoint, accessKey, secretKey, region, bucketName, objectKey string) (*ObjectInfo, error) {
	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Getting info for object %s in bucket %s", objectKey, bucketName))

	s3Manager, err := a.NewS3Client(endpoint, region, accessKey, secretKey)
	if err != nil {
		runtime.LogError(ContextX, "v1: Connection failed for GetObjectInfo: "+err.Error())
		return nil, err
	}

	ctx := context.Background()
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	result, err := s3Manager.client.HeadObjectWithContext(ctx, headInput)
	if err != nil {
		runtime.LogError(ContextX, "v1: Failed HeadObject: "+err.Error())
		if aerr, ok := err.(awserr.Error); ok {
			runtime.LogError(ContextX, fmt.Sprintf("v1: AWS Error Code: %s, Message: %s", aerr.Code(), aerr.Message()))
		}
		return nil, fmt.Errorf("v1: failed to get object metadata: %v", err)
	}

	// Extract information using aws helper functions for safe pointer dereferencing
	objectInfo := &ObjectInfo{
		Key:          objectKey,
		Size:         aws.Int64Value(result.ContentLength),
		LastModified: aws.TimeValue(result.LastModified),
		ETag:         aws.StringValue(result.ETag),
		ContentType:  aws.StringValue(result.ContentType),
		StorageClass: aws.StringValue(result.StorageClass),
		VersionId:    aws.StringValue(result.VersionId),
	}

	if result.Metadata != nil {
		objectInfo.Metadata = make(map[string]string)
		for k, v := range result.Metadata {
			objectInfo.Metadata[k] = aws.StringValue(v) // Metadata values are *string in v1
		}
	}

	runtime.LogDebug(ContextX, fmt.Sprintf("v1: Object info retrieved: Key=%s, Size=%d, LastModified=%v",
		objectInfo.Key, objectInfo.Size, objectInfo.LastModified))

	return objectInfo, nil
}

// AddNode remains the same as it doesn't interact with AWS SDK
func (a *S3Manager) AddNode(name, endpoint, accessKey, secretKey, region string) bool {
	runtime.LogDebug(ContextX, "Adding node (no SDK change)") // Added log clarification
	fileContent, err := os.ReadFile("./config.json")
	runtime.LogDebug(ContextX, string(fileContent))
	if err != nil {
		runtime.LogError(ContextX, err.Error())
		return false
	}
	nodeList, err := nodes.GetNodes(fileContent)
	if err != nil {
		runtime.LogError(ContextX, "Failed to parse node list: "+err.Error())
		return false
	}

	// Add the new node
	newNode := nodes.Node{
		NodeName:  name,
		EndPoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
	}
	nodeList = nodes.AddNode(nodeList, newNode) // Assuming nodes.AddNode handles appending

	// Write the updated list back to the file
	updatedContent, err := json.MarshalIndent(nodeList, "", "    ")
	if err != nil {
		runtime.LogError(ContextX, "Failed to serialize node list: "+err.Error())
		return false
	}

	err = os.WriteFile("./config.json", updatedContent, 0644)
	if err != nil {
		runtime.LogError(ContextX, "Failed to write config file: "+err.Error())
		return false
	}

	runtime.LogDebug(ContextX, "Node added successfully")
	return true
}

// main function would remain largely the same, initializing Wails with the App struct.
// Ensure your go.mod file references github.com/aws/aws-sdk-go instead of v2.
// You might need to run `go mod tidy` after updating imports.
