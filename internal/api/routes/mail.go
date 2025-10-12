package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vky5/mailcat/internal/db"
	"github.com/vky5/mailcat/internal/imap"
	"gorm.io/gorm"
)

func RegisterMailRoutes(r *gin.Engine) {
	acc := r.Group("/mail")
	{
		acc.POST("/stream", streamMails)
	}
}

type StreamRequest struct {
	Email      string `json:"email"`
	Mailbox    string `json:"mailbox"`
	PageSize   int    `json:"pagesize"`
	PageNumber int    `json:"pagenumber"`
}

func streamMails(c *gin.Context) {
	var req StreamRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// setting up header to keep connection open for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// flush support (Streaming)
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "streaming unsupported"})
		return
	}

	// search for the account through email
	account, err := db.GetAccountByEmail(req.Email)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Someting went wrong"})
			return
		}
	}

	// connect to imap
	conn, err := imap.ConnectIMAP(*account) // this will stuck here until the the connection is established
	// what I am thinking is that I create a channel of the pagesize automatically and then fetch the messages from the pagesize using range and then continously stream it but better stream new emails because
	// if we try to stream mails like that seqset takes gives mail in ascending order 41 42 ... 50 and if we want 50th at the first place it is way more headache
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to imap"})
		return
	}
	defer imap.Logout(conn)

	emails, err := imap.FetchEmails(conn, req.PageSize, req.PageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// encode full batch as JSON
	data, err := json.Marshal(emails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encode emails"})
		return
	}

	// send whole batch as one SSE event
	fmt.Fprintf(c.Writer, "event: batch\n")
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()

}
