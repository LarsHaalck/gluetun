// Package updater contains code to obtain the server information
// for the Cyberghost provider.
package updater

import (
	"context"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	possibleServers := getPossibleServers()

	possibleHosts := possibleServers.hostsSlice()
	hostToIPs, _, err := u.presolver.Resolve(ctx, possibleHosts, minServers)
	if err != nil {
		return nil, err
	}

	possibleServers.adaptWithIPs(hostToIPs)

	servers = possibleServers.toSlice()

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}