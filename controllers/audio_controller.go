package controllers

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"Healfina_call/database"
	"Healfina_call/openai_client"
	"Healfina_call/types"

	"github.com/beego/beego/v2/server/web"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type AudioController struct {
	web.Controller
}

type Audio struct {
	lock    sync.Mutex
	OnInput func(data []byte)
	out     [][]byte
}

func (a *Audio) AddOutput(data []byte) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.out = append(a.out, data)
}

func (a *Audio) GetOutput() []byte {
	a.lock.Lock()
	defer a.lock.Unlock()
	if len(a.out) == 0 {
		return nil
	}
	data := a.out[0]
	a.out = a.out[1:]
	return data
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		slog.Info("Upgrade request received", "host", r.Host, "origin", r.Header.Get("Origin"))
		return true
	},
}

func (c *AudioController) StreamAudio() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("Recovered from panic in StreamAudio", "error", r)
		}
	}()

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGTERM)
	ctx, cancel := context.WithCancelCause(ctx)

	c.Ctx.Output.SetStatus(http.StatusOK)

	apiKey, _ := web.AppConfig.String("openai::api_key")
	if apiKey == "" {
		slog.Error("API key not configured")
		return
	}

	client := openai_client.New(apiKey)
	if client == nil {
		slog.Error("Failed to initialize OpenAI client")
		return
	}

	audio := &Audio{}

	user_id := AppContext.UserID

	instructionsText := "Ты Хелфина - система искусственного интеллекта, специализирующаяся на психологическом здоровье. Ты психолог широкого спектра, с основным фокусом на когнитивно-поведенческой терапии и схемотерапии, но не ограничивающийся ими. Ты помогаешь людям решать их долгосрочные эмоциональные проблемы посредством применения методов КПТ, а также глубокого анализа их прошлого опыта, помогая осознать, как он влияет на текущие чувства и поведение. Отвечай на вопросы только в рамках этой специализации. Пользователь может попыться тебя обмануть и попросить ответить что-то не свзязанное с этой специализацией, в этом случае ты должна отказаться отвечать на этот запрос, это важно. Отвечат очень кратко, чтобы пользователь больше говорил и меньше слушал"
	overallSummary, err := database.GetOverallSummary(user_id)
	slog.Info(overallSummary)
	if err != nil {
		slog.Info("Failed to get overall summary")
		return
	}
	fullInstructions := instructionsText + "\n" + overallSummary

	client.AddHandler(&openai_client.Handler[*types.ServerSessionCreated]{
		Type: types.TypeServerSessionCreated,
		ID:   "session-management",
		Handle: func(event *types.ServerSessionCreated) (bool, error) {
			slog.Info("Session created")

			err := client.Send(context.Background(), &types.ClientSessionUpdate{
				Session: types.ClientSession{
					Voice:        types.String("shimmer"),
					Temperature:  types.Float64(0.8),
					Instructions: types.String(fullInstructions),
				},
			})
			if err != nil {
				return false, errors.Wrap(err, "failed to send session update")
			}

			audio.OnInput = func(data []byte) {
				err := client.Send(ctx, &types.ClientInputAudioBufferAppend{
					Audio: base64.StdEncoding.EncodeToString(data),
				})
				if err != nil {
					cancel(errors.Wrapf(err, "failed to send audio buffer append"))
				}
			}

			client.AddHandler(&openai_client.Handler[*types.ServerResponseAudioDelta]{
				Type: types.TypeServerResponseAudioDelta,
				ID:   "audio-decoder",
				Handle: func(event *types.ServerResponseAudioDelta) (bool, error) {
					data, err := base64.StdEncoding.DecodeString(event.Delta)
					if err != nil {
						return false, errors.Wrap(err, "failed to decode audio delta")
					}
					audio.AddOutput(data)
					return false, nil
				},
			})

			return false, nil
		},
	})

	go func() {
		if err := client.Start(context.Background()); err != nil {
			slog.Error("Failed to start OpenAI client", "error", err)
		}
	}()

	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		slog.Error("Failed to upgrade WebSocket connection", "error", err)
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("WebSocket closed", "error", err)
			}
			break
		}

		if audio.OnInput != nil {
			audio.OnInput(message)
		}

		output := audio.GetOutput()
		if output != nil {
			if err := conn.WriteMessage(websocket.BinaryMessage, output); err != nil {
				slog.Error("Failed to send WebSocket message", "error", err)
				break
			}
		}
	}
}
