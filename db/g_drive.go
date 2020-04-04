package db

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"log"
)

type Drive interface {
	UploadFloorPlan(name string, content io.Reader) (string, error)
	UploadArchiveDataFile(name string, content io.Reader) error
	ListAllFiles() ([]*drive.File, error)
}

type GDrive struct {
	srv *drive.Service
}

const (
	FloorPlanFolderName = "floor-plans"
	ArchiveFolderName   = "archive"
	RootFolderName      = "root"
)

var DirNotFound = errors.New("directory not found")

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
	directory, err := d.getDirectory(name)
	if err != nil && err != DirNotFound {
		return nil, err
	}
	if directory != nil {
		return directory, nil
	}
	log.Println("folder doesnt exist; creating new one", directory, parentId)
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

func (d GDrive) getDirectory(name string) (*drive.File, error) {
	list, err := d.srv.Files.List().
		Q(fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder'", name)).
		Do()
	if err != nil {
		return nil, err
	}
	if len(list.Files) > 0 {
		return list.Files[0], nil
	}
	return nil, DirNotFound
}

func (d GDrive) ListAllFiles() ([]*drive.File, error) {
	list, err := d.srv.Files.List().Fields("files(id,name,md5Checksum,mimeType,size,createdTime,parents)").OrderBy("name").Do()
	if err != nil {
		return nil, nil
	}
	//for _, f := range list.Files {
	//	log.Printf("======= Name: %s, Id: %s, Parent: %s, Mime: %s", f.Name, f.Id, f.Parents, f.MimeType)
	//log.Printf("-- %+v", f)
	//}
	return list.Files, nil
}

func (d GDrive) UploadArchiveDataFile(name string, content io.Reader) error {
	dir, err := d.createDir(ArchiveFolderName, RootFolderName)
	if err != nil {
		log.Println("Failed to create folder: " + err.Error())
		return err
	}
	_, err = d.createFile(name, "text/plain", content, dir.Id)
	if err != nil {
		log.Printf("Could not create file: %v\n", err)
		return err
	}
	return nil
}
