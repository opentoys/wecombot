package wecombot

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// MediaUploader provides temporary media upload via the long connection.
// The upload process has three steps: Init → Chunks → Finish.
//
// Limits:
//   - Image: PNG/JPG/GIF, max 10MB
//   - Voice: AMR, max 2MB
//   - Video: MP4, max 10MB
//   - File: any, max 20MB
//   - Chunk size: max 512KB (before base64 encoding)
//   - Max chunks: 100
//   - Upload session validity: 30 minutes
//   - Media file validity after upload: 3 days
const (
	maxChunkSize    = 512 * 1024 // 512KB before base64
	maxTotalChunks  = 100
	defaultChunkSize = 500 * 1024 // leave some margin for base64 overhead
)

// MediaUpload holds state for an in-progress chunked upload.
type MediaUpload struct {
	client   *Client
	uploadID string
}

// MediaType for upload.
type MediaType string

const (
	MediaTypeFile  MediaType = "file"
	MediaTypeImage MediaType = "image"
	MediaTypeVoice MediaType = "voice"
	MediaTypeVideo MediaType = "video"
)

// UploadResult contains the result of a completed upload.
type UploadResult struct {
	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
}

// UploadFromFile uploads a local file as temporary media using chunked upload.
func (c *Client) UploadFromFile(mediaType MediaType, filePath string) (*UploadResult, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("wecombot: open file %s: %w", filePath, err)
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("wecombot: stat file %s: %w", filePath, err)
	}
	return c.UploadFromReader(mediaType, stat.Name(), stat.Size(), f)
}

// UploadFromReader uploads data from an io.Reader as temporary media.
func (c *Client) UploadFromReader(mediaType MediaType, filename string, totalSize int64, r io.Reader) (*UploadResult, error) {
	upload := &MediaUpload{client: c}

	// Calculate chunks
	chunkSize := defaultChunkSize
	totalChunks := int((totalSize + int64(chunkSize) - 1) / int64(chunkSize))
	if totalChunks > maxTotalChunks {
		chunkSize = int(totalSize/int64(maxTotalChunks)) + 1
		totalChunks = int((totalSize + int64(chunkSize) - 1) / int64(chunkSize))
	}

	// Compute MD5 if possible
	var md5Str string
	if mr, ok := r.(io.ReadSeeker); ok {
		h := md5.New()
		if _, err := io.Copy(h, mr); err == nil {
			md5Str = hexEncode(h.Sum(nil))
			mr.Seek(0, io.SeekStart)
		}
	}

	// Step 1: Initialize
	uploadID, err := c.uploadInit(&UploadInitBody{
		Type:        string(mediaType),
		Filename:    filename,
		TotalSize:   totalSize,
		TotalChunks: totalChunks,
		MD5:         md5Str,
	})
	if err != nil {
		return nil, fmt.Errorf("wecombot: upload init failed: %w", err)
	}
	upload.uploadID = uploadID

	// Step 2: Upload chunks
	buf := make([]byte, chunkSize)
	for i := 0; i < totalChunks; i++ {
		n, readErr := io.ReadFull(r, buf)
		if readErr != nil && readErr != io.ErrUnexpectedEOF && readErr != io.EOF {
			return nil, fmt.Errorf("wecombot: read chunk %d failed: %w", i, readErr)
		}
		data := buf[:n]

		if err := upload.UploadChunk(i, data); err != nil {
			return nil, fmt.Errorf("wecombot: upload chunk %d/%d failed: %w", i, totalChunks, err)
		}
	}

	// Step 3: Finish
	return upload.Finish()
}

// ---- Low-level upload methods ----

// uploadInit sends the init request and returns the upload_id.
func (c *Client) uploadInit(body *UploadInitBody) (string, error) {
	reqID := genReqID()
	if err := c.sendRequest(CmdUploadMediaInit, reqID, body); err != nil {
		return "", err
	}

	var resp UploadInitResponse
	if err := c.readJSON(&resp); err != nil {
		return "", err
	}
	if resp.ErrCode != 0 {
		return "", fmt.Errorf("wecombot: upload init error: code=%d msg=%s", resp.ErrCode, resp.ErrMsg)
	}
	return resp.Body.UploadID, nil
}

// UploadChunk uploads a single chunk of data.
func (u *MediaUpload) UploadChunk(index int, data []byte) error {
	encoded := base64.StdEncoding.EncodeToString(data)
	reqID := genReqID()
	return u.client.sendRequest(CmdUploadMediaChunk, reqID, &UploadChunkBody{
		UploadID:   u.uploadID,
		ChunkIndex: index,
		Base64Data: encoded,
	})
}

// Finish completes the chunked upload and returns the media_id.
func (u *MediaUpload) Finish() (*UploadResult, error) {
	reqID := genReqID()
	if err := u.client.sendRequest(CmdUploadMediaFinish, reqID, &UploadFinishBody{
		UploadID: u.uploadID,
	}); err != nil {
		return nil, err
	}

	var resp UploadFinishResponse
	if err := u.client.readJSON(&resp); err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("wecombot: upload finish error: code=%d msg=%s", resp.ErrCode, resp.ErrMsg)
	}

	return &UploadResult{
		Type:      resp.Body.Type,
		MediaID:   resp.Body.MediaID,
		CreatedAt: resp.Body.CreatedAt,
	}, nil
}

// helper to encode bytes to hex string
func hexEncode(b []byte) string {
	const hexChars = "0123456789abcdef"
	result := make([]byte, len(b)*2)
	for i, v := range b {
		result[i*2] = hexChars[v>>4]
		result[i*2+1] = hexChars[v&0x0f]
	}
	return string(result)
}
