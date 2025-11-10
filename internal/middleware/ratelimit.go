package middleware

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type rateLimiter struct {
	mu               sync.Mutex
	userLastAction   map[int64]time.Time
	userMessageCount map[int64]int
}

var limiter = &rateLimiter{
	userLastAction:   make(map[int64]time.Time),
	userMessageCount: make(map[int64]int),
}

func getUserIdentifier(user *models.User) string {
	if user.Username != "" {
		return fmt.Sprintf("@%s", user.Username)
	}
	if user.FirstName != "" {
		if user.LastName != "" {
			return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
		return user.FirstName
	}
	return fmt.Sprintf("ID:%d", user.ID)
}

func RateLimit(minInterval time.Duration, maxMessages int, period time.Duration) func(next bot.HandlerFunc) bot.HandlerFunc {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			limiter.cleanup(period)
		}
	}()

	return func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update.Message == nil {
				next(ctx, b, update)
				return
			}

			userID := update.Message.From.ID

			if IsAdmin(userID) {
				next(ctx, b, update)
				return
			}

			now := time.Now()
			userIdent := getUserIdentifier(update.Message.From)

			limiter.mu.Lock()

			lastAction, exists := limiter.userLastAction[userID]
			if exists && now.Sub(lastAction) < minInterval {
				limiter.mu.Unlock()

				log.Printf("⚠️ Rate limit: User %d (%s) too fast (interval: %v)",
					userID, userIdent, now.Sub(lastAction))

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "⏳ Будь ласка, зачекайте трохи перед наступним повідомленням.",
				})
				return
			}

			count := limiter.userMessageCount[userID]
			if count >= maxMessages {
				limiter.mu.Unlock()

				log.Printf("⚠️ Rate limit: User %d (%s) exceeded max messages (%d/%d)",
					userID, userIdent, count, maxMessages)

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "⏳ Ви надіслали забагато повідомлень. Спробуйте через хвилину.",
				})
				return
			}

			limiter.userLastAction[userID] = now
			limiter.userMessageCount[userID] = count + 1
			limiter.mu.Unlock()

			next(ctx, b, update)
		}
	}
}

func (rl *rateLimiter) cleanup(period time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	for userID, lastTime := range rl.userLastAction {
		if now.Sub(lastTime) > period {
			delete(rl.userLastAction, userID)
			delete(rl.userMessageCount, userID)
		}
	}

	log.Printf("🧹 Rate limiter cleanup: %d users tracked", len(rl.userLastAction))
}
