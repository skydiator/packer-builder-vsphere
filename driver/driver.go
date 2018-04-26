package driver

import (
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"context"
	"net/url"
	"fmt"
	"github.com/vmware/govmomi/object"
	"time"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25"
)

type Driver struct {
	ctx        context.Context
	client     *govmomi.Client
	finder     *find.Finder
	datacenter *object.Datacenter
}

type ConnectConfig struct {
	VCenterServer      string
	Username           string
	Password           string
	InsecureConnection bool
	Datacenter         string
}

func NewDriver(ctx context.Context, config *ConnectConfig) (*Driver, error) {

	vcenter_url, err := url.Parse(fmt.Sprintf("https://%v/sdk", config.VCenterServer))
	if err != nil {
		return nil, err
	}
	credentials := url.UserPassword(config.Username, config.Password)
	vcenter_url.User = credentials

	soapClient := soap.NewClient(vcenter_url, config.InsecureConnection)
	vimClient, err := vim25.NewClient(ctx, soapClient)
	if err != nil {
		return nil, err
	}

	vimClient.RoundTripper = session.KeepAlive(vimClient.RoundTripper, 10*time.Minute)
	client := &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}

	err = client.SessionManager.Login(ctx, credentials)
	if err != nil {
		return nil, err
	}

	finder := find.NewFinder(client.Client, false)
	datacenter, err := finder.DatacenterOrDefault(ctx, config.Datacenter)
	if err != nil {
		return nil, err
	}
	finder.SetDatacenter(datacenter)

	d := Driver{
		ctx:        ctx,
		client:     client,
		datacenter: datacenter,
		finder:     finder,
	}
	return &d, nil
}
