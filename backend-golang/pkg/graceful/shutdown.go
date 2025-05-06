package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"sample-stack-golang/pkg/logger"
)

// Package graceful menyediakan utilitas untuk shutdown aplikasi dengan graceful.
//
// Package ini mengimplementasikan mekanisme untuk memastikan bahwa aplikasi shutdown dengan graceful,
// memungkinkan operasi yang sedang berlangsung untuk diselesaikan sebelum ditutup. Package ini termasuk:
//   - Penanganan sinyal untuk SIGINT dan SIGTERM
//   - Pelacakan wait group untuk request aktif dan task background
//   - Shutdown delay yang dapat diatur
//   - Middleware untuk melacak request HTTP
//
// Contoh penggunaan:
//
//	sm := graceful.NewShutdownManager(echoServer, service.Close)
//	e.Use(sm.WaitGroupMiddleware())
//	go sm.WaitForShutdown()
//
// Untuk task background:
//
//	sm.AddTask()
//	go func() {
//	    defer sm.DoneTask()
//	    // Lakukan pekerjaan background
//	}()

// ShutdownManager menangani proses graceful shutdown aplikasi
type ShutdownManager struct {
	wg            sync.WaitGroup
	shutdownChan  chan os.Signal
	shutdownDelay time.Duration
	server        *echo.Echo
	closeFunc     func() error
}

// NewShutdownManager membuat instance baru dari ShutdownManager
func NewShutdownManager(server *echo.Echo, closeFunc func() error) *ShutdownManager {
	return &ShutdownManager{
		wg:            sync.WaitGroup{},
		shutdownChan:  make(chan os.Signal, 1),
		shutdownDelay: 10 * time.Second,
		server:        server,
		closeFunc:     closeFunc,
	}
}

// WaitGroupMiddleware adalah middleware yang melacak request HTTP yang sedang aktif
func (sm *ShutdownManager) WaitGroupMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sm.wg.Add(1)
			defer sm.wg.Done()
			return next(c)
		}
	}
}

// AddTask menambah task ke wait group
func (sm *ShutdownManager) AddTask() {
	sm.wg.Add(1)
}

// DoneTask menandai task sudah selesai
func (sm *ShutdownManager) DoneTask() {
	sm.wg.Done()
}

// WaitForShutdown menunggu sinyal shutdown dan melakukan graceful shutdown
func (sm *ShutdownManager) WaitForShutdown() {
	// Registrasi untuk SIGINT dan SIGTERM
	signal.Notify(sm.shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Menunggu sinyal shutdown
	<-sm.shutdownChan
	logger.Log.Info("Sinyal shutdown diterima, memulai proses graceful shutdown")

	// Membuat context dengan timeout untuk shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), sm.shutdownDelay)
	defer cancel()

	// Berhenti menerima request baru
	logger.Log.Info("Menghentikan HTTP server (tidak menerima request baru)")
	if err := sm.server.Shutdown(ctx); err != nil {
		logger.Log.WithError(err).Error("Terjadi error saat shutdown server")
	}

	// Menunggu semua request/proses aktif selesai
	logger.Log.Info("Menunggu semua proses aktif selesai")
	waitChan := make(chan struct{})
	go func() {
		sm.wg.Wait()
		close(waitChan)
	}()

	// Menunggu hingga semua proses selesai atau timeout
	select {
	case <-waitChan:
		logger.Log.Info("Semua proses aktif selesai dengan sukses")
	case <-ctx.Done():
		logger.Log.Warn("Timeout menunggu proses selesai, beberapa proses mungkin dihentikan paksa")
	}

	// Menutup semua resource
	logger.Log.Info("Menutup semua resource")
	if err := sm.closeFunc(); err != nil {
		logger.Log.WithError(err).Error("Terjadi error saat menutup resource")
	}

	logger.Log.Info("Graceful shutdown selesai")
}