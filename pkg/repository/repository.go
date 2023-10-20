package repository

import (
	"github.com/quic-s/quics/pkg/core/history"
	"github.com/quic-s/quics/pkg/core/metadata"
	"github.com/quic-s/quics/pkg/core/registration"
	"github.com/quic-s/quics/pkg/core/server"
	"github.com/quic-s/quics/pkg/core/sharing"
	"github.com/quic-s/quics/pkg/core/sync"
)

type Repository interface {
	NewHistoryRepository() history.Repository
	NewMetadataRepository() metadata.Repository
	NewRegistrationRepository() registration.Repository
	NewServerRepository() server.Repository
	NewSharingRepository() sharing.Repository
	NewSyncRepository() sync.Repository
}
