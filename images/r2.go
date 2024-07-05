package images

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/charisworks/charisworks-service-go/util"
	r2client "github.com/whatacotton/cloudflare-go/client"
	"github.com/whatacotton/cloudflare-go/r2"
)

var ()

type R2Conns struct {
	Crud r2.R2crud
	ctx  context.Context
}

func (r *R2Conns) Init() {
	flag.StringVar(&util.R2_ENDPOINT, "endpoint", util.R2_ENDPOINT, "endpoint")
	flag.StringVar(&util.R2_ACCOUNT_ID, "account-id", util.R2_ACCOUNT_ID, "account-id")
	flag.StringVar(&util.R2_ACCESS_KEY_ID, "access-key-id", util.R2_ACCESS_KEY_ID, "access-key-id")
	flag.StringVar(&util.R2_ACCESS_KEY_SECRET, "account-key-secret", util.R2_ACCESS_KEY_SECRET, "account-key-secret")
	flag.Parse()

	if util.R2_ENDPOINT == "" || util.R2_ACCOUNT_ID == "" || util.R2_ACCESS_KEY_ID == "" || util.R2_ACCESS_KEY_SECRET == "" {
		panic("missing required parameters")
	}
	client, err := r2client.New(
		util.R2_ACCOUNT_ID,
		util.R2_ENDPOINT,
		util.R2_ACCESS_KEY_ID,
		util.R2_ACCESS_KEY_SECRET,
	).Connect(context.TODO())
	if err != nil {
		log.Fatalf("r2 client conneciton error :%v\n", err)
	}
	r.ctx = context.Background()
	r.Crud = r2.NewR2CRUD(util.R2_BUCKET_NAME, client, 60)
	log.Println("r2 client connection success")
}

func (r *R2Conns) UploadImage(filepath string, path string) error {
	filedata, err := PathToByte(filepath)
	if err != nil {
		log.Print(err)
		return err
	}

	return r.Crud.UploadObject(r.ctx, filedata, path)
}

func (r *R2Conns) GetImages(path string) ([]string, error) {
	objects, err := r.Crud.ListObjects(r.ctx, path)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	log.Print(objects)
	var images []string
	for _, obj := range objects.Contents {
		log.Print(*obj.Key)
		images = append(images, *obj.Key)
	}
	return images, nil
}

func (r *R2Conns) DeleteImage(path string) error {
	err := r.Crud.DeleteObject(r.ctx, path)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
func PathToByte(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return fileBytes, nil
}

// 画像ファイルをバイト列に変換
func FileToByte(file *os.File) ([]byte, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	fileBytes := make([]byte, fileInfo.Size())
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return fileBytes, nil
}
