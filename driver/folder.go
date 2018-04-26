package driver

import (
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/vim25/mo"
	"fmt"
)

type Folder struct {
	driver *Driver
	folder *object.Folder
}

func (d *Driver) NewFolder(ref *types.ManagedObjectReference) *Folder {
	return &Folder{
		folder: object.NewFolder(d.client.Client, *ref),
		driver: d,
	}
}

func (d *Driver) FindFolder(name string) (*Folder, error) {
	f, err := d.finder.Folder(d.ctx, fmt.Sprintf("%v", name))
	if err != nil {
		return nil, err
	}
	return &Folder{
		folder: f,
		driver: d,
	}, nil
}

func (f *Folder) Info(params ...string) (*mo.Folder, error) {
	var p []string
	if len(params) == 0 {
		p = []string{"*"}
	} else {
		p = params
	}
	var info mo.Folder
	err := f.folder.Properties(f.driver.ctx, f.folder.Reference(), p, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (f *Folder) Path() (string, error) {
	info, err := f.Info("name", "parent")
	if err != nil {
		return "", err
	}
	if info.Parent.Type == "Datacenter" {
		return "", nil
	} else {
		parent := f.driver.NewFolder(info.Parent)
		path, err := parent.Path()
		if err != nil {
			return "", err
		}
		if path == "" {
			return info.Name, nil
		} else {
			return fmt.Sprintf("%v/%v", path, info.Name), nil
		}
	}
}
