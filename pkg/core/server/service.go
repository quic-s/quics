package server

import (
	"fmt"
	"log"

	"github.com/quic-s/quics/pkg/app"
)

type ServerService struct {
	quics            *app.App
	serverRepository Repository
}

func NewService(app *app.App, serverRepository Repository) *ServerService {
	return &ServerService{
		quics:            app,
		serverRepository: serverRepository,
	}
}

// StopServer stop quic-s server
func (ss *ServerService) StopServer() error {
	fmt.Println("************************************************************")
	fmt.Println("                           Stop                             ")
	fmt.Println("************************************************************")

	err := ss.quics.Stop()
	if err != nil {
		log.Println("quics: ", err)
		return err
	}

	return nil
}

// ListenProtocol is executed when server starts
func (ss *ServerService) ListenProtocol() error {
	fmt.Println("************************************************************")
	fmt.Println("                     Listen Protocol                        ")
	fmt.Println("************************************************************")

	// listen protocol using goroutine
	go func() {
		// listen quics-protocol
		ss.quics.Start()

		err := ss.quics.Close()
		if err != nil {
			log.Println("quics: ", err)
			return
		}

		return
	}()

	return nil
}

func (ss *ServerService) ShowClientLogs(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                        Client Log                          ")
	fmt.Println("************************************************************")

	if all != "" {
		clients, err := ss.serverRepository.GetAllClients()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, client := range clients {
			for _, root := range client.Root {
				fmt.Printf("*   UUID: %s   |   ID: %d   |   IP: %s   |   Root Directoreis: %s   *\n", client.UUID, client.Id, client.Ip, root)
			}
		}

		return nil
	}

	if id != "" {
		client, err := ss.serverRepository.GetClientByUUID(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, root := range client.Root {
			fmt.Printf("*   UUID: %s   |   ID: %d   |   IP: %s   |   Root Directory: %s   *\n", client.UUID, client.Id, client.Ip, root.AfterPath)
		}

		return nil
	}

	return nil
}

func (ss *ServerService) ShowDirLogs(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                       Directory Log                        ")
	fmt.Println("************************************************************")

	if all != "" {
		dirs, err := ss.serverRepository.GetAllRootDirectories()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, dir := range dirs {
			for _, UUID := range dir.UUIDs {
				fmt.Printf("*   Root Directory: %s   |   Owner: %s   |   Password: %s   |   UUID: %s   *\n", dir.AfterPath, dir.Owner, dir.Password, UUID)
			}
		}

		return nil
	}

	if id != "" {
		dir, err := ss.serverRepository.GetRootDirectoryByPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, UUID := range dir.UUIDs {
			fmt.Printf("*   Root Directory: %s   |   Owner: %s   |   Password: %s   |   UUID: %s   *\n", dir.AfterPath, dir.Owner, dir.Password, UUID)
		}

		return nil
	}

	return nil
}

func (ss *ServerService) ShowFileLogs(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                         File Log                           ")
	fmt.Println("************************************************************")

	if all != "" {
		files, err := ss.serverRepository.GetAllFiles()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		for _, file := range files {
			fmt.Printf("*   File: %s   |   Root Directory: %s   |   LatestHash: %s   |   LatestSyncTimestamp: %d   |   ContentsExisted: %t   |   Metadata: %s   *\n", file.AfterPath, file.RootDirKey, file.LatestHash, file.LatestSyncTimestamp, file.ContentsExisted, file.Metadata.ModTime)
		}

		return nil
	}

	if id != "" {
		file, err := ss.serverRepository.GetFileByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		fmt.Printf("*   File: %s   |   Root Directory: %s   |   LatestHash: %s   |   LatestSyncTimestamp: %d   |   ContentsExisted: %t   |   Metadata: %s   *\n", file.AfterPath, file.RootDirKey, file.LatestHash, file.LatestSyncTimestamp, file.ContentsExisted, file.Metadata.ModTime)

		return nil
	}

	return nil
}

func (ss *ServerService) DisconnectClient(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                    Disconnect Client                       ")
	fmt.Println("************************************************************")

	if all != "" {
		err := ss.serverRepository.DeleteAllClients()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	if id != "" {
		err := ss.serverRepository.DeleteClientByUUID(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	return nil
}

func (ss *ServerService) DisconnectDir(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                  Disconnect Directory                      ")
	fmt.Println("************************************************************")

	if all != "" {
		err := ss.serverRepository.DeleteAllRootDirectories()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	if id != "" {
		err := ss.serverRepository.DeleteRootDirectoryByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	return nil
}

func (ss *ServerService) DisconnectFile(all string, id string) error {
	fmt.Println("************************************************************")
	fmt.Println("                     Disconnect File                        ")
	fmt.Println("************************************************************")

	if all != "" {
		err := ss.serverRepository.DeleteAllFiles()
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	if id != "" {
		err := ss.serverRepository.DeleteFileByAfterPath(id)
		if err != nil {
			log.Println("quics: ", err)
			return err
		}

		return nil
	}

	return nil
}
