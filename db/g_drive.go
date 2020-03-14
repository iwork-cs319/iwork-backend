package db

import (
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
)

type Drive interface {
	UploadFloorPlan(name string, content io.Reader) (string, error)
}

type GDrive struct {
	srv *drive.Service
}

const (
	FloorPlanFolderName = "floor-plans"
	RootFolderName      = "root"
)

func NewDriveClient(driveConfigJSON string) (Drive, error) {
	service, err := drive.NewService(
		context.Background(),
		option.WithCredentialsJSON([]byte(driveConfigJSON)),
	)

	if err != nil {
		log.Printf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}

	return &GDrive{
		srv: service,
	}, err
}

func (d GDrive) createDir(name string, parentId string) (*drive.File, error) {
	dir := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentId},
	}

	file, err := d.srv.Files.Create(dir).Do()

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (d GDrive) createFile(name string, mimeType string, content io.Reader, parentId string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentId},
	}
	file, err := d.srv.Files.Create(f).Media(content).Do()

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (d GDrive) updatePermissionToPublic(id string) error {
	permissionData := &drive.Permission{
		Type: "anyone",
		Role: "reader",
		//AllowFileDiscovery: true,
	}
	_, err := d.srv.Permissions.Create(id, permissionData).Do()
	if err != nil {
		log.Println("Error updating permissions: " + err.Error())
		return err
	}
	return nil
}

func (d GDrive) UploadFloorPlan(name string, content io.Reader) (string, error) {
	dir, err := d.createDir(FloorPlanFolderName, RootFolderName)
	if err != nil {
		log.Println("Failed to create folder: " + err.Error())
		return "", err
	}
	file, err := d.createFile(name, "image/jpg", content, dir.Id)
	if err != nil {
		log.Printf("Could not create file: %v\n", err)
		return "", err
	}
	err = d.updatePermissionToPublic(file.Id)
	if err != nil {
		log.Println("Failed to update permissions: " + err.Error())
		return "", err
	}
	return file.Id, nil
}
