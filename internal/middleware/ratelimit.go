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
	mu             sync.Mutex
	userLastAction map[int64]time.Time
	userMessages   map[int64][]time.Time
	lastWarn       map[int64]time.Time
}

var limiter = &rateLimiter{
	userLastAction: make(map[int64]time.Time),
	userMessages:   make(map[int64][]time.Time),
	lastWarn:       make(map[int64]time.Time),
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

func (rl *rateLimiter) shouldSendWarning(userID int64, now time.Time) bool {
	lastWarnTime, exists := rl.lastWarn[userID]
	if !exists || now.Sub(lastWarnTime) > 5*time.Second {
		rl.lastWarn[userID] = now
		return true
	}
	return false
}

func RateLimit(minInterval time.Duration, maxMessages int, period time.Duration) func(next bot.HandlerFunc) bot.HandlerFunc {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			limiter.cleanup()
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
				shouldWarn := limiter.shouldSendWarning(userID, now)
				limiter.mu.Unlock()

				log.Printf("⚠️ Too fast: %s (%d), interval:  %v", userIdent, userID, now.Sub(lastAction))

				if shouldWarn {
					go func() {
						time.Sleep(100 * time.Millisecond)

						_, err := b.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: update.Message.Chat.ID,
							Text:   "⏳ Будь ласка, зачекайте трохи перед наступним повідомленням.",
						})
						if err != nil {
							log.Printf("❌ SendMessage error (interval): %v", err)
						}
					}()
				}
				return
			}

			messages := limiter.userMessages[userID]
			validMessages := make([]time.Time, 0, len(messages))

			for _, msgTime := range messages {
				if now.Sub(msgTime) <= period {
					validMessages = append(validMessages, msgTime)
				}
			}

			if len(validMessages) >= maxMessages {
				oldestMsg := validMessages[0]
				waitTime := period - now.Sub(oldestMsg)
				if waitTime < 1*time.Second {
					waitTime = 1 * time.Second
				}

				shouldWarn := limiter.shouldSendWarning(userID, now)
				limiter.mu.Unlock()

				log.Printf("⚠️ Rate limit: User %d (%s) exceeded max messages (%d/%d), wait: %v",
					userID, userIdent, len(validMessages), maxMessages, waitTime.Round(time.Second))

				if shouldWarn {
					go func() {
						time.Sleep(100 * time.Millisecond)

						text := fmt.Sprintf("⏳ Ви надіслали забагато повідомлень. Спробуйте через %d секунд.",
							int(waitTime.Seconds())+1)

						_, err := b.SendMessage(ctx, &bot.SendMessageParams{
							ChatID: update.Message.Chat.ID,
							Text:   text,
						})
						if err != nil {
							log.Printf("❌ SendMessage error: %v", err)
						}
					}()
				}
				return
			}

			limiter.userLastAction[userID] = now
			limiter.userMessages[userID] = append(validMessages, now)
			limiter.mu.Unlock()

			next(ctx, b, update)
		}
	}
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cleanupThreshold := 10 * time.Minute

	for userID, lastTime := range rl.userLastAction {
		if now.Sub(lastTime) > cleanupThreshold {
			delete(rl.userLastAction, userID)
			delete(rl.userMessages, userID)
			delete(rl.lastWarn, userID)
		}
	}

	if len(rl.userLastAction) > 0 {
		log.Printf("🧹 Rate limiter cleanup: %d users tracked", len(rl.userLastAction))
	}
}
