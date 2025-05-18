package internals

import (
	"time"
)

func SetUserOfflineEvery5Seconds(userRepo *UserRepo) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			onlineUsers, err := userRepo.GetAllOnlineUsers()
			if err != nil {
				continue
			}
			for _, user := range onlineUsers {
				userID := user["id"].(float64)
				err := userRepo.SetOffline(int(userID))
				if err != nil {
					continue
				}
			}
		}
	}
}
