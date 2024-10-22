package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redsync/redsync/v4"
)

type handler struct {
	rs *redsync.Redsync
}

func NewHanler(rs *redsync.Redsync) *handler {
	return &handler{
		rs: rs,
	}
}

func (h handler) DoRequiredLockOperation(c *gin.Context) {
	resourceID := c.Query("resource_id")
	if resourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resource_id is required"})
		return
	}

	key := fmt.Sprintf("lock:%s", resourceID)
	lockTimeout := 30 * time.Second
	expiration := redsync.WithExpiry(lockTimeout)
	mutex := h.rs.NewMutex(key, expiration)

	ctx := context.Background()
	if err := mutex.LockContext(ctx); err != nil {
		log.Printf("Error: Failed to acquire lock for resource %s: %v", resourceID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire lock", "details": err.Error()})
		return
	}

	log.Printf("Lock acquired for resource %s, performing critical operation", resourceID)

	criticalOperation(resourceID)

	if _, err := mutex.UnlockContext(ctx); err != nil {
		log.Printf("Error: Failed to release lock for resource %s: %v", resourceID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release lock", "details": err.Error()})
		return
	}

	log.Printf("Lock released for resource %s", resourceID)
	c.JSON(http.StatusOK, gin.H{"message": "Lock acquired and released successfully", "resource_id": resourceID})
}

func criticalOperation(resourceID string) {
	time.Sleep(15 * time.Second)
	log.Printf("Critical operation completed for resource %s", resourceID)
}
